package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/driver"
	"github.com/yusufelyldrm/reservation/internal/forms"
	"github.com/yusufelyldrm/reservation/internal/helpers"
	"github.com/yusufelyldrm/reservation/internal/models"
	"github.com/yusufelyldrm/reservation/internal/render"
	"github.com/yusufelyldrm/reservation/internal/repository"
	"github.com/yusufelyldrm/reservation/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo Creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo Creates a new repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

// NewHandlers set the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home func is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.gohtml", &models.TemplateData{})
}

// About func is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	//send to data to the template
	render.Template(w, r, "about.page.gohtml", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays from
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can not get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can not find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	res.Room.RoomName = room.RoomName // add room name to reservation

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.StartDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})

}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	/*	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
		if !ok {
			m.App.Session.Put(r.Context(), "error", "can not parse form")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}*/

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can not parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can not parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	reservation = models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)
	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 2)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can not insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cant insert room restriction into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//send notification - first to guest
	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s: <br>
	This is to confirm your reservation from %s to %s.`,
		reservation.FirstName,
		reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"),
	)

	msg := models.MailData{
		To:      reservation.Email,
		From:    "you@there.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gohtml", &models.TemplateData{})
}

// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gohtml", &models.TemplateData{})
}

// Availability renders the search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.gohtml", &models.TemplateData{})
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

		// convert the start date to time.Time
		layout := "2006-01-02"
		startDate, err := time.Parse(layout, start)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// convert the end date to time.Time
		endDate, err := time.Parse(layout, end)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// get the availability, if any, for the given dates, if there is an error, return it to the user
		rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		// if no availability, redirect to the search availability page, with a flash message
		if len(rooms) == 0 {
			// no availability
			m.App.Session.Put(r.Context(), "error", "No availability")
			http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
			return
		}

		// send the data to the template
		data := make(map[string]interface{})
		data["rooms"] = rooms

		// store the data in the session
		res := models.Reservation{
			StartDate: startDate,
			EndDate:   endDate,
		}
		m.App.Session.Put(r.Context(), "reservation", res)

		err = render.Template(w, r, "choose-room.page.gohtml", &models.TemplateData{
			Data: data,
		})

		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	} else {
		// if the request method is not POST , return an error
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		fmt.Println("Error parsing start date:", err)
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		fmt.Println("Error parsing end date:", err)
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		fmt.Println("Error converting room ID:", err)
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	// convert the response to JSON
	out, err := json.MarshalIndent(resp, "", "     ")

	if err != nil {
		helpers.ServerError(w, err)
		return
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
	render.Template(w, r, "contact.page.gohtml", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Cannot get item from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get item from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(chi.URLParam(r, "id")) // get the room id from the url
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// get the reservation from the session
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID                                           // set the room id in the reservation
	m.App.Session.Put(r.Context(), "reservation", res)            // store the reservation in the session
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther) // redirect to the make reservation page
}

// BookRoom takes URL parameters, builds a session variable, and takes user to make reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, err := time.Parse(layout, ed)

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation

	res.Room.RoomName = room.RoomName
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
