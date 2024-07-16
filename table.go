package als_md

import "strings"

func TableAlignCenter(body []string) []string {
	for i, v := range body {
		v = strings.ReplaceAll(v, "<table>",
			`<div class="div-table" style="text-align: center;">
<table style="margin: auto">`)
		body[i] = strings.ReplaceAll(v, "</table>", "</table>\n</div>")
	}
	return body
}

/*
<thead>
<tr>
<th></th>
<th></th>
</tr>
</thead>
*/
func TableRemoveEmptyThead(body []string) []string {
	for i, v := range body {
		if !strings.Contains(v, "thead") {
			continue
		}
		var newV, thead []string
		isEmpty := true
		inThead := false
		for _, line := range strings.Split(v, "\n") {
			if line == "<thead>" {
				thead = append(thead, line)
				inThead = true
				continue
			}
			if line == "</thead>" {
				thead = append(thead, line)
				inThead = false
				if !isEmpty {
					newV = append(newV, thead...)
				}
				thead = nil
				isEmpty = true
				continue
			}
			if inThead {
				if strings.HasPrefix(line, "<th>") &&
					strings.HasSuffix(line, "</th>") &&
					line != "<th></th>" {
					isEmpty = false
				}
				thead = append(thead, line)
				continue
			}
			newV = append(newV, line)
		}
		body[i] = strings.Join(newV, "\n")
	}
	return body
}
