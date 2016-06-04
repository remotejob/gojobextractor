package jobdetails_simple

import (
//	"fmt"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
	"time"
	"gopkg.in/mgo.v2"
	"github.com/remotejob/gojobextractor/domains"
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

func (jo *JobOffer) ParsePage(urlstr string) {

	resp, err := http.Get(urlstr)
	if err != nil {
		panic(err)
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		panic(err)

	}

	matcher := func(n *html.Node) bool {
		// must check for nil values
		if n.DataAtom == atom.A {
			//			fmt.Println(scrape.Attr(n, "class"))
			//						return scrape.Attr(n, "class") == "employer"
			//			return strings.HasPrefix(scrape.Attr(n, "class"), "-item -job")
			return true

		}
		return false
	}

	jo.FindLocation(root)
	jo.FindDescription(root)
	
	grid, ok := scrape.Find(root, scrape.ByClass("jobdetail"))

	if ok {
		gridItems := scrape.FindAll(grid, matcher)
		now := time.Now()

		for _, link := range gridItems {

//			fmt.Println(link)

			id := scrape.Attr(link, "id")
			data_uri := scrape.Attr(link, "data-uri")
			href := scrape.Attr(link, "href")
			data_email := scrape.Attr(link, "data-email")

			class := scrape.Attr(link, "class")

			text := scrape.Text(link)

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

	}

}

func (jo *JobOffer) FindLocation(node *html.Node) {

	loc, ok := scrape.Find(node, scrape.ByClass("location"))

	if ok {

		location := scrape.Text(loc)
		if location != "" {
			jo.Location = location
		}
	}

}
func (jo *JobOffer) FindDescription(node *html.Node) {

	des, ok := scrape.Find(node, scrape.ByClass("description"))

	if ok {

		description := scrape.Text(des)
		if description  != "" {
			jo.Description = strings.Replace(description,"Job Description","",1)
		}
	}

}

func (jo *JobOffer) ExamDbRecord(session mgo.Session) {

	joboffers := domains.JobOffer{jo.Id, jo.Company, jo.Title, jo.Location, jo.Tags, jo.Externallink, jo.Email, jo.Hits, jo.Created_at, jo.Applied, jo.Description}

	dbhandler.InsertRecord(session, joboffers)

}