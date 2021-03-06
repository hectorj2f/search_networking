package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/hectorj2f/search_networking/types"

	logger "github.com/Sirupsen/logrus"
)

func RegisterRequest(db *DB, state string, id string) (error) {
	// TODO: Store remote IP making the request
	_, err := db.Conn.Exec("INSERT INTO requests (message_id, created, state) VALUES ($1, $2, $3)", id, time.Now(), state)
	return err
}
func UpdateRequest(db *DB, state string, id string) (error) {
	// TODO: Lock
	_, err := db.Conn.Exec("UPDATE requests SET state = $1, updated = $2 WHERE message_id = $3", state, time.Now(), id)
	// TODO: Unlock the transaction
	return err
}

func GetUsers(db *DB) ([]map[string]interface{}, error) {
	result, err := executeQuery(db, "SELECT id, created, username, role, organization, password FROM users")
	return result, err
}

func GetUsersByRole(db *DB, role string) ([]map[string]interface{}, error) {
	result, err := executeQuery(db, "SELECT id, created, username, role, organization, password FROM users WHERE role = $1", role)
	return result, err
}

func GetUsersById(db *DB, id string) ([]map[string]interface{}, error) {
	result, err := executeQuery(db, "SELECT id, created, username, role, organization, password FROM users WHERE id = $1", id)
	return result, err
}

func GetUsersByUsername(db *DB, username string) ([]map[string]interface{}, error) {
	result, err := executeQuery(db, "SELECT id, created, username, role, organization, password FROM users WHERE username = $1", username)
	return result, err
}

func GetUsersByOrganization(db *DB, organization string) ([]map[string]interface{}, error) {
	result, err := executeQuery(db, "SELECT id, created, username, role, organization, password FROM users WHERE organization = $1", organization)
	return result, err
}

func GetUsersByMultipleCriteria(db *DB, query map[string]interface{}) ([]map[string]interface{}, error) {
	where := ""
	for key, value := range query {
		if len(where) > 0 {
			if key == "username" {
				where = fmt.Sprintf("%s and %s = '%%s%'", where, key, value)
				continue
			}
			where = fmt.Sprintf("%s and %s = '%s'", where, key, value)
		} else {
			if key == "username" {
				where = fmt.Sprintf("%s = '%%s%'", where, key, value)
				continue
			}
			where = fmt.Sprintf("%s = '%s'", key, value)
		}
	}
	criteria := fmt.Sprintf("SELECT id, created, username, role, organization, password FROM users WHERE %s", where)
	logger.Debugf("SQL query %s", criteria)
	result, err := executeQuery(db, criteria)
	return result, err
}


func executeQuery(db *DB, criteria string, args ...interface{}) ([]map[string]interface{}, error) {
	var (
		err error
		rows *sql.Rows
	)
	if len(args) == 0 {
		rows, err = db.Conn.Query(criteria)
	} else {
		rows, err = db.Conn.Query(criteria, args[0].(string))
	}

	if err != nil {
			return nil, err
	}
	users := make([]types.User, 0)
	for rows.Next() {
			var u types.User
			err := rows.Scan(
					&u.Id,
					&u.Created,
					&u.Username,
					&u.Role,
					&u.Organization,
					&u.Password,
			)
			if err != nil {
					return nil, err
			}
			users = append(users, u)
	}

	result := prepareUsersOutput(users)

	return result, nil
}


func prepareUsersOutput(users []types.User) ( []map[string]interface{} ) {
		logger.Infof("Found users: %s", strconv.Itoa(len(users)))
		list_users := make([]map[string]interface{}, 0)
		for _, user := range users {
			json := map[string]interface{}{"id": user.Id,
																		"created": fmt.Sprintf("%s",user.Created),
																		"username": user.Username,
																		"role": user.Role,
																		"organization": user.Organization}
			list_users = append(list_users, json)
		}
		return list_users
}
