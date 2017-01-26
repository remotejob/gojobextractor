package pdfcv

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/remotejob/gojobextractor/domains"
)

func CreateCV(emplayer domains.JobOffer) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {

		pdf.Image("/home/juno/gowork/src/github.com/remotejob/go_cv_pdf/images/me_alex.jpg", 10, 10, 60, 0, false, "", 0, "")
		// pdf.ImageOptions("/home/juno/gowork/src/github.com/remotejob/go_cv_pdf/images/me_alex.jpg", 10, 6, 40, 0, false, nil, 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetX(110)
		pdf.CellFormat(70, 10, "bconfig.Subtitle", "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetX(90)
		pdf.CellFormat(70, 10, "bconfig.Maintitle", "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 10)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Phone:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "+358451202801", "LRT", 0, "", false, 0, "tel:+358451202801")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Email:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "support@mazurov.eu", "LRT", 0, "", false, 0, "mail:support@mazurov.eu")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Web:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "http://mazurov.eu", "LRT", 0, "", false, 0, "http://mazurov.eu")
		pdf.Ln(-1)

		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Skype:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "mazurovfi", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Address:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "Hogberginkuja 1 Lappohja", "LRT", 0, "", false, 0, "")
		// pdf.Ln(-1)
		// pdf.SetX(95)
		// pdf.CellFormat(20, 6, "", "LR", 0, "", false, 0, "")
		// pdf.CellFormat(65, 6, "Lappohja", "LR", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "", "LRB", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "10820 Finland", "LRB", 0, "", false, 0, "")
		pdf.Ln(20)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()

	err := pdf.OutputFileAndClose("/tmp/my_cv.pdf")
	if err == nil {
		fmt.Println("Successfully generated my_cv.pdf")
	} else {
		fmt.Println(err)
	}
}
