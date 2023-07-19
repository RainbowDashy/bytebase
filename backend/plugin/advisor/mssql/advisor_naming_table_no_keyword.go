// Package mssql is the advisor for MSSQL database.
package mssql

import (
	"fmt"

	"github.com/antlr4-go/antlr/v4"
	parser "github.com/bytebase/tsql-parser"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bytebase/bytebase/backend/common/log"
	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/advisor/db"
	bbparser "github.com/bytebase/bytebase/backend/plugin/parser/sql"
)

var (
	_ advisor.Advisor = (*NamingTableNoKeywordAdvisor)(nil)
)

func init() {
	advisor.Register(db.MSSQL, advisor.MSSQLTableNamingNoKeyword, &NamingTableNoKeywordAdvisor{})
}

// NamingTableNoKeywordAdvisor is the advisor checking for table naming convention without keyword..
type NamingTableNoKeywordAdvisor struct {
}

// Check checks for table naming convention without keyword..
func (*NamingTableNoKeywordAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	tree, ok := ctx.AST.(antlr.Tree)
	if !ok {
		return nil, errors.Errorf("failed to convert to Tree")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}

	listener := &namingTableNoKeywordChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.generateAdvice()
}

// namingTableNoKeywordChecker is the listener for table naming convention without keyword.
type namingTableNoKeywordChecker struct {
	*parser.BaseTSqlParserListener

	level advisor.Status
	title string

	adviceList []advisor.Advice
}

// generateAdvice returns the advices generated by the listener, the advices must not be empty.
func (l *namingTableNoKeywordChecker) generateAdvice() ([]advisor.Advice, error) {
	if len(l.adviceList) == 0 {
		l.adviceList = append(l.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return l.adviceList, nil
}

// EnterCreate_table is called when production create_table is entered.
func (l *namingTableNoKeywordChecker) EnterCreate_table(ctx *parser.Create_tableContext) {
	tableName := ctx.Table_name().GetTable()
	normalizedTableName, err := bbparser.NormalizeTSQLTableNamePart(tableName)
	if err != nil {
		log.Error("failed to normalize table name", zap.Error(err))
	}
	if bbparser.IsTSQLKeyword(normalizedTableName, false) {
		l.adviceList = append(l.adviceList, advisor.Advice{
			Status:  l.level,
			Code:    advisor.NameIsKeywordIdentifier,
			Title:   l.title,
			Content: fmt.Sprintf("Table name [%s] is a reserved keyword and should be avoided.", normalizedTableName),
			Line:    tableName.GetStart().GetLine(),
		})
	}
}
