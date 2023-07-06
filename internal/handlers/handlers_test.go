package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yusufelyldrm/reservation/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarter", "GET", http.StatusOK},
	{"mj", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	/*
		{"post-search-availability", "/search-availability", "POST", []postData{
			{key: "start", value: "2023-01-01"},
			{key: "end ", value: "2023-01-02"},
		}, http.StatusOK},
		{"post-search-availability-json", "/search-availability-json", "POST", []postData{
			{key: "start", value: "2023-01-01"},
			{key: "end ", value: "2023-01-02"},
		}, http.StatusOK},
		{"post-make-reservation", "/make-reservation", "POST", []postData{
			{key: "first_name", value: "Yusuf"},
			{key: "last_name ", value: "Smith"},
			{key: "email ", value: "test@ysf.com"},
			{key: "phone  ", value: "555-555-55-55"},
		}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	//test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		res, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if res.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, res.StatusCode)
		}
	}

}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	//test case where reservation is not in session(reset everthing)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test with non-existent room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := `start_date=2023-01-01`
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Yusuf")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Elyıldırım")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=yusufelyildirim@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	//test form missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for missing post body : got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test form missing invalid start date
	reqBody = `start_date=invalid`
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2023-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=Yusuf")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Elyıldırım")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=yusufelyildirim@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
