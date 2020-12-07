package gonawintest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taironas/route"

	"github.com/taironas/gonawin/helpers/handlers"
	"github.com/taironas/gonawin/models"
)

// HandlerFunc type returns an httptest.ResponseRecorder when you pass the following params:
// url: the url you want to query.
// params: the parameters that the url should use.
// r: the router to record the httptest.ResponseRecorder.
type HandlerFunc func(url, method string, params map[string]string, r *route.Router) *httptest.ResponseRecorder

// GenerateHandlerFunc returns a HandlerFunc type based on the handler passed as argument.
func GenerateHandlerFunc(t *testing.T, handler func(http.ResponseWriter, *http.Request, *models.User) error) HandlerFunc {

	return func(url, method string, params map[string]string, r *route.Router) *httptest.ResponseRecorder {
		if req, err := http.NewRequest(method, url, nil); err != nil {
			t.Errorf("%v", err)
		} else {
			values := req.URL.Query()
			for k, v := range params {
				values.Add(k, v)
			}
			req.URL.RawQuery = values.Encode()
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)
			return recorder
		}
		return nil
	}
}

// TestingAuthorized runs the function pass by parameter
func TestingAuthorized(f func(w http.ResponseWriter, r *http.Request, u *models.User) error) handlers.ErrorHandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return f(w, r, nil)
	}
}
