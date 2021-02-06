package main
import (
	"fmt"
	"net/http"
	"log"
	"os"
	"encoding/json"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type RequestBody struct {
	Email string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func sendemail(w http.ResponseWriter, r *http.Request) {
	if (*r).Method == "OPTIONS" {
		return
	}

	requestBody := &RequestBody{}
	json.NewDecoder(r.Body).Decode(requestBody)
	siteMail := mail.NewEmail("Daniela", "contact@danielalins.com")

	to := mail.NewEmail("Sender", requestBody.Email)
	subject := requestBody.Subject
	message := fmt.Sprintf("%[1]s %[2]s\n", requestBody.Email, requestBody.Message)
	subjectTo := "E-mail confirmation."
	messageTo := """<div style="font-size: large;">
	<p>Hi!</p>
	<p>Thank you very much for your email. I'll respond to it as soon as possible.</p>
	<p>Best Regards,</p>
	<br>
</div>
<div>
	<h4 style="margin: 0 0">Daniela Lins</h4>
	<p style="margin: 0 0;">Full Stack Web Developer</p>
</div"""

	// Send me (Daniela) the email
	Send(siteMail, subject, message)
	Send(to, subjectTo, messageTo)

}

func Send(email *mail.Email, subject string, message string){

	from := mail.NewEmail("Daniela", "contact@danielalins.com")
	to := email
	plainTextContent := message
	htmlContent := message
	sendgridMessage := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	dat, err := ioutil.ReadFile("/go/bin/sendgrid.env")
	 check(err)
	client := sendgrid.NewSendClient(string(dat))
	response, err := client.Send(sendgridMessage)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func Handlers() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/email", sendemail).Methods("OPTIONS")
	r.HandleFunc("/email", sendemail).Methods("POST")

	return r
}


func main() {
	port := os.Getenv("PORT")
	if (port == "") {
		port = "8000"
	}

	// Handle routes
	r := Handlers()
	handler := cors.AllowAll().Handler(r)
	http.Handle("/", handler)

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))


}