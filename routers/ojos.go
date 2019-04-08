package routers

import (
	"github.com/google/logger"
	"github.com/gorilla/mux"
	"image"
	"image/color"
	"image/png"
	"io"
	"net/http"
	"ojos/database"
	"strconv"
	"time"
)

var (
	// Client for rest calls
	Client http.Client
	setClientFunc = setClient
)

// OjosHandler can receive:
//  - url of a web page to capture
//  - optional css selector
// Returns a png image of the page or element of the page specified with the css
// selector
func OjosHandler(w http.ResponseWriter, r *http.Request) {

	// Assign an id to each request so in future I can track these requests
	// Increment an id # in redis
	newestID, err := database.Get([]byte("newestID"))
	if err != nil {
		newestID = []byte("0");
	} else {
		intNewestID, err:= strconv.Atoi(string(newestID))
		if err != nil {
			logger.Error(err)
			return
		}
		database.Put([]byte("newestID"), []byte(strconv.Itoa(intNewestID + 1)))
	}

	// Read in query arguments from the
	// request
	vars := mux.Vars(r)

	// Send a request to capturama service to capture and generate an image in png format from
	// the specified webpage
	url, _ := vars["url"]
	dynamicSizeSelector, ok := vars["dynamicSizeSelector"]

	req, err := http.NewRequest("GET", "localhost:12345/capture", nil)
	if err != nil {
		logger.Error(err)
		return
	}

	q := req.URL.Query()
	q.Add("url", url)

	if ok {
		q.Add("dynamic_size_selector", dynamicSizeSelector)
	}

	req.URL.RawQuery = q.Encode()

	setClientFunc()
	resp, err := Client.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}

	// Return a 200 status code
	w.WriteHeader(http.StatusOK)

	// If capturama successful return image. If not, return a 1x1 pixel
	if resp.StatusCode == 200 {
		w.Header().Set("Content-Type", "image/png")
		io.Copy(w, resp.Body)
	} else {

		img := image.NewRGBA(image.Rect(0, 0, 1, 1))
		img.Set(0, 0, color.RGBA{255, 0, 0, 255})
		err := png.Encode(w, img)
		if err != nil {
			logger.Error(err)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	}

	// Get stats from number of time a particular domain is captured
	// Use redis to store count, and put outputs in log(?)
	oldCount, err := database.Get([]byte(""))
	if err != nil {
		err := database.Put([]byte(url), []byte("0"))
		if err != nil {
			logger.Error(err)
			return
		}
	} else {

		oldCountInt, err := strconv.Atoi(string(oldCount))
		newCountInt := oldCountInt + 1

		if err != nil {
			logger.Error(err)
			return
		}

		err = database.Put([]byte(url), []byte(strconv.Itoa(newCountInt)))
		if err != nil {
			logger.Error(err)
			return
		}

		logger.Infof("Successfully completed request #%d, returned image from url %s\n",
			newCountInt, "")
	}
}

func setClient()  {
	Client = http.Client{
		Transport: &http.Transport{},
		Timeout: 30 * time.Second,
	}
}
