package pagination

import "fmt"

type Filter struct {
	ColumnKey string      `json:"columnKey"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

type Filters struct {
	And []Filter `json:"and"`
	Or  []Filter `json:"or"`
}

// ValidateSortField checks if the sort field is allowed
func ValidateSortField(field string, allowedFields []string) bool {
	if field == "" {
		return true
	}

	for _, allowed := range allowedFields {
		if field == allowed {
			return true
		}
	}
	return false
}

func BuildWhereClause(filter Filter) (string, interface{}) {
	column := filter.ColumnKey
	operator := filter.Operator
	value := filter.Value

	switch operator {
	case "=", "!=", "<", ">", "<=", ">=":
		return fmt.Sprintf("%s %s ?", column, operator), value
	case "LIKE %_%":
		return fmt.Sprintf("%s LIKE ?", column), "%" + value.(string) + "%"
	case "LIKE _%":
		return fmt.Sprintf("%s LIKE ?", column), value.(string) + "%"
	case "LIKE %_":
		return fmt.Sprintf("%s LIKE ?", column), "%" + value.(string)
	case "IN":
		return fmt.Sprintf("%s IN (?)", column), value
	case "NOT IN":
		return fmt.Sprintf("%s NOT IN (?)", column), value
	case "IS NULL":
		return fmt.Sprintf("%s IS NULL", column), nil
	case "IS NOT NULL":
		return fmt.Sprintf("%s IS NOT NULL", column), nil
	case "BETWEEN":
		if values, ok := value.([]interface{}); ok && len(values) == 2 {
			return fmt.Sprintf("%s BETWEEN ? AND ?", column), values
		}
	}
	return "", nil
}
