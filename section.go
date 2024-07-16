package als_md

import (
	"fmt"
	"regexp"
	"strings"
)

/*
code block
*/
type Block struct {
	Language string
	Content  string
}

type Section struct {
	Title           string
	AudioHref       string
	Transcript      Block
	TranslationList []Block
	NoteList        []Block
}

func SplitSection(content string) ([]Section, error) {
	if strings.Contains(content, "\r") {
		content = FormatNewline(content)
	}

	var section_list []Section
	var section *Section

	in_code_block := false
	code_block_index := 0
	code_block_language := ""
	code_block_content := ""

	regAudio := regexp.MustCompile(`src="(.+)"></audio>`)

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "# ") && !in_code_block {
			if section != nil {
				section_list = append(section_list, *section)
			}
			section = &Section{
				Title: line[2:],
			}
			code_block_index = 0
			continue
		}
		if strings.HasPrefix(line, "<audio ") && !in_code_block {
			submatch := regAudio.FindStringSubmatch(line)
			if len(submatch) > 1 {
				section.AudioHref = submatch[1]
			}
			continue
		}
		if strings.HasPrefix(line, "`````") && !in_code_block {
			in_code_block = true
			code_block_index += 1
			code_block_language = strings.ReplaceAll(line, "`````", "")
			continue
		}
		if strings.HasPrefix(line, "`````") && in_code_block {
			block := Block{
				Language: code_block_language,
				Content:  code_block_content,
			}
			if code_block_index == 1 {
				section.Transcript = block
			} else if code_block_language != "" {
				section.TranslationList = append(section.TranslationList, block)
			} else {
				section.NoteList = append(section.NoteList, block)
			}
			in_code_block = false
			code_block_language = ""
			code_block_content = ""
			continue
		}
		if in_code_block {
			code_block_content += line + "\n"
		} else {
			if line != "" {
				fmt.Printf("invalid line: %s\n", line)
			}
		}
	}
	if section != nil {
		section_list = append(section_list, *section)
	}
	return section_list, nil
}
