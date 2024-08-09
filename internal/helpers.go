package builder

import "strings"

func builderConcat(strBuilder *strings.Builder, pars ...string) {
	for _, par := range pars {
		strBuilder.WriteString(par)
	}
}
