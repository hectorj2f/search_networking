package main

import (
  "fmt"
  "time"
  "os"
  "strconv"

  "github.com/hectorj2f/search_networking/database"

  "github.com/jmcvetta/randutil"
  logger "github.com/Sirupsen/logrus"
)

const(
  AMOUNT_USERS = 500
  USERS_FLAG = "USERS"
  )


func main(){
  db, err := database.SetupConnection()
  if err != nil {
    logger.Error(err)
    os.Exit(2)
  }

  if err = populateDb(db); err != nil {
    logger.Error(err)
    os.Exit(2)
  }

  defer db.Close()
}

func makeUsersTable(db *database.DB) error {
  _, err := db.Conn.Exec("DROP TABLE IF EXISTS users")
  if err != nil {
      return err
  }
  _, err = db.Conn.Exec("DROP TABLE IF EXISTS requests")
  if err != nil {
      return err
  }
  _, err = db.Conn.Exec("CREATE TABLE users (id serial primary key, created timestamp without time zone, username text NOT NULL, role text NOT NULL, organization text NOT NULL, password text NOT NULL)")
  if err != nil {
      return err
  }

  _, err = db.Conn.Exec("CREATE TABLE requests (message_id text NOT NULL, created timestamp without time zone, updated timestamp without time zone, state text NOT NULL)")
  if err != nil {
      return err
  }

  return nil
}

func choiceString(choices []string) string {
  var winner string
  length := len(choices)
  i, err := randutil.IntRange(0, length)
  if err != nil {
    panic(err)
  }
  winner = choices[i]
  return winner
}

func populateDb(db *database.DB) error {
  logger.Info("Populating the database...")
  makeUsersTable(db)

  amout_users := AMOUNT_USERS
  if os.Getenv(USERS_FLAG) != "" {
    amout_users, _ = strconv.Atoi(os.Getenv(USERS_FLAG))
  }

  for i := 0; i < amout_users; i++ {
      alias, _ := randutil.AlphaString(5)
      username := fmt.Sprintf("%s%s", choiceString([]string{"Jerry", "Johh", "Richard", "Brown", "Hector", "Stephan", "Anna"}), alias)
      role := choiceString([]string{"operations", "admin", "developer", "user"})
      organization := choiceString([]string{"elasticbox", "giantswarm", "google", "apple"})

      _, err := db.Conn.Exec("INSERT INTO users (created, username, role, organization, password) VALUES ($1, $2, $3, $4, 'pass')", time.Now(), username, role, organization)
      if err != nil {
          return err
      }
  }

  logger.Infof("Database populated with %d users.", AMOUNT_USERS)
  return nil
}
