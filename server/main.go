package main

import(
  "os"
  "strconv"

  "github.com/hectorj2f/search_networking/resources"
  "github.com/hectorj2f/search_networking/database"
  "github.com/hectorj2f/search_networking/networking"
  logger "github.com/Sirupsen/logrus"
  )

func main() {
  db, err := database.SetupConnection()
  if err != nil {
    logger.Error(err)
    os.Exit(2)
  }

  if err = db.PopulateDb(); err != nil {
    logger.Error(err)
    os.Exit(2)
  }

  defer db.Close()

  cert := os.Getenv(resources.TLS_CERT_FLAG)
  key := os.Getenv(resources.TLS_KEY_FLAG)
  port := resources.SERVER_PORT
  if os.Getenv(resources.PORT_FLAG) != "" {
    port, _ = strconv.Atoi(os.Getenv(resources.PORT_FLAG))
  }

  networking.Server(cert, key, port)
}
