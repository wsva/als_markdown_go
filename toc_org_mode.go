package als_md

import (
	"fmt"
	"strings"
)

/*
仿照Emacs中org-mode导出的html格式
*/

/*
<div id="table-of-contents">
<h2>目录</h2>
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

// 下一级的第一个章节，再追加一个1即可
func (t *TOCSectionNumber) FirstChild() TOCSectionNumber {
	section := append([]int{}, t.n...)
	section = append(section, 1)
	return NewTOCSection(section)
}

// 同一级的下一个章节，末尾序号加1即可
func (t *TOCSectionNumber) Next() TOCSectionNumber {
	section := append([]int{}, t.n...)
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

func NewTOC(heading string) *TOC {
	return &TOC{
		Heading: heading,
	}
}

func (t *TOC) NewSection(depth int) TOCSectionNumber {
	/*
		第一个章节
		如果depth!=1，那就要一级一级地补上
		这样才能确保后面的章节序号都能很方便的计算出来
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
		不是第一个章节了
		首先要找到前一个同级的Section
		如果找不到，那就找上级的
		因为第一个章节已经全量初始化了，所以肯定能找到
	*/
	var last tocItem
	for i := len(t.List) - 1; i >= 0; i-- {
		if t.List[i].Depth() <= depth {
			last = t.List[i]
			break
		}
	}

	/*
		找到了前一个同级的
		直接把Section最后一位加一就行了
	*/
	if last.Depth() == depth {
		return last.Section.Next()
	}

	/*
		如果没找到同级的，那么肯定就是上级的了
		可能相差多级，要把缺的补上
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
