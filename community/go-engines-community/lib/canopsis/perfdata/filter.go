package perfdata

import "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/utils"

func Parse(perfData []string) []utils.RegexExpression {
	perfDataRe := make([]utils.RegexExpression, len(perfData))
	for i, reStr := range perfData {
		perfDataRe[i], _ = utils.NewRegexExpression(reStr)
	}

	return perfDataRe
}

func Filter(perfData []string, perfDataRe []utils.RegexExpression, entityPerfData []string) []string {
	filteredPerfData := make([]string, 0)

	for i, metric := range perfData {
		re := perfDataRe[i]
		matched := false
		for _, entityMetric := range entityPerfData {
			if re == nil {
				matched = entityMetric == metric
			} else {
				matched = re.Match([]byte(entityMetric))
			}

			if matched {
				break
			}
		}

		if matched {
			filteredPerfData = append(filteredPerfData, metric)
		}
	}

	return filteredPerfData
}
