package send_emails

import (
	"github.com/remotejob/gojobextractor/dbhandler"
	"github.com/remotejob/gojobextractor/domains"
	//	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"time"

	"github.com/scorredoira/email"
	"gopkg.in/mgo.v2"
)

func SendAll(dbsession mgo.Session, emails []domains.Email, login string, pass string) {

	for _, email := range emails {

		send(dbsession, login, pass, email)
		time.Sleep(4000 * time.Millisecond)
	}

}
func send(dbsession mgo.Session, glogin string, gpass string, emailtxt domains.Email) {
	from := glogin
	pass := gpass

	msg := "\n" + emailtxt.Body

	m := email.NewMessage("ref: "+emailtxt.Subject, msg)

	m.From = mail.Address{
		Name:    "Mazurov Alex",
		Address: "support@remotejob.eu",
	}
	m.To = []string{emailtxt.To}
	//	m.Subject = "ref: " + emailtxt.Subject

	err := m.Attach("mazurov_cv.pdf")
	if err != nil {
		log.Println(err)
	}

	err = email.Send("smtp.gmail.com:587", smtp.PlainAuth("", from, pass, "smtp.gmail.com"), m)
	if err != nil {
		log.Println(err)
	}

	dbhandler.UpdateExtEmploerEmail(dbsession, emailtxt)

}
