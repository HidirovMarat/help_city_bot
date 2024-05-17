package request

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
)

func SendComplaint(typeComplaint, firstName, lastName, complaint string) {
	body := fmt.Sprintf(`{"firstName":"%s", "lastName":"%s", "complaint":"%s"}`, firstName, lastName, complaint)
	data := []byte(body)
	r := bytes.NewReader(data)
	add := "http://example.com/" + typeComplaint
	_, err := http.Post(add, "application/json", r)
	if err != nil {
		log.Printf(err.Error())
	}
}
