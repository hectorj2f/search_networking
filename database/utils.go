package database

import (
  "fmt"
  "time"

  "github.com/jmcvetta/randutil"
  logger "github.com/Sirupsen/logrus"
)

const(
  AMOUNT_USERS = 75
  )

func makeUsersTable(db *DB) error {
  _, err := db.Conn.Exec("DROP TABLE users")
  if err != nil {
      return err
  }
  _, err = db.Conn.Exec("CREATE TABLE users (id serial primary key, created timestamp without time zone, username text NOT NULL, role text NOT NULL, organization text NOT NULL, password text NOT NULL)")
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

func (db *DB) PopulateDb() error {
  logger.Debug("Populating the database...")
  makeUsersTable(db)

  for i := 0; i < AMOUNT_USERS; i++ {
      alias, _ := randutil.AlphaString(5)
      username := fmt.Sprintf("%s%s", choiceString([]string{"Jerry", "Johh", "Richard", "Brown", "Hector", "Stephan", "Anna"}), alias)
      role := choiceString([]string{"operations", "admin", "developer", "user"})
      organization := choiceString([]string{"elasticbox", "giantswarm", "google", "apple"})

      _, err := db.Conn.Exec("INSERT INTO users (created, username, role, organization, password) VALUES ($1, $2, $3, $4, 'pass')", time.Now(), username, role, organization)
      if err != nil {
          return err
      }
  }

  logger.Debugf("Database has %d users.", AMOUNT_USERS)
  return nil
}
