package als_md

import (
	"fmt"
	"slices"
	"strings"
)

/*
similiar to html exported by Emacs org-mode

<div id="table-of-contents">
<h2>≡</h2>
<div id="text-table-of-contents">
<ul>
<li><a href="#orge7b5326">1. Linux</a></li>
<li><a href="#org95727cb">2. Windows</a>
  <ul>
  <li><a href="#orgda2e8bf">2.1. Disk</a></li>
  <li><a href="#orge4ab8ba">2.2. Network</a></li>
  </ul>
  </li>
</ul>
</div>
</div>
*/

const TOCTemplate = `
<div id="table-of-contents">
<h2>%v</h2>
<div id="text-table-of-contents">
%v
</div>
</div>
`

type TOCSectionNumber struct {
	n []int //section number
}

func NewTOCSection(section []int) TOCSectionNumber {
	return TOCSectionNumber{n: section}
}

func (t *TOCSectionNumber) String() string {
	var strList []string
	for _, v := range t.n {
		strList = append(strList, fmt.Sprint(v))
	}
	return strings.Join(strList, ".")
}

func (t *TOCSectionNumber) Copy() TOCSectionNumber {
	return NewTOCSection(t.n)
}

func (t *TOCSectionNumber) Len() int {
	return len(t.n)
}

/*
To create the first child section of the next level,
simply add a 1 to the end of the current section number.
*/
func (t *TOCSectionNumber) FirstChild() TOCSectionNumber {
	section := slices.Clone(t.n)
	section = append(section, 1)
	return NewTOCSection(section)
}

/*
For the next section at the same level, simply increment the last number.
*/
func (t *TOCSectionNumber) Next() TOCSectionNumber {
	section := slices.Clone(t.n)
	section[len(section)-1] += 1
	return NewTOCSection(section)
}

type tocItem struct {
	Title   string
	Section TOCSectionNumber
	Virtual bool
}

func (t *tocItem) Depth() int {
	return t.Section.Len()
}

func (t *tocItem) HTML() string {
	if t.Virtual {
		return ""
	}
	return fmt.Sprintf(`<a href="#sec-%v">%v %v</a>`,
		t.Section.String(), t.Section.String(), t.Title)
}

type tocTreeNode struct {
	Self     tocItem
	Children []*tocTreeNode
}

func (t *tocTreeNode) HTML() []string {
	var result []string
	result = append(result, "<li>")
	result = append(result, t.Self.HTML())
	if len(t.Children) > 0 {
		result = append(result, "<ul>")
		for _, v := range t.Children {
			result = append(result, v.HTML()...)
		}
		result = append(result, "</ul>")
	}
	result = append(result, "</li>")
	return result
}

type tocTree struct {
	Children []*tocTreeNode
}

func (t *tocTree) findParent(section []int) *tocTreeNode {
	parent := t.Children[section[0]-1]
	for i := 1; i < len(section)-1; i++ {
		parent = parent.Children[section[i]-1]
	}
	return parent
}

func (t *tocTree) Add(item tocItem) {
	if item.Depth() == 1 {
		t.Children = append(t.Children, &tocTreeNode{
			Self: item,
		})
		return
	}
	parent := t.findParent(item.Section.n)
	parent.Children = append(parent.Children, &tocTreeNode{
		Self: item,
	})
}

func (t *tocTree) HTML() []string {
	if len(t.Children) == 0 {
		return nil
	}

	var result []string
	result = append(result, "<ul>")
	for _, v := range t.Children {
		result = append(result, v.HTML()...)
	}
	result = append(result, "</ul>")
	return result
}

type TOC struct {
	Heading string
	List    []tocItem
}

// if heading == "", use defalut: ☰
func NewTOC(heading string) *TOC {
	if heading == "" {
		heading = "☰"
	}
	return &TOC{
		Heading: heading,
	}
}

func (t *TOC) NewSection(depth int) TOCSectionNumber {
	/*
		This is the first chapter.
		If depth != 1, we need to insert missing levels step by step.
		This way, the numbering of later chapters can be computed easily.
	*/
	if len(t.List) == 0 {
		s := NewTOCSection([]int{1})
		for i := 1; i < depth; i++ {
			t.Add("virtual", s, true)
			s = s.FirstChild()
		}
		return s
	}

	/*
		Since this is not the first chapter,
		we first try to find the previous section at the same depth.
		If not found, we fall back to a higher-level section.
		The first chapter has already been fully initialized, so there will always be a valid reference.
	*/
	var last tocItem
	for i := len(t.List) - 1; i >= 0; i-- {
		if t.List[i].Depth() <= depth {
			last = t.List[i]
			break
		}
	}

	/*
		Once the previous section at the same level is found,
		simply increment the last segment of its section number.
	*/
	if last.Depth() == depth {
		return last.Section.Next()
	}

	/*
		If no sibling section is found, we can assume it belongs to a parent level.
		There may be several levels missing, so we need to fill them accordingly.
	*/
	s := last.Section
	for i := 1; i < depth-last.Depth(); i++ {
		s = s.FirstChild()
		t.Add("virtual", s, true)
	}
	return s.FirstChild()
}

func (t *TOC) Add(title string, s TOCSectionNumber, virtual bool) {
	t.List = append(t.List, tocItem{
		Title:   title,
		Section: s,
		Virtual: virtual,
	})
}

func (t *TOC) HTML() string {
	if len(t.List) == 0 {
		return fmt.Sprintf(TOCTemplate, t.Heading, "")
	}

	var tree tocTree
	for _, v := range t.List {
		tree.Add(v)
	}
	return fmt.Sprintf(TOCTemplate, t.Heading, strings.Join(tree.HTML(), "\n"))
}
