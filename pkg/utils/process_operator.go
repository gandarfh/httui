package utils

func ProcessParamsOperators(values []map[string]string, workspaceId uint) []map[string]string {
	for i, item := range values {
		for k, v := range item {
			values[i][k] = ReplaceByOperator(v, workspaceId)
		}
	}

	return values
}
