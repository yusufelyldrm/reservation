package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/forms"
	"github.com/yusufelyldrm/reservation/internal/models"
	"github.com/yusufelyldrm/reservation/internal/render"
	"log"
	"net/http"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo Creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home func is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIp)

	render.RenderTemplate(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About func is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	//perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello , its me again"

	//get remote ip
	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")

	stringMap["remote_ip"] = remoteIp

	//send to data to the template
	render.RenderTemplate(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make a reservation page and displays from
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})

}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)
	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 2, r)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.gohtml", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {

	// check if the request method is POST , do the following
	if r.Method == http.MethodPost {
		// process form data here
		err := r.ParseForm()
		if err != nil {
			log.Fatal(err)
			return
		}

		// get the data from the form
		start := r.Form.Get("start")
		end := r.Form.Get("end")

		//print response
		n, err := w.Write([]byte(fmt.Sprintf("Start date is %s end date is %s ", start, end)))
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
		log.Printf("Wrote %d bytes ", n)
	} else {
		// if the request method is not POST , return an error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}
	// convert the response to JSON
	out, err := json.MarshalIndent(resp, "", "     ")

	// if there is an error , log it and return an internal server error
	if err != nil {
		log.Println(err)
	}

	// set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// write the JSON response
	_, err = w.Write(out)

	if err != nil {
		log.Fatal(err)
		return
	}
}

// Contact renders the contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.gohtml", &models.TemplateData{})
}
