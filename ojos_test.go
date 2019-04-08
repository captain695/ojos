package routers

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"ojos/database"
	"os"
	"testing"
)

func TestOjosHandler(t *testing.T) {

	gock.Observe(gock.DumpRequest)
	setClient()
	gock.InterceptClient(&Client)
	setClientFunc = func() {}

	// Create log
	lf, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	logger.Init("ojos", true, true, lf).Close()

	t.Run("Capturama fails with one parameter", func(t *testing.T) {
		defer gock.Off()

		gock.New("localhost:12345").
			Get("/capture?url=http%3A%2F%2Fwww.url.com").
			Reply(400)

		database.Open()
		defer database.Close()

		request := httptest.NewRequest("GET", "/ojos?url=http://www.url.com", nil)
		request = mux.SetURLVars(request, map[string]string{"url": "http://www.url.com"})
		response := httptest.NewRecorder()
		OjosHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Capturama fails with two parameters", func(t *testing.T) {
		defer gock.Off()

		gock.New("localhost:12345").
			Get("/capture?url=http%3A%2F%2Fwww.url.com&dynamic_size_selector=body%20h1").
			Reply(400)

		database.Open()
		defer database.Close()

		request := httptest.NewRequest("GET", "/ojos?url=http://www.url.com", nil)
		request = mux.SetURLVars(request, map[string]string{"url": "http://www.url.com"})
		response := httptest.NewRecorder()
		OjosHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("Capturama successful", func(t *testing.T) {
		defer gock.Off()

		gock.New("localhost:12345").
			Get("/capture?url=http%3A%2F%2Fwww.url.com").
			Reply(200)

		database.Open()
		defer database.Close()

		request := httptest.NewRequest("GET", "/ojos?url=http://www.url.com", nil)
		request = mux.SetURLVars(request, map[string]string{"url": "http://www.url.com"})
		response := httptest.NewRecorder()
		OjosHandler(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	assert.True(t, gock.IsDone())
}
