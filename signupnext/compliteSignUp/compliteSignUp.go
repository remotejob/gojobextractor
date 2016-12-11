package compliteSignUp

import (
	"fmt"
	"time"

	"github.com/tebeka/selenium"
)

func Complite(link string) {

	caps := selenium.Capabilities{"browserName": "chrome"}
	//				caps := selenium.Capabilities{"browserName": "phantomjs"}
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer wd.Quit()

	wd.Get(link)

	time.Sleep(time.Millisecond * 3000)

	// Read in an integer
	var i int
	_, err = fmt.Scanln(&i)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())

		// If int read fails, read as string and forget
		var discard string
		fmt.Scanln(&discard)
		return
	}

	time.Sleep(time.Millisecond * 1500)

}
