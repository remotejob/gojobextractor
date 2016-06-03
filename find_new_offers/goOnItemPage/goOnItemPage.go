package goOnItemPage

import (
	gm "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	am "github.com/sclevine/agouti/matchers"
)

func GoOn(item *agouti.Selection) {

	gm.Expect(item.FindByClass("-title")).Should(am.BeFound())
	title := item.FindByClass("-title")

	gm.Expect(title.FindByClass("job-link")).Should(am.BeFound())
	job_link := title.FindByClass("job-link")

	gm.Expect(job_link.Click()).Should(gm.Succeed())

}
