package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	postedData := url.Values{}
	postedData.Add("start_date", "2023-01-01")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//test form missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
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
	postedData.Add("start_date", "invalid")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test form missing invalid end date
	postedData.Add("end_date", "invalid")
	postedData.Add("start_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test form missing invalid room id
	postedData.Add("start_date", "2023-01-01")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "invalid_room_id")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid room_id: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
	//test form missing invalid data

	postedData.Add("start_date", "2023-01-01")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "y")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for failure to insert reservation to database

	postedData.Add("start_date", "2023-01-01")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "2")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler failed when trying to inserting reservation: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//test for failure to insert restriction to database

	postedData.Add("start_date", "2023-01-01")
	postedData.Add("end_date", "2023-01-02")
	postedData.Add("first_name", "Yusuf")
	postedData.Add("last_name", "Elyıldırım")
	postedData.Add("email", "yusufelyildirim@gmail.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1000")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepoitory_AvailabilityJSON(t *testing.T) {
	//first case-room are not available

	reqBody := "start_date=2023-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2023-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	//create a request
	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))

	//get context with session
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	//set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//make handler handlefunc
	handler := http.HandlerFunc(Repo.AvailabilityJSON)

	//make response recorder
	rr := httptest.NewRecorder()

	//serve http
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json")
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
