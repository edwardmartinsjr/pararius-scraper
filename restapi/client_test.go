package restapi

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Jexia REST API client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

func setup(c C) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	url, err := url.Parse(server.URL)
	c.So(err, ShouldBeNil)

	client = New(url)
	c.So(err, ShouldBeNil)
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("request method: %v, want %v", got, want)
	}
}

func TestPost(t *testing.T) {
	Convey("Given a mocked http client", t, func(c C) {
		var jsonStr = []byte(`{
			"method": "ums",
			"email": "mail@mail",
			"password": "awesomepwd"
		  }`)

		Convey("Test API request with non-specific error condition", func() {
			setup(c)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.WriteHeader(http.StatusOK)
			})

			_, err := client.post("", "", jsonStr, "")
			So(err, ShouldBeNil)
		})

		Convey("Test API request with InternalServerError error condition", func() {
			setup(c)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.WriteHeader(http.StatusInternalServerError)
			})

			_, err := client.post("", "", jsonStr, "")
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, errors.New("500 Internal Server Error"))
		})

		Convey("Test API request with BadRequest error condition", func() {
			setup(c)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.WriteHeader(http.StatusBadRequest)
			})

			_, err := client.post("", "", jsonStr, "")
			So(err, ShouldNotBeNil)
			So(err, ShouldBeError, errors.New("400 Bad Request"))
		})

		Convey("Test API request with OK success condition", func() {
			setup(c)

			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, "POST")
				w.WriteHeader(http.StatusOK)
				fmt.Fprint(w, `{"success": true}`)
			})

			client, err := client.post("", "", jsonStr, "")
			So(client.Success, ShouldEqual, true)
			So(err, ShouldBeNil)
		})
	})
}
