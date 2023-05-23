package render

import (
	"encoding/gob"
	"github.com/alexedwards/scs/v2"
	"github.com/yusufelyldrm/reservation/internal/config"
	"github.com/yusufelyldrm/reservation/internal/models"
	"net/http"
	"os"
	"testing"
	"time"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	//what am I going to put in the session
	gob.Register(models.Reservation{})

	//change this to true when in production
	testApp.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = testApp.InProduction

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}
