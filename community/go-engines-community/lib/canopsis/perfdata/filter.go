package perfdata

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"

func Parse(perfData []string) []utils.RegexExpression {
	perfDataRe := make([]utils.RegexExpression, len(perfData))
	for i, reStr := range perfData {
		perfDataRe[i], _ = utils.NewRegexExpression(reStr)
	}

	return perfDataRe
}

func Filter(entityPerfData, perfData []string, perfDataRe []utils.RegexExpression) []string {
	filteredPerfData := make([]string, 0)

	for _, entityMetric := range entityPerfData {
		matched := false
		for k, metric := range perfData {
			if re := perfDataRe[k]; re != nil {
				matched = re.Match([]byte(entityMetric))
			} else {
				matched = entityMetric == metric
			}

			if matched {
				break
			}
		}

		if matched {
			filteredPerfData = append(filteredPerfData, entityMetric)
		}
	}

	return filteredPerfData
}
