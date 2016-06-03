package handle_internal_link

import (
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/coverletter"
	"github.com/remotejob/gojobextractor/apply_for_job/handle_internal_link/mytags"
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	"fmt"
	gm "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	am "github.com/sclevine/agouti/matchers"
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

//				fmt.Println("to click", text, id)
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
//		fmt.Println("idtoapply",idtoapply)
//		fmt.Println(alllinks.At(idtoapply).Text())
		jo.ElaborateFrame(dbsession, page, alllinks.At(idtoapply))

	} else {

		fmt.Println("Can't find apply link id->", idtoapply)

	}

}

func (jo *InternalJobOffer) ElaborateFrame(dbsession mgo.Session, page *agouti.Page, link *agouti.Selection) {

	gm.Expect(link.Click()).To(gm.Succeed())

	gm.Expect(page.FindByID("apply-dialog")).Should(am.BeFound())
	form := page.FindByID("apply-dialog")

// stop hear
	time.Sleep(1000 * time.Millisecond)
	
	apply_form := form.FindByClass("apply-form")
	
	fmt.Println("apply_form!!->",apply_form)
	
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
	fmt.Println("idtoinput", idtoinput)
	fmt.Println("idbuttonsubmit", idbuttonsubmit)
	//
	//	gm.Expect(allinputs.At(idtoinput).Click()).To(gm.Succeed())
	gm.Expect(allinputs.At(idtoinput).UploadFile("/home/juno/git/cv/version_desk_react_00/dist/mazurov_cv.pdf")).To(gm.Succeed())

	mytagstoinsert := mytags.GetMyTags("/home/juno/git/jobprotractor/gojobextractor/mytags.csv", jo.Tags)
//	fmt.Println(mytagstoinsert)
	coverlettertxt := coverletter.Create(mytagstoinsert, "/home/juno/git/jobprotractor/gojobextractor/coverletter.csv")

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
