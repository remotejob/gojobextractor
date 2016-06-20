package handle_internal_link

import (
	"fmt"
	gm "github.com/onsi/gomega"
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/coverletter"
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/mytags"
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	"github.com/sclevine/agouti"
	am "github.com/sclevine/agouti/matchers"
	"github.com/tebeka/selenium"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
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
}

func NewInternalJobOffers(job domains.JobOffer) *InternalJobOffer {
	//func NewInternalJobOffers() *InternalJobOffer {

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
	}

}

func (jo *InternalJobOffer) Apply_headless(dbsession mgo.Session, page selenium.WebDriver, link string) {

	page.Get(link)
	time.Sleep(time.Millisecond * 4000)
	jobdetails, err := page.FindElement(selenium.ByClassName, "jobdetail")
	if err != nil {
		fmt.Println(err.Error())
	}
	time.Sleep(time.Millisecond * 4000)
	//	fmt.Println(jobdetails)
	alllinks, err := jobdetails.FindElements(selenium.ByTagName, "a")
	if err != nil {
		fmt.Println(err.Error())
	}
	count_links := len(alllinks)

	fmt.Println("count_links", count_links)
	var idtoapply int
	idtoapply = 0

	for i := 0; i < count_links; i++ {

		if data_jobid, err := alllinks[i].GetAttribute("data-jobid"); err == nil {
			//			fmt.Println(data_jobid)
			text, _ := alllinks[i].Text()
			id, _ := alllinks[i].GetAttribute("id")

			if text == "apply now" && id == "apply" {
				idtoapply = i
				fmt.Println("apply id", idtoapply, data_jobid)
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

		jo.ElaborateFrame_headless(dbsession, page, alllinks[idtoapply])

	} else {

		fmt.Println("Can't find apply link id->", idtoapply)

	}

}

func (jo *InternalJobOffer) Apply(dbsession mgo.Session, page *agouti.Page) {

	gm.Expect(page.FindByClass("jobdetail")).Should(am.BeFound())

	jobdetails := page.FindByClass("jobdetail")

	alllinks := jobdetails.All("a")

	count_links, _ := alllinks.Count()

	var idtoapply int
	idtoapply = 0

	for i := 0; i < count_links; i++ {

		data_jobid, _ := alllinks.At(i).Attribute("data-jobid")
		href, _ := alllinks.At(i).Attribute("href")

		if data_jobid != "" {
			text, _ := alllinks.At(i).Text()
			id, _ := alllinks.At(i).Attribute("id")

			if text == "apply now" && id == "apply" {

				idtoapply = i

			}

		}

		if href != "" {

			if strings.HasPrefix(href, "mailto:") {

				emailtxt, _ := alllinks.At(i).Text()
				jo.Email = emailtxt
			}

		}

	}

	if idtoapply > 0 {

		jo.ElaborateFrame(dbsession, page, alllinks.At(idtoapply))

	} else {

		fmt.Println("Can't find apply link id->", idtoapply)

	}

}

func (jo *InternalJobOffer) ElaborateFrame_headless(dbsession mgo.Session, page selenium.WebDriver, link selenium.WebElement) {

	link.Click()
	time.Sleep(4000 * time.Millisecond)
	if form, err := page.FindElement(selenium.ByID, "apply-dialog"); err == nil {

		if allinputs, err := form.FindElements(selenium.ByTagName, "input"); err == nil {

			fmt.Println("allinputs", len(allinputs))

			for _, input := range allinputs {

				if type_atr, err := input.GetAttribute("type"); err == nil {
					if type_atr == "file" {
						input.SendKeys("mazurov_cv.pdf")
						time.Sleep(3000 * time.Millisecond)

					}

				}

			}

		}

		mytagstoinsert := mytags.GetMyTags("mytags.csv", jo.Tags)
		coverlettertxt := coverletter.Create(mytagstoinsert, "coverletter.csv")

		if coverletter, err := form.FindElement(selenium.ByID, "CoverLetter"); err == nil {

			coverletter.SendKeys(coverlettertxt)
			time.Sleep(1000 * time.Millisecond)

			if submitbtm, err := form.FindElement(selenium.ByID, "apply-submit"); err == nil {
				submitbtm.Submit()

				jo.Applied = true
				jo.UpdateApplyedEmployer(dbsession)

			}
			time.Sleep(1000 * time.Millisecond)

		}

	}

}

func (jo *InternalJobOffer) ElaborateFrame(dbsession mgo.Session, page *agouti.Page, link *agouti.Selection) {

	gm.Expect(link.Click()).To(gm.Succeed())

	gm.Expect(page.FindByID("apply-dialog")).Should(am.BeFound())
	form := page.FindByID("apply-dialog")

	time.Sleep(1000 * time.Millisecond)

	apply_form := form.FindByClass("apply-form")

	fmt.Println("apply_form!!->", apply_form)

	gm.Expect(form.FindByClass("apply-form")).Should(am.BeFound())
	time.Sleep(1000 * time.Millisecond)

	gm.Expect(form.All("input")).Should(am.BeFound())
	allinputs := form.All("input")

	count_inputs, _ := allinputs.Count()
	var idtoinput int
	var idbuttonsubmit int

	for i := 0; i < count_inputs; i++ {

		type_atr, _ := allinputs.At(i).Attribute("type")

		if type_atr != "" {

			if type_atr == "file" {

				idtoinput = i

			}

			if type_atr == "submit" {

				idbuttonsubmit = i

			}

		}

	}
	//

	gm.Expect(allinputs.At(idtoinput).UploadFile("mazurov_cv.pdf")).To(gm.Succeed())

	mytagstoinsert := mytags.GetMyTags("mytags.csv", jo.Tags)
	//	fmt.Println(mytagstoinsert)
	coverlettertxt := coverletter.Create(mytagstoinsert, "coverletter.csv")

	gm.Expect(form.FindByID("CoverLetter")).Should(am.BeFound())
	coverletter := form.FindByID("CoverLetter")
	coverletter.SendKeys(coverlettertxt)
	gm.Expect(allinputs.At(idbuttonsubmit).Submit()).To(gm.Succeed())
	gm.Expect(page.ConfirmPopup()).To(gm.Succeed())
	time.Sleep(1000 * time.Millisecond)
	jo.Applied = true

	jo.UpdateApplyedEmployer(dbsession)

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
