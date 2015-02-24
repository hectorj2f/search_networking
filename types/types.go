package types

import (
  "time"
  
  "code.google.com/p/go-uuid/uuid"
  "github.com/docker/libchan"
  )

type Message struct {
  Data    []byte
  Id      uuid.UUID
  Ret     libchan.Sender
}

type AckMessage struct {
  Id      uuid.UUID
  Data    []byte
}

type User struct {
    Id  					int
    Created 			*time.Time
    //Data *json.RawMessage
    Username 			string
    Role					string
    Organization 	string
    Password 			string
}
