package als_md

import "strings"

/*
Mac OS: \r => \n
Windows: \r\n => \n
*/
func FormatNewline(content string) string {
	var b strings.Builder
	var crlf string
	for _, v := range content {
		if v == '\r' || v == '\n' {
			crlf += string(v)
		} else {
			if len(crlf) > 0 {
				crlf = strings.ReplaceAll(crlf, "\r\n", "\n")
				crlf = strings.ReplaceAll(crlf, "\r", "\n")
				b.WriteString(crlf)
				crlf = ""
			}
			b.WriteRune(v)
		}
	}
	return b.String()
}
