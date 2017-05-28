package handle_internal_link

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/jung-kurt/gofpdf"
	"github.com/remotejob/go_cv_pdf/toml_parser"
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/coverletter"
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/mytags"

	"github.com/remotejob/gojobextractor/domains"
	"github.com/remotejob/gojobextractor/elasticLoader/loader/dbhandler"
	"github.com/tebeka/selenium"
	"gopkg.in/mgo.v2"
)

type InternalJobOffer struct {
	Id           string
	Company      string
	Title        string
	Location     string
	Tags         []string
	Externallink string
	Email        string
	Hits         int
	Created_at   time.Time
	Applied      bool
	Description  string
	Login        string
}

func NewInternalJobOffers(job domains.JobOffer, login string) *InternalJobOffer {

	return &InternalJobOffer{

		Id:           job.Id,
		Company:      job.Company,
		Title:        job.Title,
		Location:     job.Location,
		Tags:         job.Tags,
		Externallink: job.Externallink,
		Email:        job.Email,
		Hits:         job.Hits,
		Created_at:   job.Created_at,
		Applied:      job.Applied,
		Description:  job.Description,
		Login:        login,
	}

}

func (jo *InternalJobOffer) Apply_headless(dbsession mgo.Session, page selenium.WebDriver, link string, cvpdf string) bool {

	reCaph := false

	page.Get(link)
	time.Sleep(time.Millisecond * 4500)
	jobdetails, err := page.FindElement(selenium.ByClassName, "jobdetail")
	if err != nil {

		fmt.Println(err.Error())
		if strings.HasPrefix(err.Error(), "no such element") {
			// if err.Error() == "no such element" {

			fmt.Println("Check for Page not found 1")

			_, err := page.FindElement(selenium.ByID, "jobs-not-found")
			if err != nil {
				fmt.Println(err.Error())
			} else {

				fmt.Println("It's Page not found!! make Applie True an contunue")
				jo.Applied = true
				jo.UpdateApplyedEmployer(dbsession)

			}

		}

	} else {
		time.Sleep(time.Millisecond * 4500)

		alllinks, err := jobdetails.FindElements(selenium.ByTagName, "a")
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Check for Page not found 2")

		}
		count_links := len(alllinks)

		fmt.Println("count_links", count_links)
		var idtoapply int
		idtoapply = 0

		var applybtm []selenium.WebElement

		for i := 0; i < count_links; i++ {

			// if data_jobid, err := alllinks[i].GetAttribute("data-jobid"); err == nil {

			if data_jobid, err := alllinks[i].GetAttribute("data-analyticurl"); err == nil {

				text, _ := alllinks[i].Text()
				href, _ := alllinks[i].GetAttribute("href")
				// id, _ := alllinks[i].GetAttribute("id")

				// if text == "apply now" && id == "apply" {

				if text == "Apply now" && strings.HasPrefix(href, "https://stackoverflow.com/jobs/apply/") {
					idtoapply = i
					fmt.Println("apply id", idtoapply, data_jobid)

					applybtm = append(applybtm, alllinks[i])

				}

			}

			if href, err := alllinks[i].GetAttribute("href"); err == nil {
				if strings.HasPrefix(href, "mailto:") {

					emailtxt, _ := alllinks[i].Text()
					jo.Email = emailtxt
				}

			}

		}
		if idtoapply > 0 {
			file, err := os.OpenFile("applyto.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if _, err = file.WriteString(link + "\n"); err != nil {
				panic(err)
			}
			time.Sleep(time.Millisecond * 1000)
			reCaph = jo.ElaborateFrame_headless(dbsession, page, applybtm[0], cvpdf)

		} else {

			fmt.Println("Can't find apply link id->", link)

			file, err := os.OpenFile("cant_apply.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if _, err = file.WriteString(link + "\n"); err != nil {
				panic(err)
			}

		}

	}

	return reCaph
}

