package api

import (
    "github.com/hectorj2f/search_networking/database"

    logger "github.com/Sirupsen/logrus"
)

func RegisterRequest(state string, id string) {
  db := database.GetDatabase()

  if err := database.RegisterRequest(db, state, id); err != nil {
    logger.Errorf("Error updating request %s", id)
  }
}

func UpdateRequest(state string, id string) {
  db := database.GetDatabase()

  if err := database.UpdateRequest(db, state, id); err != nil {
    logger.Errorf("Error updating request %s", id)
  }
}
