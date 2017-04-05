package jobdetails

import (
	"github.com/remotejob/gojobextractor/domains"
	//	"fmt"
	gm "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	am "github.com/sclevine/agouti/matchers"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
	"github.com/remotejob/gojobextractor/dbhandler"
)

type JobOffer struct {
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

func NewJobOffers() *JobOffer {

	return &JobOffer{

		Id:           "",
		Company:      "",
		Title:        "",
		Location:     "",
		Tags:         []string{},
		Externallink: "",
		Email:        "",
		Hits:         0,
		Created_at:   time.Now(),
		Applied:      false,
		Description:  "",
	}

}

func (jo *JobOffer) GetAllLinks(page *agouti.Page) {

	gm.Expect(page.FindByClass("jobdetail")).Should(am.BeFound())

	jobdetails := page.FindByClass("jobdetail")

	alllinks := jobdetails.All("a")

	count_links, _ := alllinks.Count()

	for i := 0; i < count_links; i++ {

		jo.ParceLink(alllinks.At(i))

	}

}

func (jo *JobOffer) ParceLink(link *agouti.Selection) {

	now := time.Now()
	id, _ := link.Attribute("id")
	data_uri, _ := link.Attribute("data-uri")
	href, _ := link.Attribute("href")
	data_email, _ := link.Attribute("data-email")

	class, _ := link.Attribute("class")

	text, _ := link.Text()

	if class == "title job-link" {

		jo.Id = href
		jo.Created_at = now

		if text != "" {
			jo.Title = text
		}

	}

	if class == "post-tag job-link no-tag-menu" {
		if text != "" {
			jo.Tags = append(jo.Tags, text)
		}

	}
	if class == "employer" {
		if text != "" {
			jo.Company = text
		}

	}

	if id == "apply" && data_uri != "" {
		jo.Externallink = data_uri
	}

	if id == "apply" && data_email != "" {

		jo.Externallink = data_email

	}

	if href != "" && text != "" {

		if strings.HasPrefix(href, "mailto") {
			jo.Email = text
		}

	}

}

func (jo *JobOffer) FindLocation(page *agouti.Page) {
	gm.Expect(page.FindByClass("location")).Should(am.BeFound())
	location_on_page := page.FirstByClass("location")
	
	location, _ := location_on_page.Text()
	if location != "" {
		jo.Location = location
	}

}

func (jo *JobOffer) ExamDbRecord(session mgo.Session) {

	joboffers := domains.JobOffer{jo.Id, jo.Company, jo.Title, jo.Location, jo.Tags, jo.Externallink, jo.Email, jo.Hits, jo.Created_at, jo.Applied, jo.Description}

	dbhandler.InsertRecord(session, joboffers)

}
