package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/yhat/scrape"

	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	mgmail "google.golang.org/api/gmail/v1"
)

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("gmail-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, mgmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := mgmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	user := "me"

	// log.Println("start")

	job, err := srv.Users.Messages.List(user).Q("from:do-not-reply@stackoverflow.com").Do()

	// job, err := srv.Users.Messages.List(user).Q("from:aleksander.mazurov@gmail.com").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels. %v", err)
	}
	for _, m := range job.Messages {

		// log.Println("m", m)

		msg, err := srv.Users.Messages.Get(user, m.Id).Format("full").Do()
		if err != nil {
			log.Fatalf("Unable to retrieve message %v: %v", m.Id, err)
		}

		// log.Println(msg.Payload.Parts)
		matcher := func(n *html.Node) bool {
			// must check for nil values
			if n.DataAtom == atom.A {

				return true

			}
			return false
		}
		// for _, head := range msg.Payload.Headers {

		// 	// log.Println(head)

		// 	if head.Name == "To" {

		// 		file, err := os.OpenFile("goodaccounts.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		// 		if err != nil {
		// 			panic(err)
		// 		}
		// 		defer file.Close()

		// 		log.Println(head.Value)
		// 		if _, err = file.WriteString(head.Value + "\n"); err != nil {
		// 			panic(err)
		// 		}
		// 		file.Close()

		// 	}

		// }
		for _, part := range msg.Payload.Parts {

			// log.Println(part)
			if part.MimeType == "text/html" {

				// log.Println(part.Body.Data)

				data, _ := base64.URLEncoding.DecodeString(part.Body.Data)

				root, err := html.Parse(bytes.NewReader(data))
				if err != nil {

					log.Fatalln(err.Error())

				}

				links := scrape.FindAll(root, matcher)
				for _, link := range links {

					// log.Println(link)
					href := scrape.Attr(link, "href")

					if strings.HasPrefix(href, "https://stackoverflow.com/users/signup-finish") {
						log.Println(href)

						file, err := os.OpenFile("signuplinks.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
						if err != nil {
							panic(err)
						}
						defer file.Close()

						if _, err = file.WriteString(href + "\n"); err != nil {
							panic(err)
						}
						file.Close()

					}

				}
			}

		}

	}

}
