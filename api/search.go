package api

import (
    "github.com/hectorj2f/search_networking/database"

    logger "github.com/Sirupsen/logrus"
)

func SearchUsers(query map[string]interface{}) ([]map[string]interface{}, error) {
  logger.Infof("Query to search users: %s", query)
  db := database.GetDatabase()

  if len(query) > 0 {
    /*if _, present := query["organization"]; present {
      return database.GetUsersByOrganization(db, query["organization"].(string))
    }
    */
    return database.GetUsersByMultipleCriteria(db, query)
  }
  return database.GetUsers(db)
}
