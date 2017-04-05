package findalllinks

import (
	"fmt"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"net/http"
	"strings"
)

func FindAll(urlstr string) []string {

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
		if n.DataAtom == atom.Div {
			//			fmt.Println(scrape.Attr(n, "class"))
			//			return scrape.Attr(n, "class") == "-item -job"
			return strings.HasPrefix(scrape.Attr(n, "class"), "-item -job")

		}
		return false
	}

	var links []string

	grid, ok := scrape.Find(root, scrape.ByClass("listResults"))

	if ok {

		//		fmt.Println(grid)
		gridItems := scrape.FindAll(grid, matcher)

		for _, itemA := range gridItems {

			title, ok := scrape.Find(itemA, scrape.ByClass("-title"))
			if ok {

				link, ok := scrape.Find(title, scrape.ByTag(atom.A))
				if ok {

					fmt.Println(scrape.Attr(link, "href"))
					links = append(links, scrape.Attr(link, "href"))
				}

			}
			//			fmt.Println(scrape.Attr(itemA, "href"))

		}

	}
	return links
}
