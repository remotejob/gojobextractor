package send_emails

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
//	"fmt"
	"github.com/scorredoira/email"
	"gopkg.in/mgo.v2"
	"log"
	"net/mail"
	"net/smtp"
	"time"
)

func SendAll(dbsession mgo.Session, emails []domains.Email, login string, pass string) {

	for _, email := range emails {

		send(dbsession, login, pass, email)
		time.Sleep(3000 * time.Millisecond)
	}


}
func send(dbsession mgo.Session, glogin string, gpass string, emailtxt domains.Email) {
	from := glogin
	pass := gpass

	msg := "\n" +emailtxt.Body

	m := email.NewMessage("ref: " + emailtxt.Subject,msg)

	m.From = mail.Address{
		Name:    "Alex Mazurov",
		Address: "support@mazurov.eu",
	}
	m.To = []string{emailtxt.To}
//	m.Subject = "ref: " + emailtxt.Subject

	err := m.Attach("/home/juno/git/cv/version_desk_react_00/dist/mazurov_cv.pdf")
	if err != nil {
		log.Println(err)
	}

	err = email.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), m)
	if err != nil {
		log.Println(err)
	}

	//	fmt.Println("TO:",email)
	//
	////	to := "aleksander.mazurov@gmail.com"
	//	to :=email.To
	//
	//	myfrom := "support@mazurov.eu"
	//
	//	msg := "From: " + myfrom + "\n" +
	//		"To: " + to + "\n" +
	//		"Subject:ref: " + email.Subject + "\n\n" +
	//		email.Body
	//
	//	err := smtp.SendMail("smtp.gmail.com:587",
	//		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
	//		from, []string{to}, []byte(msg))
	//
	//	if err != nil {
	//		log.Printf("smtp error: %s", err)
	//		return
	//	}

	//	log.Print("sent, visit http://foobarbazz.mailinator.com")
	dbhandler.UpdateExtEmploerEmail(dbsession, emailtxt)

}
