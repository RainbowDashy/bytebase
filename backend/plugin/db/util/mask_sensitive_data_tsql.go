package util

import (
	"cmp"
	"fmt"
	"unicode"

	"github.com/antlr4-go/antlr/v4"
	"github.com/pkg/errors"

	tsqlparser "github.com/bytebase/tsql-parser"

	"github.com/bytebase/bytebase/backend/plugin/db"
	parser "github.com/bytebase/bytebase/backend/plugin/parser/sql"
	storepb "github.com/bytebase/bytebase/proto/generated-go/store"
)

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFields(sql string) ([]db.SensitiveField, error) {
	tree, err := parser.ParseTSQL(sql)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse snowsql")
	}
	if tree == nil {
		return nil, nil
	}

	listener := &tsqlSensitiveFieldExtractorListener{
		extractor: extractor,
	}
	antlr.ParseTreeWalkerDefault.Walk(listener, tree)

	return listener.result, listener.err
}

type tsqlSensitiveFieldExtractorListener struct {
	*tsqlparser.BaseTSqlParserListener

	extractor *sensitiveFieldExtractor
	result    []db.SensitiveField
	err       error
}

// EnterSelect_statement_standalone is called when production select_statement_standalone is entered.
func (l *tsqlSensitiveFieldExtractorListener) EnterDml_clause(ctx *tsqlparser.Dml_clauseContext) {
	if ctx.Select_statement_standalone() == nil {
		return
	}

	result, err := l.extractor.extractTSqlSensitiveFieldsFromSelectStatementStandalone(ctx.Select_statement_standalone())
	if err != nil {
		l.err = err
		return
	}

	for _, field := range result {
		l.result = append(l.result, db.SensitiveField{
			Name:         field.name,
			MaskingLevel: field.maskingLevel,
		})
	}
}

// extractTSqlSensitiveFieldsFromSelectStatementStandalone extracts sensitive fields from select_statement_standalone.
func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromSelectStatementStandalone(ctx tsqlparser.ISelect_statement_standaloneContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	// TODO(zp): handle the CTE
	if ctx.With_expression() != nil {
		allCommonTableExpression := ctx.With_expression().AllCommon_table_expression()
		// TSQL do not have `RECURSIVE` keyword, if we detect `UNION`, we will treat it as `RECURSIVE`.
		for _, commonTableExpression := range allCommonTableExpression {
			var result []fieldInfo
			var err error
			normalizedCTEName := parser.NormalizeTSQLIdentifier(commonTableExpression.GetExpression_name())

			var fieldsInAnchorClause []fieldInfo
			// If statement has more than one UNION, the first one is the anchor, and the rest are recursive.
			recursiveCTE := false
			queryExpression := commonTableExpression.Select_statement().Query_expression()
			if queryExpression.Query_specification() != nil {
				fieldsInAnchorClause, err = extractor.extractTSqlSensitiveFieldsFromQuerySpecification(queryExpression.Query_specification())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract sensitive fields from `query_specification` in `query_expression`")
				}
				if allSQLUnions := queryExpression.AllSql_union(); len(allSQLUnions) > 0 {
					recursiveCTE = true
					for i := 0; i < len(allSQLUnions)-1; i++ {
						// For UNION operator, the number of the columns in the result set is the same, and will use the left part's column name.
						// So we only need to extract the sensitive fields of the right part.
						right, err := extractor.extractTSqlSensitiveFieldsFromQuerySpecification(allSQLUnions[i].Query_specification())
						if err != nil {
							return nil, errors.Wrapf(err, "failed to extract the %d set operator near line %d", i+1, allSQLUnions[i].GetStart().GetLine())
						}
						if len(fieldsInAnchorClause) != len(right) {
							return nil, errors.Wrapf(err, "the number of columns in the query statement nearly line %d returns %d fields, but %d set operator near line %d returns %d fields", ctx.GetStart().GetLine(), len(fieldsInAnchorClause), i+1, allSQLUnions[i].GetStart().GetLine(), len(right))
						}
						for j := range right {
							if cmp.Less[storepb.MaskingLevel](fieldsInAnchorClause[j].maskingLevel, right[j].maskingLevel) {
								fieldsInAnchorClause[j].maskingLevel = right[j].maskingLevel
							}
						}
					}
				}
			} else if allQueryExpression := queryExpression.AllQuery_expression(); len(allQueryExpression) > 0 {
				if len(allQueryExpression) > 1 {
					recursiveCTE = true
				}
				fieldsInAnchorClause, err = extractor.extractTSqlSensitiveFieldsFromQueryExpression(allQueryExpression[0])
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract sensitive fields from `query_specification` in `query_expression`")
				}
			}
			if !recursiveCTE {
				result = fieldsInAnchorClause
			} else {
				tempCTEOuterSchemaInfo := db.TableSchema{
					Name: normalizedCTEName,
				}
				for i := 0; i < len(fieldsInAnchorClause); i++ {
					tempCTEOuterSchemaInfo.ColumnList = append(tempCTEOuterSchemaInfo.ColumnList, db.ColumnInfo{
						Name:         fieldsInAnchorClause[i].name,
						MaskingLevel: fieldsInAnchorClause[i].maskingLevel,
					})
					result = append(result, fieldsInAnchorClause[i])
				}
				originalSize := len(extractor.cteOuterSchemaInfo)
				extractor.cteOuterSchemaInfo = append(extractor.cteOuterSchemaInfo, tempCTEOuterSchemaInfo)
				for {
					change := false
					if queryExpression.Query_specification() != nil && len(queryExpression.AllSql_union()) > 0 {
						fieldsInRecursiveClause, err := extractor.extractTSqlSensitiveFieldsFromQuerySpecification(queryExpression.AllSql_union()[len(queryExpression.AllSql_union())-1].Query_specification())
						if err != nil {
							return nil, errors.Wrapf(err, "failed to extract sensitive fields of the recursive clause of recursive CTE %q near line %d", normalizedCTEName, queryExpression.AllSql_union()[len(queryExpression.AllSql_union())-1].Query_specification().GetStart().GetLine())
						}
						if len(fieldsInRecursiveClause) != len(tempCTEOuterSchemaInfo.ColumnList) {
							return nil, errors.Wrapf(err, "recursive clause returns %d fields, but anchor clause returns %d fields in recursive CTE %q near line %d", len(fieldsInRecursiveClause), len(tempCTEOuterSchemaInfo.ColumnList), normalizedCTEName, queryExpression.AllSql_union()[len(queryExpression.AllSql_union())-1].Query_specification().GetStart().GetLine())
						}
						extractor.cteOuterSchemaInfo = extractor.cteOuterSchemaInfo[:originalSize]
						for i := 0; i < len(fieldsInRecursiveClause); i++ {
							if cmp.Less[storepb.MaskingLevel](tempCTEOuterSchemaInfo.ColumnList[i].MaskingLevel, fieldsInRecursiveClause[i].maskingLevel) {
								change = true
								tempCTEOuterSchemaInfo.ColumnList[i].MaskingLevel = fieldsInRecursiveClause[i].maskingLevel
								result[i].maskingLevel = fieldsInRecursiveClause[i].maskingLevel
							}
						}
					} else if allQueryExpression := queryExpression.AllQuery_expression(); len(allQueryExpression) > 1 {
						fieldsInRecursiveClause, err := extractor.extractTSqlSensitiveFieldsFromQueryExpression(allQueryExpression[len(allQueryExpression)-1])
						if err != nil {
							return nil, errors.Wrapf(err, "failed to extract sensitive fields of the recursive clause of recursive CTE %q near line %d", normalizedCTEName, allQueryExpression[len(allQueryExpression)-1].GetStart().GetLine())
						}
						if len(fieldsInRecursiveClause) != len(tempCTEOuterSchemaInfo.ColumnList) {
							return nil, errors.Wrapf(err, "recursive clause returns %d fields, but anchor clause returns %d fields in recursive CTE %q near line %d", len(fieldsInRecursiveClause), len(tempCTEOuterSchemaInfo.ColumnList), normalizedCTEName, allQueryExpression[len(allQueryExpression)-1].GetStart().GetLine())
						}
						extractor.cteOuterSchemaInfo = extractor.cteOuterSchemaInfo[:originalSize]
						for i := 0; i < len(fieldsInRecursiveClause); i++ {
							if cmp.Less[storepb.MaskingLevel](tempCTEOuterSchemaInfo.ColumnList[i].MaskingLevel, fieldsInRecursiveClause[i].maskingLevel) {
								change = true
								tempCTEOuterSchemaInfo.ColumnList[i].MaskingLevel = fieldsInRecursiveClause[i].maskingLevel
								result[i].maskingLevel = fieldsInRecursiveClause[i].maskingLevel
							}
						}
					}
					if !change {
						break
					}
					originalSize = len(extractor.cteOuterSchemaInfo)
					extractor.cteOuterSchemaInfo = append(extractor.cteOuterSchemaInfo, tempCTEOuterSchemaInfo)
				}
				extractor.cteOuterSchemaInfo = extractor.cteOuterSchemaInfo[:originalSize]
			}
			if v := commonTableExpression.Column_name_list(); v != nil {
				if len(result) != len(v.AllId_()) {
					return nil, errors.Errorf("the number of column name list %d does not match the number of columns %d", len(v.AllId_()), len(result))
				}
				for i, columnName := range v.AllId_() {
					normalizedColumnName := parser.NormalizeTSQLIdentifier(columnName)
					result[i].name = normalizedColumnName
				}
			}
			// Append to the extractor.schemaInfo.DatabaseList
			columnList := make([]db.ColumnInfo, 0, len(result))
			for _, field := range result {
				columnList = append(columnList, db.ColumnInfo{
					Name:         field.name,
					MaskingLevel: field.maskingLevel,
				})
			}
			extractor.cteOuterSchemaInfo = append(extractor.cteOuterSchemaInfo, db.TableSchema{
				Name:       normalizedCTEName,
				ColumnList: columnList,
			})
		}
	}
	return extractor.extractTSqlSensitiveFieldsFromSelectStatement(ctx.Select_statement())
}

