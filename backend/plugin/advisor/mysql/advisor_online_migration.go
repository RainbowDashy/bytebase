package mysql

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"

	mysql "github.com/bytebase/mysql-parser"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/parser/base"
	mysqlparser "github.com/bytebase/bytebase/backend/plugin/parser/mysql"
	"github.com/bytebase/bytebase/backend/store/model"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

var (
	_ advisor.Advisor = (*OnlineMigrationAdvisor)(nil)
)

func init() {
	advisor.Register(storepb.Engine_MYSQL, advisor.MySQLOnlineMigration, &OnlineMigrationAdvisor{})
}

// OnlineMigrationAdvisor is the advisor checking for using gh-ost to migrate large tables.
type OnlineMigrationAdvisor struct {
}

// Check checks for using gh-ost to migrate large tables.
func (*OnlineMigrationAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	if ctx.ChangeType == storepb.PlanCheckRunConfig_DDL_GHOST {
		return []advisor.Advice{
			{
				Status:  advisor.Success,
				Code:    advisor.Ok,
				Title:   "OK",
				Content: "",
			},
		}, nil
	}

	stmtList, ok := ctx.AST.([]*mysqlparser.ParseResult)
	if !ok {
		return nil, errors.Errorf("failed to convert to StmtNode")
	}

	payload, err := advisor.UnmarshalNumberTypeRulePayload(ctx.Rule.Payload)
	if err != nil {
		return nil, err
	}
	minRows := int64(payload.Number)

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &useGhostChecker{
		level:            level,
		title:            string(ctx.Rule.Type),
		currentDatabase:  ctx.CurrentDatabase,
		changedResources: make(map[string]base.SchemaResource),
	}

	for _, stmt := range stmtList {
		antlr.ParseTreeWalkerDefault.Walk(checker, stmt.Tree)
	}

	dbSchema := model.NewDBSchema(ctx.DBSchema, nil, nil)
	for _, resource := range checker.changedResources {
		tableRows := dbSchema.GetDatabaseMetadata().GetSchema(resource.Schema).GetTable(resource.Table).GetRowCount()
		if tableRows >= minRows {
			checker.adviceList = append(checker.adviceList, advisor.Advice{
				Status:  checker.level,
				Code:    advisor.AdviseOnlineMigration,
				Title:   checker.title,
				Content: fmt.Sprintf("Estimated table row count of %q is %d, greater than the set value %d. Consider using gh-ost to migrate", fmt.Sprintf("%s.%s", resource.Schema, resource.Table), tableRows, minRows),
			})
		}
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type useGhostChecker struct {
	*mysql.BaseMySQLParserListener

	adviceList []advisor.Advice
	level      advisor.Status
	title      string

	currentDatabase  string
	changedResources map[string]base.SchemaResource
}

func (c *useGhostChecker) EnterAlterTable(ctx *mysql.AlterTableContext) {
	resource := base.SchemaResource{
		Database: c.currentDatabase,
	}
	db, table := mysqlparser.NormalizeMySQLTableRef(ctx.TableRef())
	if db != "" {
		resource.Database = db
	}
	resource.Table = table
	c.changedResources[resource.String()] = resource
}
