package rdbms

import vs "netspective/valuesrc"

type SqlStatementParameter struct {
	sqlStmtParamValueSrc vs.ValueSource
	sqlStmtParamDefaultValueSrc vs.ValueSource
}

type SqlStatementParameters struct {
	sqlStatementParametersList []SqlStatementParameter
}

func NewSqlStatementParameter(value vs.ValueSource, defaultValue vs.ValueSource) (*SqlStatementParameter, error) {
	instance := new(SqlStatementParameter)
	instance.sqlStmtParamValueSrc = value
	instance.sqlStmtParamDefaultValueSrc = defaultValue
	return instance, nil
}

func NewSqlStatementParametersSimple(specs []string) (*SqlStatementParameters, error) {
	instance := new(SqlStatementParameters)

	for s := range specs {
		spec := specs[s]
		vs, err := vs.Factory.CreateValueSource(spec, nil)
		if(err != nil) {
			return instance, err
		}

		param, paramErr := NewSqlStatementParameter(vs, nil)
		if(paramErr != nil) {
			return instance, paramErr
		}
		instance.sqlStatementParametersList = append(instance.sqlStatementParametersList, *param)
	}

	return instance, nil
}
