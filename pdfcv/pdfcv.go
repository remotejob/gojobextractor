package pdfcv

import (
	"fmt"
	"strconv"

	"github.com/jung-kurt/gofpdf"
	"github.com/remotejob/go_cv_pdf/toml_parser"
	"github.com/remotejob/gojobextractor/domains"
)

func CreateCV(emplayer domains.JobOffer, mytags []domains.Tags) {

	bconfig := toml_parser.Parse("/home/juno/gowork/src/github.com/remotejob/go_cv_pdf/cv.toml")
	header := []string{"Item", "Duration", "Info"}
	// jobplaceheader := []string{"Company", "Duration", "Position", "Details", "Location", "Country"}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {

		pdf.Image("/home/juno/gowork/src/github.com/remotejob/gojobextractor/images/me_alex.jpg", 10, 10, 60, 0, false, "", 0, "")
		// pdf.ImageOptions("/home/juno/gowork/src/github.com/remotejob/go_cv_pdf/images/me_alex.jpg", 10, 6, 40, 0, false, nil, 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetX(110)
		pdf.CellFormat(70, 10, "CV", "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetX(90)
		pdf.CellFormat(70, 10, "Alex Mazurov", "", 0, "C", false, 0, "")
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
		pdf.CellFormat(65, 6, "10820 Finland", "LRB", 1, "", false, 0, "")
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Speaking", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "English Italian Russian.", "LRT", 1, "", false, 0, "")
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "languages: ", "LRB", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "", "LRB", 0, "", false, 0, "")
		pdf.Ln(10)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, fmt.Sprintf("Page %d/{nb}", pdf.PageNo()), "", 0, "C", false, 0, "")
	})
	pdf.AliasNbPages("")
	pdf.AddPage()
	// pdf.SetX(40)
	pdf.CellFormat(0, 5, "ref:   "+emplayer.Id, "T", 1, "L", false, 0, "")
	// pdf.Ln(-1)
	pdf.CellFormat(0, 5, "app: "+emplayer.Title, "B", 1, "L", false, 0, "")
	// pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 8)
	if len(mytags) > 0 {
		pdf.CellFormat(70, 9, "My experience:", "B", 1, "L", false, 0, "")

		for _, tag := range mytags {
			pdf.CellFormat(70, 6, tag.Tag+" ---> "+tag.Duration, "", 1, "L", false, 0, "")
		}
		pdf.CellFormat(70, 8, "", "T", 1, "L", false, 0, "")
	} else {
		pdf.CellFormat(70, 10, "Only for your consideration.", "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 4, "I have many years of practical experience in DataBases/Communication/Voip/Webdev.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Last 5 years mostly was dedicated to Docker/Golang/JavaScript/NoSqlDB.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "But I try to keep as mush as possible updated my old knowledge C/Java...", "", 1, "L", false, 0, "")
	pdf.Ln(-1)
	pdf.CellFormat(30, 5, "I can be useful in:", "B", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Programming languages (Golang-JAVA-C-RUBY-JavaScript)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Telecommunication(Voip-SMS-Asterisk-CiscoIOS)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Web Application(WebService-Reactjs-Angular2-RubyOnRails)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Databases(Elascticsearch-MongoDB-Redis-MySql-PostgreSQL-SQLite-Oracle)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Clouds(AWS-DigitalOcean-ContainerEngine)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Virtualization(LinuxKVM-Docker-Kubernetes),", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Unit tests(Selenium-Protractor-Ginkgo)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "SEO(Organic Web Search Optimization).", "", 1, "L", false, 0, "")
	pdf.CellFormat(110, 8, "My Development IDE: Eclipse/Visual Studio Code  <---> Git(GitHub/GtLab).", "BT", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Start from 2013 most of my projects openly disposed.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "So for undestanding my capasity you need only small expertice of my projects.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "I hope you find all necessary information on Sites:", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "http://mazurov.eu", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "https://github.com/sinelga  ---> 83 projects", "", 1, "L", false, 0, "https://github.com/sinelga")
	pdf.CellFormat(0, 4, "https://github.com/remotejob ---> 51 projects", "", 1, "L", false, 0, "https://github.com/remotejob")
	pdf.CellFormat(0, 4, "https://gitlab.com/remotejob/projects ---> 13 projects", "", 1, "L", false, 0, "https://gitlab.com/remotejob/projects")
	pdf.CellFormat(0, 4, "https://hub.docker.com/u/remotejob ---> 10 images", "", 1, "L", false, 0, "https://hub.docker.com/u/remotejob")

	pdf.Ln(-1)
	pdf.CellFormat(0, 4, "PS:", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Last 10 years my position was mostly as remote programmer, actual presence at site was only about 10%,", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, " so at least for now I would prefer REMOTE JOB.", "", 1, "L", false, 0, "")
	pdf.AddPage()
	w := []float64{35, 15, 125}

	for _, cv := range bconfig.Cv {

		if cv.Name == "Telecommunication Application" {

			pdf.AddPage()
		}

		pdf.SetFont("Arial", "B", 15)
		pdf.CellFormat(0, 10, cv.Name, "", 1, "", false, 0, "")

		pdf.SetFillColor(100, 149, 237)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(128, 0, 0)
		pdf.SetLineWidth(.3)
		pdf.SetFont("Arial", "", 8)

		for j, str := range header {
			pdf.CellFormat(w[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)

		fill := false

		for i, item := range cv.Item {

			if i == len(cv.Item)-1 {

				pdf.CellFormat(w[0], 6, item.Title, "LRB", 0, "", fill, 0, "")
				pdf.CellFormat(w[1], 6, strconv.Itoa(item.Duration)+" y.", "LRB", 0, "R", fill, 0, "")
				pdf.CellFormat(w[2], 6, item.Extra, "LRB", 0, "", fill, 0, "")

			} else {

				pdf.CellFormat(w[0], 6, item.Title, "LR", 0, "", fill, 0, "")
				pdf.CellFormat(w[1], 6, strconv.Itoa(item.Duration)+" y.", "LR", 0, "R", fill, 0, "")
				pdf.CellFormat(w[2], 6, item.Extra, "LR", 0, "", fill, 0, "")
			}
			pdf.Ln(-1)

			fill = !fill

		}
	}
	err := pdf.OutputFileAndClose("/tmp/my_cv.pdf")
	if err == nil {
		fmt.Println("Successfully generated my_cv.pdf")
	} else {
		fmt.Println(err)
	}
}