// extractTSqlSensitiveFieldsFromSelectStatement extracts sensitive fields from select_statement.
func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromSelectStatement(ctx tsqlparser.ISelect_statementContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	queryResult, err := extractor.extractTSqlSensitiveFieldsFromQueryExpression(ctx.Query_expression())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract sensitive fields from `query_expression` in `select_statement`")
	}

	return queryResult, nil
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromQueryExpression(ctx tsqlparser.IQuery_expressionContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	if ctx.Query_specification() != nil {
		left, err := extractor.extractTSqlSensitiveFieldsFromQuerySpecification(ctx.Query_specification())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `query_specification` in `query_expression`")
		}
		if allSQLUnions := ctx.AllSql_union(); len(allSQLUnions) > 0 {
			for i, sqlUnion := range allSQLUnions {
				// For UNION operator, the number of the columns in the result set is the same, and will use the left part's column name.
				// So we only need to extract the sensitive fields of the right part.
				right, err := extractor.extractTSqlSensitiveFieldsFromQuerySpecification(sqlUnion.Query_specification())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract the %d set operator near line %d", i+1, sqlUnion.GetStart().GetLine())
				}
				if len(left) != len(right) {
					return nil, errors.Wrapf(err, "the number of columns in the query statement nearly line %d returns %d fields, but %d set operator near line %d returns %d fields", ctx.GetStart().GetLine(), len(left), i+1, sqlUnion.GetStart().GetLine(), len(right))
				}
				for i := range right {
					if cmp.Less[storepb.MaskingLevel](left[i].maskingLevel, right[i].maskingLevel) {
						left[i].maskingLevel = right[i].maskingLevel
					}
				}
			}
		}
		return left, nil
	}

	if allQueryExpressions := ctx.AllQuery_expression(); len(allQueryExpressions) > 0 {
		left, err := extractor.extractTSqlSensitiveFieldsFromQueryExpression(allQueryExpressions[0])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `query_specification` in `query_expression`")
		}
		for i := 1; i < len(allQueryExpressions); i++ {
			// For UNION operator, the number of the columns in the result set is the same, and will use the left part's column name.
			// So we only need to extract the sensitive fields of the right part.
			right, err := extractor.extractTSqlSensitiveFieldsFromQueryExpression(allQueryExpressions[i])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to extract the %d set operator near line %d", i+1, allQueryExpressions[i].GetStart().GetLine())
			}
			if len(left) != len(right) {
				return nil, errors.Wrapf(err, "the number of columns in the query statement nearly line %d returns %d fields, but %d set operator near line %d returns %d fields", ctx.GetStart().GetLine(), len(left), i+1, allQueryExpressions[i].GetStart().GetLine(), len(right))
			}
			for i := range right {
				if cmp.Less[storepb.MaskingLevel](left[i].maskingLevel, right[i].maskingLevel) {
					left[i].maskingLevel = right[i].maskingLevel
				}
			}
		}
		return left, nil
	}

	panic("never reach here")
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromQuerySpecification(ctx tsqlparser.IQuery_specificationContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	if from := ctx.GetFrom(); from != nil {
		fromFieldList, err := extractor.extractTSqlSensitiveFieldsFromTableSources(ctx.Table_sources())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_sources` in `query_specification`")
		}
		originalFromFieldList := len(extractor.fromFieldList)
		extractor.fromFieldList = append(extractor.fromFieldList, fromFieldList...)
		defer func() {
			extractor.fromFieldList = extractor.fromFieldList[:originalFromFieldList]
		}()
	}

	var result []fieldInfo

	selectList := ctx.Select_list()
	for _, selectListElem := range selectList.AllSelect_list_elem() {
		if asterisk := selectListElem.Asterisk(); asterisk != nil {
			var normalizedDatabaseName, normalizedSchemaName, normalizedTableName string
			if tableName := asterisk.Table_name(); tableName != nil {
				normalizedDatabaseName, normalizedSchemaName, normalizedTableName = splitTableNameIntoNormalizedParts(tableName)
			}
			left, err := extractor.tsqlGetAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get all fields of table %s.%s.%s", normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
			}
			result = append(result, left...)
		} else if selectListElem.Udt_elem() != nil {
			// TODO(zp): handle the UDT.
			result = append(result, fieldInfo{
				name:         fmt.Sprintf("UNSUPPORTED UDT %s", selectListElem.GetText()),
				maskingLevel: defaultMaskingLevel,
			})
		} else if selectListElem.LOCAL_ID() != nil {
			// TODO(zp): handle the local variable, SELECT @a=id FROM blog.dbo.t1;
			result = append(result, fieldInfo{
				name:         fmt.Sprintf("UNSUPPORTED LOCALID %s", selectListElem.GetText()),
				maskingLevel: defaultMaskingLevel,
			})
		} else if expressionElem := selectListElem.Expression_elem(); expressionElem != nil {
			columnName, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expressionElem)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to check if the expression element is sensitive")
			}
			result = append(result, fieldInfo{
				name:         columnName,
				maskingLevel: maskingLevel,
			})
		}
	}

	return result, nil
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromTableSources(ctx tsqlparser.ITable_sourcesContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	var allTableSources []tsqlparser.ITable_sourceContext
	if v := ctx.Non_ansi_join(); v != nil {
		allTableSources = v.GetSource()
	} else if len(ctx.AllTable_source()) != 0 {
		allTableSources = ctx.GetSource()
	}

	var result []fieldInfo
	// If there are multiple table sources, the default join type is CROSS JOIN.
	for _, tableSource := range allTableSources {
		tableSourceResult, err := extractor.extractTSqlSensitiveFieldsFromTableSource(tableSource)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_source` in `table_sources`")
		}
		result = append(result, tableSourceResult...)
	}
	return result, nil
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromTableSource(ctx tsqlparser.ITable_sourceContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	var result []fieldInfo
	left, err := extractor.extractTSqlSensitiveFieldsFromTableSourceItem(ctx.Table_source_item())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_source_item` in `table_source`")
	}
	result = append(result, left...)

	if allJoinParts := ctx.AllJoin_part(); len(allJoinParts) > 0 {
		for _, joinPart := range allJoinParts {
			if joinOn := joinPart.Join_on(); joinOn != nil {
				right, err := extractor.extractTSqlSensitiveFieldsFromTableSource(joinOn.Table_source())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_source` in `join_on`")
				}
				result = append(result, right...)
			}
			if crossJoin := joinPart.Cross_join(); crossJoin != nil {
				right, err := extractor.extractTSqlSensitiveFieldsFromTableSourceItem(crossJoin.Table_source_item())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_source` in `cross_join`")
				}
				result = append(result, right...)
			}
			if apply := joinPart.Apply_(); apply != nil {
				right, err := extractor.extractTSqlSensitiveFieldsFromTableSourceItem(apply.Table_source_item())
				if err != nil {
					return nil, errors.Wrapf(err, "failed to extract sensitive fields from `table_source` in `apply`")
				}
				result = append(result, right...)
			}
			// TODO(zp): handle pivot and unpivot.
			if pivot := joinPart.Pivot(); pivot != nil {
				return nil, errors.New("pivot is not supported yet")
			}
			if unpivot := joinPart.Unpivot(); unpivot != nil {
				return nil, errors.New("unpivot is not supported yet")
			}
		}
	}

	return result, nil
}

// extractTSqlSensitiveFieldsFromTableSourceItem extracts sensitive fields from table source item.
func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromTableSourceItem(ctx tsqlparser.ITable_source_itemContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	var result []fieldInfo
	var err error
	// TODO(zp): handle other cases likes ROWSET_FUNCTION.
	if ctx.Full_table_name() != nil {
		normalizedDatabaseName, tableSchema, err := extractor.tsqlFindTableSchema(ctx.Full_table_name(), "", extractor.currentDatabase, "dbo")
		if err != nil {
			return nil, err
		}
		for _, column := range tableSchema.ColumnList {
			result = append(result, fieldInfo{
				database:     normalizedDatabaseName,
				table:        tableSchema.Name,
				name:         column.Name,
				maskingLevel: column.MaskingLevel,
			})
		}
	}

	if ctx.Table_source() != nil {
		return extractor.extractTSqlSensitiveFieldsFromTableSource(ctx.Table_source())
	}

	if ctx.Derived_table() != nil {
		result, err = extractor.extractTSqlSensitiveFieldsFromDerivedTable(ctx.Derived_table())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `derived_table` in `table_source_item`")
		}
	}

	// If there are as_table_alias, we should patch the table name to the alias name, and reset the schema and database.
	// For example:
	// SELECT t1.id FROM blog.dbo.t1 AS TT1; -- The multi-part identifier "t1.id" could not be bound.
	// SELECT TT1.id FROM blog.dbo.t1 AS TT1; -- OK
	if asTableAlias := ctx.As_table_alias(); asTableAlias != nil {
		asName := parser.NormalizeTSQLIdentifier(asTableAlias.Table_alias().Id_())
		for i := 0; i < len(result); i++ {
			result[i].table = asName
			result[i].schema = ""
			result[i].database = ""
		}
	}

	if columnAliasList := ctx.Column_alias_list(); columnAliasList != nil {
		allColumnAlias := columnAliasList.AllColumn_alias()
		if len(allColumnAlias) != len(result) {
			return nil, errors.Errorf("the number of column alias %d does not match the number of columns %d", len(allColumnAlias), len(result))
		}
		for i := 0; i < len(result); i++ {
			if allColumnAlias[i].Id_() != nil {
				result[i].name = parser.NormalizeTSQLIdentifier(allColumnAlias[i].Id_())
				continue
			} else if allColumnAlias[i].STRING() != nil {
				result[i].name = allColumnAlias[i].STRING().GetText()
				continue
			}
			panic("never reach here")
		}
	}

	return result, nil
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromDerivedTable(ctx tsqlparser.IDerived_tableContext) ([]fieldInfo, error) {
	if ctx == nil {
		return nil, nil
	}

	allSubquery := ctx.AllSubquery()
	if len(allSubquery) > 0 {
		left, err := extractor.extractTSqlSensitiveFieldsFromSubquery(allSubquery[0])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to extract sensitive fields from `subquery` in `derived_table`")
		}
		for i := 1; i < len(allSubquery); i++ {
			// For UNION operator, the number of the columns in the result set is the same, and will use the left part's column name.
			// So we only need to extract the sensitive fields of the right part.
			right, err := extractor.extractTSqlSensitiveFieldsFromSubquery(allSubquery[i])
			if err != nil {
				return nil, errors.Wrapf(err, "failed to extract the %d set operator near line %d", i+1, allSubquery[i].GetStart().GetLine())
			}
			if len(left) != len(right) {
				return nil, errors.Wrapf(err, "the number of columns in the derived table statement nearly line %d returns %d fields, but %d set operator near line %d returns %d fields", ctx.GetStart().GetLine(), len(left), i+1, allSubquery[i].GetStart().GetLine(), len(right))
			}
			for i := range right {
				if cmp.Less[storepb.MaskingLevel](left[i].maskingLevel, right[i].maskingLevel) {
					left[i].maskingLevel = right[i].maskingLevel
				}
			}
		}
		return left, nil
	}

	if tableValueConstructor := ctx.Table_value_constructor(); tableValueConstructor != nil {
		return extractor.extractTSqlSensitiveFieldsFromTableValueConstructor(tableValueConstructor)
	}

	panic("never reach here")
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromTableValueConstructor(ctx tsqlparser.ITable_value_constructorContext) ([]fieldInfo, error) {
	if allExpressionList := ctx.AllExpression_list_(); len(allExpressionList) > 0 {
		// The number of expression in each expression list should be the same.
		// But we do not check, just use the first one, and engine will throw a compilation error if the number of expressions are not the same.
		expressionList := allExpressionList[0]
		var result []fieldInfo
		for _, expression := range expressionList.AllExpression() {
			columnName, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to check if the expression is sensitive")
			}
			result = append(result, fieldInfo{
				name:         columnName,
				maskingLevel: maskingLevel,
			})
		}
		return result, nil
	}
	panic("never reach here")
}

func (extractor *sensitiveFieldExtractor) extractTSqlSensitiveFieldsFromSubquery(ctx tsqlparser.ISubqueryContext) ([]fieldInfo, error) {
	return extractor.extractTSqlSensitiveFieldsFromSelectStatement(ctx.Select_statement())
}

func (extractor *sensitiveFieldExtractor) tsqlFindTableSchema(fullTableName tsqlparser.IFull_table_nameContext, normalizedFallbackLinkedServerName, normalizedFallbackDatabaseName, normalizedFallbackSchemaName string) (string, db.TableSchema, error) {
	normalizedLinkedServer, normalizedDatabaseName, normalizedSchemaName, normalizedTableName := normalizeFullTableName(fullTableName, "", "", "")
	if normalizedLinkedServer != "" {
		// TODO(zp): How do we handle the linked server?
		return "", db.TableSchema{}, errors.Errorf("linked server is not supported yet, but found %q", fullTableName.GetText())
	}
	// For tsql, we should find the table schema in cteOuterSchemaInfo by ascending order.
	if normalizedDatabaseName == "" && normalizedSchemaName == "" {
		for _, tableSchema := range extractor.cteOuterSchemaInfo {
			if extractor.isIdentifierEqual(normalizedTableName, tableSchema.Name) {
				return normalizedDatabaseName, tableSchema, nil
			}
		}
	}
	normalizedLinkedServer, normalizedDatabaseName, normalizedSchemaName, normalizedTableName = normalizeFullTableName(fullTableName, normalizedFallbackLinkedServerName, normalizedFallbackDatabaseName, normalizedFallbackSchemaName)
	if normalizedLinkedServer != "" {
		// TODO(zp): How do we handle the linked server?
		return "", db.TableSchema{}, errors.Errorf("linked server is not supported yet, but found %q", fullTableName.GetText())
	}
	for _, databaseSchema := range extractor.schemaInfo.DatabaseList {
		if normalizedDatabaseName != "" && !extractor.isIdentifierEqual(normalizedDatabaseName, databaseSchema.Name) {
			continue
		}
		for _, schemaSchema := range databaseSchema.SchemaList {
			if normalizedSchemaName != "" && !extractor.isIdentifierEqual(normalizedSchemaName, schemaSchema.Name) {
				continue
			}
			for _, tableSchema := range schemaSchema.TableList {
				if !extractor.isIdentifierEqual(normalizedTableName, tableSchema.Name) {
					continue
				}
				return normalizedDatabaseName, tableSchema, nil
			}
		}
	}
	return "", db.TableSchema{}, errors.Errorf("table %s.%s.%s is not found", normalizedDatabaseName, normalizedSchemaName, normalizedTableName)
}

// splitTableNameIntoNormalizedParts splits the table name into normalized 3 parts: database, schema, table.
func splitTableNameIntoNormalizedParts(tableName tsqlparser.ITable_nameContext) (string, string, string) {
	var database string
	if d := tableName.GetDatabase(); d != nil {
		normalizedD := parser.NormalizeTSQLIdentifier(d)
		if normalizedD != "" {
			database = normalizedD
		}
	}

	var schema string
	if s := tableName.GetSchema(); s != nil {
		normalizedS := parser.NormalizeTSQLIdentifier(s)
		if normalizedS != "" {
			schema = normalizedS
		}
	}

	var table string
	if t := tableName.GetTable(); t != nil {
		normalizedT := parser.NormalizeTSQLIdentifier(t)
		if normalizedT != "" {
			table = normalizedT
		}
	}
	return database, schema, table
}

// normalizeFullTableName normalizes the each part of the full table name, returns (linkedServer, database, schema, table).
func normalizeFullTableName(fullTableName tsqlparser.IFull_table_nameContext, normalizedFallbackLinkedServerName, normalizedFallbackDatabaseName, normalizedFallbackSchemaName string) (string, string, string, string) {
	if fullTableName == nil {
		return "", "", "", ""
	}
	// TODO(zp): unify here and the related code in sql_service.go
	linkedServer := normalizedFallbackLinkedServerName
	if server := fullTableName.GetLinkedServer(); server != nil {
		linkedServer = parser.NormalizeTSQLIdentifier(server)
	}

	database := normalizedFallbackDatabaseName
	if d := fullTableName.GetDatabase(); d != nil {
		normalizedD := parser.NormalizeTSQLIdentifier(d)
		if normalizedD != "" {
			database = normalizedD
		}
	}

	schema := normalizedFallbackSchemaName
	if s := fullTableName.GetSchema(); s != nil {
		normalizedS := parser.NormalizeTSQLIdentifier(s)
		if normalizedS != "" {
			schema = normalizedS
		}
	}

	var table string
	if t := fullTableName.GetTable(); t != nil {
		normalizedT := parser.NormalizeTSQLIdentifier(t)
		if normalizedT != "" {
			table = normalizedT
		}
	}

	return linkedServer, database, schema, table
}

func (extractor *sensitiveFieldExtractor) tsqlGetAllFieldsOfTableInFromOrOuterCTE(normalizedDatabaseName, normalizedSchemaName, normalizedTableName string) ([]fieldInfo, error) {
	type maskType = uint8
	const (
		maskNone         maskType = 0
		maskDatabaseName maskType = 1 << iota
		maskSchemaName
		maskTableName
	)
	mask := maskNone
	if normalizedTableName != "" {
		mask |= maskTableName
	}
	if normalizedSchemaName != "" {
		if mask&maskTableName == 0 {
			return nil, errors.Errorf(`table name %s is specified without column name`, normalizedTableName)
		}
		mask |= maskSchemaName
	}
	if normalizedDatabaseName != "" {
		if mask&maskSchemaName == 0 {
			return nil, errors.Errorf(`database name %s is specified without schema name`, normalizedDatabaseName)
		}
		mask |= maskDatabaseName
	}

	var result []fieldInfo
	for _, field := range extractor.fromFieldList {
		if mask&maskDatabaseName != 0 && !extractor.isIdentifierEqual(normalizedDatabaseName, field.database) {
			continue
		}
		if mask&maskSchemaName != 0 && !extractor.isIdentifierEqual(normalizedSchemaName, field.schema) {
			continue
		}
		if mask&maskTableName != 0 && !extractor.isIdentifierEqual(normalizedTableName, field.table) {
			continue
		}
		result = append(result, field)
	}
	return result, nil
}

func (extractor *sensitiveFieldExtractor) tsqlIsFullColumnNameSensitive(ctx tsqlparser.IFull_column_nameContext) (fieldInfo, error) {
	normalizedLinkedServer, normalizedDatabaseName, normalizedSchemaName, normalizedTableName := normalizeFullTableName(ctx.Full_table_name(), "", "", "")
	if normalizedLinkedServer != "" {
		return fieldInfo{}, errors.Errorf("linked server is not supported yet, but found %q", ctx.GetText())
	}
	normalizedColumnName := parser.NormalizeTSQLIdentifier(ctx.Id_())

	return extractor.tsqlIsFieldSensitive(normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName)
}

func (extractor *sensitiveFieldExtractor) tsqlIsFieldSensitive(normalizedDatabaseName string, normalizedSchemaName string, normalizedTableName string, normalizedColumnName string) (fieldInfo, error) {
	type maskType = uint8
	const (
		maskNone         maskType = 0
		maskDatabaseName maskType = 1 << iota
		maskSchemaName
		maskTableName
		maskColumnName
	)
	mask := maskNone
	if normalizedColumnName != "" {
		mask |= maskColumnName
	}
	if normalizedTableName != "" {
		if mask&maskColumnName == 0 {
			return fieldInfo{}, errors.Errorf(`table name %s is specified without column name`, normalizedTableName)
		}
		mask |= maskTableName
	}
	if normalizedSchemaName != "" {
		if mask&maskTableName == 0 {
			return fieldInfo{}, errors.Errorf(`schema name %s is specified without table name`, normalizedSchemaName)
		}
		mask |= maskSchemaName
	}
	if normalizedDatabaseName != "" {
		if mask&maskSchemaName == 0 {
			return fieldInfo{}, errors.Errorf(`database name %s is specified without schema name`, normalizedDatabaseName)
		}
		mask |= maskDatabaseName
	}

	if mask == maskNone {
		return fieldInfo{}, errors.Errorf(`no object name is specified`)
	}

	// We just need to iterate through the fromFieldList sequentially until we find the first matching object.

	// It is safe if there are two or more objects in the fromFieldList have the same column name, because the executor
	// will throw a compilation error if the column name is ambiguous.
	// For example, there are two tables T1 and T2, and both of them have a column named "C1". The following query will throw
	// a compilation error:
	// SELECT C1 FROM T1, T2;
	//
	// But users can specify the table name to avoid the compilation error:
	// SELECT T1.C1 FROM T1, T2;
	//
	// Further more, users can not use the original table name if they specify the alias name:
	// SELECT T1.C1 FROM T1 AS T3, T2; -- invalid identifier 'ADDRESS.ID'
	for _, field := range extractor.fromFieldList {
		if mask&maskDatabaseName != 0 && !extractor.isIdentifierEqual(normalizedDatabaseName, field.database) {
			continue
		}
		if mask&maskSchemaName != 0 && !extractor.isIdentifierEqual(normalizedSchemaName, field.schema) {
			continue
		}
		if mask&maskTableName != 0 && !extractor.isIdentifierEqual(normalizedTableName, field.table) {
			continue
		}
		if mask&maskColumnName != 0 && !extractor.isIdentifierEqual(normalizedColumnName, field.name) {
			continue
		}
		return field, nil
	}
	return fieldInfo{}, errors.Errorf(`no matching column %q.%q.%q.%q`, normalizedDatabaseName, normalizedSchemaName, normalizedTableName, normalizedColumnName)
}

// isIdentifierEqual compares the identifier with the given normalized parts, returns true if they are equal.
// It will consider the case sensitivity based on the current database.
func (extractor *sensitiveFieldExtractor) isIdentifierEqual(a, b string) bool {
	if !extractor.schemaInfo.IgnoreCaseSensitive {
		return a == b
	}
	if len(a) != len(b) {
		return false
	}
	runeA, runeB := []rune(a), []rune(b)
	for i := 0; i < len(runeA); i++ {
		if unicode.ToLower(runeA[i]) != unicode.ToLower(runeB[i]) {
			return false
		}
	}
	return true
}

// evalExpressionElemMaskingLevel returns true if the expression element is sensitive, and returns the column name.
// It is the closure of the expression_elemContext, it will recursively check the sub expression element.
func (extractor *sensitiveFieldExtractor) evalExpressionElemMaskingLevel(ctx antlr.RuleContext) (string, storepb.MaskingLevel, error) {
	if ctx == nil {
		return "", defaultMaskingLevel, nil
	}
	switch ctx := ctx.(type) {
	case *tsqlparser.Expression_elemContext:
		columName, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the expression element is sensitive")
		}
		if columnAlias := ctx.Column_alias(); columnAlias != nil {
			columName = parser.NormalizeTSQLIdentifier(columnAlias.Id_())
		} else if asColumnAlias := ctx.As_column_alias(); asColumnAlias != nil {
			columName = parser.NormalizeTSQLIdentifier(asColumnAlias.Column_alias().Id_())
		}
		return columName, maskingLevel, nil
	case *tsqlparser.ExpressionContext:
		if ctx.Primitive_expression() != nil {
			return extractor.evalExpressionElemMaskingLevel(ctx.Primitive_expression())
		}
		if ctx.Function_call() != nil {
			return extractor.evalExpressionElemMaskingLevel(ctx.Function_call())
		}
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the expression is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if valueCall := ctx.Value_call(); valueCall != nil {
			return extractor.evalExpressionElemMaskingLevel(valueCall)
		}
		if queryCall := ctx.Query_call(); queryCall != nil {
			return extractor.evalExpressionElemMaskingLevel(queryCall)
		}
		if existCall := ctx.Exist_call(); existCall != nil {
			return extractor.evalExpressionElemMaskingLevel(existCall)
		}
		if modifyCall := ctx.Modify_call(); modifyCall != nil {
			return extractor.evalExpressionElemMaskingLevel(modifyCall)
		}
		if hierarchyIDCall := ctx.Hierarchyid_call(); hierarchyIDCall != nil {
			return extractor.evalExpressionElemMaskingLevel(hierarchyIDCall)
		}
		if caseExpression := ctx.Case_expression(); caseExpression != nil {
			return extractor.evalExpressionElemMaskingLevel(caseExpression)
		}
		if fullColumnName := ctx.Full_column_name(); fullColumnName != nil {
			return extractor.evalExpressionElemMaskingLevel(fullColumnName)
		}
		if bracketExpression := ctx.Bracket_expression(); bracketExpression != nil {
			return extractor.evalExpressionElemMaskingLevel(bracketExpression)
		}
		if unaryOperationExpression := ctx.Unary_operator_expression(); unaryOperationExpression != nil {
			return extractor.evalExpressionElemMaskingLevel(unaryOperationExpression)
		}
		if overClause := ctx.Over_clause(); overClause != nil {
			return extractor.evalExpressionElemMaskingLevel(overClause)
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Unary_operator_expressionContext:
		if expression := ctx.Expression(); expression != nil {
			return extractor.evalExpressionElemMaskingLevel(expression)
		}
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Bracket_expressionContext:
		if expression := ctx.Expression(); expression != nil {
			return extractor.evalExpressionElemMaskingLevel(expression)
		}
		if subquery := ctx.Subquery(); subquery != nil {
			return extractor.evalExpressionElemMaskingLevel(subquery)
		}
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Case_expressionContext:
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the case_expression is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if allSwitchSections := ctx.AllSwitch_section(); len(allSwitchSections) > 0 {
			for _, switchSection := range allSwitchSections {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(switchSection)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the case_expression is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if allSwitchSearchConditionSections := ctx.AllSwitch_search_condition_section(); len(allSwitchSearchConditionSections) > 0 {
			for _, switchSearchConditionSection := range allSwitchSearchConditionSections {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(switchSearchConditionSection)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the case_expression is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Switch_sectionContext:
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the switch_setion is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Switch_search_condition_sectionContext:
		finalLevel := defaultMaskingLevel
		if searchCondition := ctx.Search_condition(); searchCondition != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(searchCondition)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the switch_search_condition_section is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if expression := ctx.Expression(); expression != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the switch_search_condition_section is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Search_conditionContext:
		if predicate := ctx.Predicate(); predicate != nil {
			return extractor.evalExpressionElemMaskingLevel(predicate)
		}
		finalLevel := defaultMaskingLevel
		if allSearchConditions := ctx.AllSearch_condition(); len(allSearchConditions) > 0 {
			for _, searchCondition := range allSearchConditions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(searchCondition)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the search_condition is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.PredicateContext:
		if subquery := ctx.Subquery(); subquery != nil {
			return extractor.evalExpressionElemMaskingLevel(subquery)
		}
		if freeTextPredicate := ctx.Freetext_predicate(); freeTextPredicate != nil {
			return extractor.evalExpressionElemMaskingLevel(freeTextPredicate)
		}

		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the predicate is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expressionList)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the predicate is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Freetext_predicateContext:
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the freetext_predicate is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if allCullColumnName := ctx.AllFull_column_name(); len(allCullColumnName) > 0 {
			for _, fullColumnName := range allCullColumnName {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(fullColumnName)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the freetext_predicate is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.SubqueryContext:
		// For subquery, we clone the current extractor, reset the from list, but keep the cte, and then extract the sensitive fields from the subquery
		cloneExtractor := &sensitiveFieldExtractor{
			currentDatabase:    extractor.currentDatabase,
			schemaInfo:         extractor.schemaInfo,
			cteOuterSchemaInfo: extractor.cteOuterSchemaInfo,
		}
		fieldInfo, err := cloneExtractor.extractTSqlSensitiveFieldsFromSubquery(ctx)
		// The expect behavior is the fieldInfo contains only one field, which is the column name,
		// but in order to do not block user, we just return isSensitive if there is any sensitive field.
		// return fieldInfo[0].sensitive, err
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the subquery is sensitive")
		}
		finalLevel := defaultMaskingLevel
		for _, field := range fieldInfo {
			if cmp.Less[storepb.MaskingLevel](finalLevel, field.maskingLevel) {
				finalLevel = field.maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Hierarchyid_callContext:
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the hierarchyid_call is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Query_callContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Exist_callContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Modify_callContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Value_callContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Primitive_expressionContext:
		if ctx.Primitive_constant() != nil {
			_, sensitive, err := extractor.evalExpressionElemMaskingLevel(ctx.Primitive_constant())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the primitive constant is sensitive")
			}
			return ctx.GetText(), sensitive, nil
		}
		panic("never reach here")
	case *tsqlparser.Primitive_constantContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Function_callContext:
		// In TSqlParser.g4, the function_callContext is defined as:
		// 	function_call
		// : ranking_windowed_function                         #RANKING_WINDOWED_FUNC
		// | aggregate_windowed_function                       #AGGREGATE_WINDOWED_FUNC
		// ...
		// ;
		// So it will be parsed as RANKING_WINDOWED_FUNC, AGGREGATE_WINDOWED_FUNC, etc.
		// We just need to check the first token to see if it is a sensitive function.
		panic("never reach here")
	case *tsqlparser.RANKING_WINDOWED_FUNCContext:
		return extractor.evalExpressionElemMaskingLevel(ctx.Ranking_windowed_function())
	case *tsqlparser.Ranking_windowed_functionContext:
		finalLevel := defaultMaskingLevel
		if overClause := ctx.Over_clause(); overClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(overClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the ranking_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if expression := ctx.Expression(); expression != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the ranking_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Over_clauseContext:
		finalLevel := defaultMaskingLevel
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the over_clause is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if orderByClause := ctx.Order_by_clause(); orderByClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(orderByClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the over_clause is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if rowOrRangeClause := ctx.Row_or_range_clause(); rowOrRangeClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(rowOrRangeClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the over_clause is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Expression_list_Context:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the expression_list is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Order_by_clauseContext:
		finalLevel := defaultMaskingLevel
		for _, orderByExpression := range ctx.GetOrder_bys() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(orderByExpression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the order_by_clause is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Order_by_expressionContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the order_by_expression is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.Row_or_range_clauseContext:
		if windowFrameExtent := ctx.Window_frame_extent(); windowFrameExtent != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(windowFrameExtent)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the row_or_range_clause is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		}
		panic("never reach here")
	case *tsqlparser.Window_frame_extentContext:
		if windowFramePreceding := ctx.Window_frame_preceding(); windowFramePreceding != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(windowFramePreceding)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the window_frame_extent is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		}
		if windowFrameBounds := ctx.AllWindow_frame_bound(); len(windowFrameBounds) > 0 {
			finalLevel := defaultMaskingLevel
			for _, windowFrameBound := range windowFrameBounds {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(windowFrameBound)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the window_frame_extent is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		panic("never reach here")
	case *tsqlparser.Window_frame_boundContext:
		if preceding := ctx.Window_frame_preceding(); preceding != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(preceding)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the window_frame_bound is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		} else if following := ctx.Window_frame_following(); following != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(following)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the window_frame_bound is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		}
		panic("never reach here")
	case *tsqlparser.Window_frame_precedingContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Window_frame_followingContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.AGGREGATE_WINDOWED_FUNCContext:
		return extractor.evalExpressionElemMaskingLevel(ctx.Aggregate_windowed_function())
	case *tsqlparser.Aggregate_windowed_functionContext:
		finalLevel := defaultMaskingLevel
		if allDistinctExpression := ctx.All_distinct_expression(); allDistinctExpression != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(allDistinctExpression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the aggregate_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if overClause := ctx.Over_clause(); overClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(overClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the aggregate_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if expression := ctx.Expression(); expression != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the aggregate_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expressionList)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the aggregate_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.All_distinct_expressionContext:
		_, sensitive, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the all_distinct_expression is sensitive")
		}
		return ctx.GetText(), sensitive, nil
	case *tsqlparser.ANALYTIC_WINDOWED_FUNCContext:
		return extractor.evalExpressionElemMaskingLevel(ctx.Analytic_windowed_function())
	case *tsqlparser.Analytic_windowed_functionContext:
		finalLevel := defaultMaskingLevel
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the analytic_windowed_function is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if overClause := ctx.Over_clause(); overClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(overClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the analytic_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expressionList)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the analytic_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if orderByClause := ctx.Order_by_clause(); orderByClause != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(orderByClause)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the analytic_windowed_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.BUILT_IN_FUNCContext:
		return extractor.evalExpressionElemMaskingLevel(ctx.Built_in_functions())
	case *tsqlparser.APP_NAMEContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.APPLOCK_MODEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the applock_mode is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.APPLOCK_TESTContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the applock_test is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ASSEMBLYPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the assemblyproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.COL_LENGTHContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the col_length is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.COL_NAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the col_name is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.COLUMNPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the columnproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATABASEPROPERTYEXContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the databasepropertyex is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DB_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the db_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DB_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the db_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILE_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the file_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILE_IDEXContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the file_idex is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILE_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the file_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILEGROUP_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the filegroup_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILEGROUP_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the filegroup_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FILEGROUPPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the filegroupproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.FILEPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the fileproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.FILEPROPERTYEXContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the filepropertyex is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.FULLTEXTCATALOGPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the fulltextcatalogproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.FULLTEXTSERVICEPROPERTYContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the fulltextserviceproperty is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.INDEX_COLContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the index_col is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.INDEXKEY_PROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the indexkey_property is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.INDEXPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the indexproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.OBJECT_DEFINITIONContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the object_definition is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.OBJECT_IDContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the object_id is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.OBJECT_NAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the object_name is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.OBJECT_SCHEMA_NAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the object_schema_name is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.OBJECTPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the objectproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.OBJECTPROPERTYEXContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the objectpropertyex is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.PARSENAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the parsename is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.SCHEMA_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the schema_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SCHEMA_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the schema_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SERVERPROPERTYContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the serverproperty is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.STATS_DATEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the stats_date is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.TYPE_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the type_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.TYPE_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the type_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.TYPEPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the typeproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ASCIIContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the ascii is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CHARContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the char is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CHARINDEXContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the charindex is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.CONCATContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the concat is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.CONCAT_WSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the concat_ws is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DIFFERENCEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the difference is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.FORMATContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the format is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.LEFTContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the left is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.LENContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the len is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.LOWERContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the lower is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.LTRIMContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the ltrim is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.NCHARContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the nchar is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.PATINDEXContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the patindex is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.QUOTENAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the quotename is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.REPLACEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the replace is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.REPLICATEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the replicate is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.REVERSEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the reverse is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.RIGHTContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the right is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.RTRIMContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the rtrim is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SOUNDEXContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the soundex is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SPACEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the space is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.STRContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the str is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.STRINGAGGContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the stringagg is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.STRING_ESCAPEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the string_escape is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.STUFFContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the stuff is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.SUBSTRINGContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the substring is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.TRANSLATEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the translate is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.TRIMContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the trim is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.UNICODEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the unicode is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.UPPERContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the upper is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.BINARY_CHECKSUMContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the binary_checksum is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.CHECKSUMContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the checksum is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.COMPRESSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the compress is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DECOMPRESSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the decompress is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FORMATMESSAGEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the formatmessage is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ISNULLContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the isnull is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ISNUMERICContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the isnumeric is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CASTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the cast is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.TRY_CASTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the try_cast is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CONVERTContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the convert is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.COALESCEContext:
		finalLevel := defaultMaskingLevel
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the coalesce is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.CURSOR_STATUSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the cursor_status is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CERT_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the cert_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DATALENGTHContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datalength is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.IDENT_CURRENTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the  ident_current is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.IDENT_INCRContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the  ident_incr is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.IDENT_SEEDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the  ident_seed is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SQL_VARIANT_PROPERTYContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the sql_variant_property is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DATE_BUCKETContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the date_bucket is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATEADDContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the dateadd is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATEDIFFContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datediff is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATEDIFF_BIGContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datediff_big is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATEFROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datefromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATENAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datename is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DATEPARTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datepart is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DATETIME2FROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datetime2fromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATETIMEFROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datetimefromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATETIMEOFFSETFROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datetimeoffsetfromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATETRUNCContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the datetrunc is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DAYContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the day is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.EOMONTHContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the eomonth is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ISDATEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the isdate is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.MONTHContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the month is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SMALLDATETIMEFROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the smalldatetimefromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.SWITCHOFFSETContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the switchoffset is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.TIMEFROMPARTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the timefromparts is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.TODATETIMEOFFSETContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the todatetimeoffset is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.YEARContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the year is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.NULLIFContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the nullif is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.PARSEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the parse is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.IIFContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the iif is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ISJSONContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the isjson is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.JSON_ARRAYContext:
		finalLevel := defaultMaskingLevel
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the json_array is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.JSON_VALUEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the json_value is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.JSON_QUERYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the json_query is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.JSON_MODIFYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the json_modify is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.JSON_PATH_EXISTSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the json_path_exists is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.ABSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the abs is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.ACOSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the acos is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.ASINContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the asin is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.ATANContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the atan is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.ATN2Context:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the atn2 is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.CEILINGContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the ceiling is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.COSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the cos is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.COTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the cot is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.DEGREESContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the degrees is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.EXPContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the exp is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.FLOORContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the floor is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.LOGContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the log is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.LOG10Context:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the log10 is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.POWERContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the power is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.RADIANSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the radians is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.RANDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the rand is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.ROUNDContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the round is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.MATH_SIGNContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the math_sign is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SINContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the sin is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SQRTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the sqrt is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SQUAREContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the square is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.TANContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the tan is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.GREATESTContext:
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the greatest is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		}
		panic("never reach here")
	case *tsqlparser.LEASTContext:
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the least is sensitive")
			}
			return ctx.GetText(), maskingLevel, nil
		}
		panic("never reach here")
	case *tsqlparser.CERTENCODEDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the certencoded is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.CERTPRIVATEKEYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the certprivatekey is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.DATABASE_PRINCIPAL_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the database_principal_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.HAS_DBACCESSContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the has_dbaccess is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.HAS_PERMS_BY_NAMEContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the has_perms_by_name is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.IS_MEMBERContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the is_member is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.IS_ROLEMEMBERContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the is_rolemember is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.IS_SRVROLEMEMBERContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the is_srvrolemember is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.LOGINPROPERTYContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the loginproperty is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.PERMISSIONSContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the permissions is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.PWDENCRYPTContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the pwdencrypt is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.PWDCOMPAREContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the pwdcompare is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.SESSIONPROPERTYContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the sessionproperty is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SUSER_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the suser_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SUSER_SNAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the suser_sname is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SUSER_SIDContext:
		finalLevel := defaultMaskingLevel
		for _, expression := range ctx.AllExpression() {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the suser_sid is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.USER_IDContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the user_id is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.USER_NAMEContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the user_name is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.SCALAR_FUNCTIONContext:
		finalLevel := defaultMaskingLevel
		if expressionList := ctx.Expression_list_(); expressionList != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Expression_list_())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the scalar_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		if scalarFunctionName := ctx.Scalar_function_name(); scalarFunctionName != nil {
			_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Scalar_function_name())
			if err != nil {
				return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the scalar_function is sensitive")
			}
			if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
				finalLevel = maskingLevel
			}
			if finalLevel == maxMaskingLevel {
				return ctx.GetText(), finalLevel, nil
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Scalar_function_nameContext:
		return ctx.GetText(), defaultMaskingLevel, nil
	case *tsqlparser.Freetext_functionContext:
		finalLevel := defaultMaskingLevel
		if allFullColumnName := ctx.AllFull_column_name(); len(allFullColumnName) > 0 {
			for _, fullColumnName := range allFullColumnName {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(fullColumnName)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the freetext_function is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		if allExpressions := ctx.AllExpression(); len(allExpressions) > 0 {
			finalLevel := defaultMaskingLevel
			for _, expression := range allExpressions {
				_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(expression)
				if err != nil {
					return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the freetext_function is sensitive")
				}
				if cmp.Less[storepb.MaskingLevel](finalLevel, maskingLevel) {
					finalLevel = maskingLevel
				}
				if finalLevel == maxMaskingLevel {
					return ctx.GetText(), finalLevel, nil
				}
			}
		}
		return ctx.GetText(), finalLevel, nil
	case *tsqlparser.Full_column_nameContext:
		fieldInfo, err := extractor.tsqlIsFullColumnNameSensitive(ctx)
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the full_column_name is sensitive")
		}
		return fieldInfo.name, fieldInfo.maskingLevel, nil
	case *tsqlparser.PARTITION_FUNCContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Partition_function().Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the partition_function is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	case *tsqlparser.HIERARCHYID_METHODContext:
		_, maskingLevel, err := extractor.evalExpressionElemMaskingLevel(ctx.Hierarchyid_static_method().Expression())
		if err != nil {
			return "", storepb.MaskingLevel_MASKING_LEVEL_UNSPECIFIED, errors.Wrapf(err, "failed to check if the hierarchyid_method is sensitive")
		}
		return ctx.GetText(), maskingLevel, nil
	}
	panic("never reach here")
}
