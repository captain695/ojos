package main

import (
  "github.com/google/logger"
  "github.com/gorilla/mux"
  "net/http"
  "ojos/database"
  "ojos/routers"
  "os"
)

func main() {

  // Register endpoints
  rtr := mux.NewRouter()
  rtr.HandleFunc("/ojos", routers.OjosHandler).
    Queries("url", "{url}", "dynamic_size_selector", "{dynamicSizeSelector}")
  rtr.HandleFunc("/ojos", routers.OjosHandler).
  Queries("url", "{url}")

  // Create log
  lf, err := os.OpenFile("log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
  if err != nil {
    logger.Fatalf("Failed to open log file: %v", err)
  }

  defer lf.Close()

  logger.Init("ojos", true, true, lf).Close()

  // Create or open the database
  err = database.Open()
  if err != nil {
    logger.Error(err)
    return
  }
  defer database.Close()

  // Start service
  err = http.ListenAndServe(":3000", rtr)
  if err != nil {
    logger.Fatal(err)
  }
}
