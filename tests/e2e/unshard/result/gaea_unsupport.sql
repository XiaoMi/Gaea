sql:SELECT 5 DIV 2, -5 DIV 2, 5 DIV -2, -5 DIV -2;
mysqlRes:[[2 -2 -2 2]]
gaeaError:Error 1054: Unknown column '5DIV2' in 'field list'

sql:SELECT INSERT('Quadratic', 3, 4, 'What');
mysqlRes:[[QuWhattic]]
gaeaError:Error 1305: FUNCTION sbtest1.INSERT_FUNC does not exist

