package als_md

import (
	"regexp"

	_ "embed"

	"github.com/go-pdf/fpdf"
)

//go:embed Roboto-Regular.ttf
var RobotoRegularBytes []byte

//go:embed Roboto-Bold.ttf
var RobotoBoldBytes []byte

type PDF struct {
	SectionList []Section

	// write each section from a new page
	NewPage bool

	// suggest: left = top = right = 10
	MarginLeft  float64
	MarginTop   float64
	MarginRight float64

	//FontHeading []byte
	//FontContent []byte

	// suggest: 16
	FontSize float64

	// suggest: 6
	LineHeight float64

	// add space between paragraphs
	// suggest: 2
	ParagraphOffset float64
}

func (p *PDF) Generate() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")

	pdf.SetMargins(p.MarginLeft, p.MarginTop, p.MarginRight)
	pdf.AddUTF8FontFromBytes("Roboto", "", RobotoRegularBytes)
	pdf.AddUTF8FontFromBytes("Roboto", "B", RobotoBoldBytes)

	for k1, sec := range p.SectionList {
		if k1 == 0 || p.NewPage {
			pdf.AddPage()
		}

		pdf.SetFont("Roboto", "B", p.FontSize)
		pdf.Write(p.LineHeight, sec.Title+"\n")
		pdf.Write(p.ParagraphOffset, "\n")

		pdf.SetFont("Roboto", "", p.FontSize)
		reg := regexp.MustCompile(`\s*(\r|\n)\s*`)
		lineList := reg.Split(sec.Transcript.Content, -1)
		for k2, v2 := range lineList {
			pdf.Write(p.LineHeight, v2+"\n")
			if k2 < len(lineList)-1 {
				pdf.Write(p.ParagraphOffset, "\n")
			}
		}
	}
	return pdf
}
