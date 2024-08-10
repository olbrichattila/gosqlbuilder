package builder

import "strings"

func builderConcat(strBuilder *strings.Builder, pars ...string) {
	for _, par := range pars {
		strBuilder.WriteString(par)
	}
}

func validateRelation(relation string) bool {
	switch relation {
	case "=", ">", "<", ">=", "<=", "<>", "!=":
		return true
	default:
		return false
	}
}
