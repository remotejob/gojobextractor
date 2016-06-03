package findItems

import (
//	"fmt"
	gm "github.com/onsi/gomega"
	am "github.com/sclevine/agouti/matchers"
	"github.com/sclevine/agouti"
)

func FindAllOnPage(page agouti.Page) *agouti.MultiSelection {

	listResults := page.AllByClass("listResults")
	gm.Expect(listResults.At(1).AllByClass("-item")).Should(am.BeFound())
	items := listResults.At(1).AllByClass("-item")
	gm.Expect(items.Count()).Should(gm.Equal(25))

	return items

}
