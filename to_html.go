package als_md

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
)

const (
	topContainer = "<div id='top-container' class='top-container'></div>\n"
)

// return TOC, body
func ToHTML(tocHeading, content string) (string, string) {
	toc := NewTOC(tocHeading)
	var body []string
	var piece string

	if strings.Contains(content, "\r") {
		content = FormatNewline(content)
	}

	regHead := regexp.MustCompile(`^(#+) +(.*)$`)

	for _, line := range strings.Split(content, "\n") {
		if regHead.MatchString(line) {
			if piece != "" {
				html := markdown.ToHTML([]byte(piece), nil, nil)
				body = append(body, string(html))
				piece = ""
			}

			m := regHead.FindStringSubmatch(line)
			s := toc.NewSection(len(m[1]))
			toc.Add(m[2], s, false)
			body = append(body,
				fmt.Sprintf(`<h%v id="sec-%v"><span class="section-number-%v">%v</span> %v</h%v>`,
					len(m[1])+1, s.String(), len(m[1])+1, s.String(), m[2], len(m[1])+1))
		} else {
			piece += line + "\n"
		}
	}

	if piece != "" {
		html := markdown.ToHTML([]byte(piece), nil, nil)
		body = append(body, string(html))
		piece = ""
	}

	body = TableAlignCenter(body)
	body = TableRemoveEmptyThead(body)

	return toc.HTML(), topContainer + strings.Join(body, "\n")
}