func (jo *InternalJobOffer) ElaborateFrame_headless(dbsession mgo.Session, page selenium.WebDriver, link selenium.WebElement, cvpdf string) bool {

	var mytagstoinsert []domains.Tags
	reCaph := false
	time.Sleep(2000 * time.Millisecond)

	// link.Click()

	if err := link.Click(); err != nil {

		fmt.Println("error clicking ", err.Error())
		fmt.Println(link.Location())
		if err := link.Click(); err != nil {
			fmt.Println("SECOND error clicking!! ", err.Error())

		}

	} else {
		fmt.Println("Click on link OK")

	}
	time.Sleep(6000 * time.Millisecond)
	if frms, err := page.FindElements(selenium.ByTagName, "iframe"); err == nil {

		for _, frm := range frms {

			if frmtitle, err := frm.GetAttribute("title"); err == nil {
				log.Println(frmtitle)

				if frmtitle == "recaptcha widget" {

					reCaph = true

					return reCaph

				}
			}
		}

	}

	if form, err := page.FindElement(selenium.ByID, "file-upload-form"); err == nil {

		mytagstoinsert = mytags.GetMyTags("mytags.csv", jo.Tags)
		if allinputs, err := form.FindElements(selenium.ByTagName, "input"); err == nil {

			fmt.Println("allinputs", len(allinputs))
			for _, input := range allinputs {
				if type_atr, err := input.GetAttribute("type"); err == nil {
					if type_atr == "file" {

						log.Println("need create new PDF file")
						jo.CreatePdfCv(mytagstoinsert)
						time.Sleep(3000 * time.Millisecond)
						input.SendKeys(cvpdf)
						time.Sleep(3000 * time.Millisecond)

					}
				}
			}
		}

	}
	if textarea, err := page.FindElement(selenium.ByID, "CoverLetter"); err == nil {
		log.Println("CoverLetter OK")

		coverlettertxt := coverletter.Create(mytagstoinsert, "coverletter_simple.csv")
		time.Sleep(3000 * time.Millisecond)
		textarea.Clear()

		textarea.SendKeys(coverlettertxt)

		time.Sleep(4000 * time.Millisecond)

	}

	if custom_questions, err := page.FindElements(selenium.ByTagName, "select"); err == nil {

		log.Println("custom_questions len", len(custom_questions))

		for _, customQuestion := range custom_questions {

			if options, err := customQuestion.FindElements(selenium.ByTagName, "option"); err == nil {
				log.Println("options1 len", len(options))

				options[1].Click()

				time.Sleep(3000 * time.Millisecond)
				// log.Println("Find other Submit subbuttom")

			} else {

				log.Println("Client Options !!")
			}

		}

	}

	if customInputs, err := page.FindElements(selenium.ByTagName, "input"); err == nil {

		log.Println("customInputs len", len(customInputs))

		for _, customInput := range customInputs {

			if idatt, err := customInput.GetAttribute("id"); err == nil {

				// log.Println(idatt)

				if strings.HasPrefix(idatt, "custom-question") {

					customInput.SendKeys("http://mazurov.eu")

				}

			}

		}
	}

	time.Sleep(3000 * time.Millisecond)

	if inputels, err := page.FindElements(selenium.ByClassName, "j-apply-btn"); err == nil {

		log.Println("SubmitOK", len(inputels))

		inputels[0].Click()

	} else {

		log.Println("NOO submit!!!!")

	}

	time.Sleep(5000 * time.Millisecond)

	//*[@id="content"]/div[2]/div/div[1]/div[1]/div[1]

	if success, err := page.FindElement(selenium.ByXPATH, "//*[@id=\"content\"]/div[2]/div/div[1]/div[1]/div[1]"); err == nil {

		textSuccess, _ := success.Text()
		log.Println("index-hedMessage_success OK", textSuccess, jo.Login)

		if strings.HasPrefix(textSuccess, "Awesome!") {

			jo.Applied = true
			jo.UpdateApplyedEmployer(dbsession)

			file, err := os.OpenFile("good_account.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if _, err = file.WriteString(jo.Login + "\n"); err != nil {
				panic(err)
			}

		} else {

			log.Println("No text Awesome!!!!!!")

			jo.MarkBadAccount("bad_account.txt")

			reCaph = true

			return reCaph

		}

	} else {

		jo.MarkBadAccount("bad_account.txt")

		reCaph = true

		return reCaph

	}

	return reCaph

}

