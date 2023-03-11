package validate

import (
	"regexp"
	"strings"
)

var (
	reHeader = regexp.MustCompile(`^\s@(\w+):\s*$`)
	reOption = regexp.MustCompile(`^\s{2,}(\w+):\s*(.+)$`)
)

type Option struct{ K, V string }

func ParseDoc(comments []string) map[string][]Option {
	annos := make(map[string][]Option)

	var headerName string
	for _, comment := range comments {
		c := strings.TrimPrefix(comment, "//")

		result := reHeader.FindAllStringSubmatch(c, -1)
		if len(result) > 0 {
			headerName = result[0][1]
			continue
		}

		result = reOption.FindAllStringSubmatch(c, -1)
		if len(result) > 0 {
			if headerName != "" {
				annos[headerName] = append(annos[headerName], Option{
					K: result[0][1],
					V: result[0][2],
				})
			}
			continue
		}

		headerName = ""
	}

	return annos
}