func (jo *InternalJobOffer) MarkBadAccount(filestr string) {

	log.Println("index-hedMessage_success NO success !!!!?", jo.Login)

	file, err := os.OpenFile(filestr, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if _, err = file.WriteString(jo.Login + " " + jo.Id + "\n"); err != nil {
		panic(err)
	}
}

func (jo *InternalJobOffer) UpdateApplyedEmployer(dbsession mgo.Session) {

	applyedemployer := domains.JobOffer{
		Id:           jo.Id,
		Company:      jo.Company,
		Title:        jo.Title,
		Location:     jo.Location,
		Tags:         jo.Tags,
		Externallink: jo.Externallink,
		Email:        jo.Email,
		Hits:         jo.Hits,
		Created_at:   jo.Created_at,
		Applied:      jo.Applied,
		Description:  jo.Description,
	}

	dbhandler.UpdateEmployer(dbsession, applyedemployer)

}
func (jo *InternalJobOffer) CreatePdfCv(tagstoinsert []domains.Tags) {
	bconfig := toml_parser.Parse("cv.toml")
	header := []string{"Item", "Duration", "Info"}

	jobplaceheader := []string{"Company", "Duration", "Position", "Details", "Location", "Country"}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetHeaderFunc(func() {

		pdf.Image("images/me_alex.jpg", 10, 10, 60, 0, false, "", 0, "")
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
	pdf.CellFormat(0, 5, "ref:   "+jo.Id, "T", 1, "L", false, 0, "")
	// pdf.Ln(-1)
	pdf.CellFormat(0, 5, "app: "+jo.Title, "B", 1, "L", false, 0, "")
	// pdf.Ln(-1)
	pdf.SetFont("Arial", "B", 8)
	if len(tagstoinsert) > 0 {
		pdf.CellFormat(70, 9, "My experience:", "B", 1, "L", false, 0, "")

		for _, tag := range tagstoinsert {
			pdf.CellFormat(70, 6, tag.Tag+" ---> "+tag.Duration, "", 1, "L", false, 0, "")
		}
		pdf.CellFormat(70, 8, "", "T", 1, "L", false, 0, "")
	} else {
		pdf.CellFormat(70, 10, "Only for your consideration.", "", 1, "L", false, 0, "")
	}
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 4, "I have many years of practical experience in DataBases/Communication/Voip/Webdev.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Last 5 years mostly was dedicated to Docker/Golang/JavaScript/NoSqlDB.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "But I try to keep as much as possible updated my old knowledge C/Java...", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Actually I have contract with Middlesex University (London) but it's temporary till 2018", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "project related to elaboration Big Data (RDF/Elasticsearch technology)", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "so I think it's time to look for a new job.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Please do not hesitate to contact me if you need any further clarification.", "", 1, "L", false, 0, "")
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
	pdf.CellFormat(110, 8, "My Development IDE: Eclipse/Visual Studio Code  <---> Git(GitHub/GitLab).", "BT", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "Start from 2013 most of my projects openly disposed.", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 4, "So for understanding my capacity you need only small expertice of my projects.", "", 1, "L", false, 0, "")
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

	bjob := toml_parser.ParseWorkPlaces("/home/juno/gowork/src/github.com/remotejob/gojobextractor/jobs.toml")

	pdf.SetHeaderFunc(func() {

		pdf.Image("/home/juno/gowork/src/github.com/remotejob/go_cv_pdf/images/me_alex.jpg", 10, 10, 60, 0, false, "", 0, "")
		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 10)
		pdf.SetX(110)
		pdf.CellFormat(70, 10, bjob.Subtitle, "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "B", 15)
		pdf.SetX(90)
		pdf.CellFormat(70, 10, bjob.Maintitle, "", 0, "C", false, 0, "")
		pdf.Ln(-1)
		pdf.SetFont("Arial", "", 10)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Phone:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "+358451202801", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Email:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "support@mazurov.eu", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Web:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "http://mazurov.eu", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)

		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Skype:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "mazurovfi", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "Address:", "LRT", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "Hogberginkuja 1 Lappohja", "LRT", 0, "", false, 0, "")
		pdf.Ln(-1)
		pdf.SetX(95)
		pdf.CellFormat(20, 6, "", "LRB", 0, "", false, 0, "")
		pdf.CellFormat(65, 6, "10820 Finland", "LRB", 0, "", false, 0, "")

		pdf.Ln(20)
	})

	pdf.AddPage()

	wjobplace := []float64{35, 15, 20, 55, 35, 15}

	for _, job := range bjob.Jobs {

		pdf.SetFont("Arial", "B", 15)
		pdf.CellFormat(0, 10, job.Name, "", 1, "", false, 0, "")

		pdf.SetFillColor(100, 149, 237)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetDrawColor(128, 0, 0)
		pdf.SetLineWidth(.3)
		pdf.SetFont("Arial", "", 8)

		for j, str := range jobplaceheader {
			pdf.CellFormat(wjobplace[j], 7, str, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)

		pdf.SetFillColor(224, 235, 255)
		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("", "", 0)

		fill := false

		for i, item := range job.Item {

			if i == len(job.Item)-1 {

				pdf.CellFormat(wjobplace[0], 6, item.Title, "LRB", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[1], 6, item.Duration, "LRB", 0, "R", fill, 0, "")
				pdf.CellFormat(wjobplace[2], 6, item.Position, "LRB", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[3], 6, item.Details, "LRB", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[4], 6, item.Location, "LRB", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[5], 6, item.Country, "LRB", 0, "", fill, 0, "")

			} else {

				pdf.CellFormat(wjobplace[0], 6, item.Title, "LR", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[1], 6, item.Duration, "LR", 0, "R", fill, 0, "")
				pdf.CellFormat(wjobplace[2], 6, item.Position, "LR", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[3], 6, item.Details, "LR", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[4], 6, item.Location, "LR", 0, "", fill, 0, "")
				pdf.CellFormat(wjobplace[5], 6, item.Country, "LR", 0, "", fill, 0, "")

			}
			pdf.Ln(-1)

			fill = !fill

		}

	}

	err := pdf.OutputFileAndClose("/tmp/mazurov_cv.pdf")
	if err == nil {
		fmt.Println("Successfully generated mazurov_cv.pdf")
	} else {
		fmt.Println(err)
	}
}
