package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type User struct {
	GoogleID string `json:"google_id"`
	Name     string `json:"name"`
}

type Meow struct {
	ID       uint64 `json:"id"`
	UserName string `json:"user_name"`
	Body     string `json:"body"`
}

type Timeline struct {
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
	MeowID    uint64    `json:"meow_id"`
}

func (m *Meow) MarshalJSON() ([]byte, error) {
	type Alias Meow
	return json.Marshal(&struct {
		ID string `json:"id"`
		*Alias
	}{
		ID:    fmt.Sprint(m.ID),
		Alias: (*Alias)(m),
	})
}

func (m *Meow) UnmarshalJSON(data []byte) (err error) {
	type Alias Meow
	obj := &struct {
		ID string `json:"id"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		return err
	}
	m.ID, err = strconv.ParseUint(obj.ID, 10, 64)
	return err
}

func (t *Timeline) MarshalJSON() ([]byte, error) {
	type Alias Timeline
	return json.Marshal(&struct {
		MeowID string `json:"meow_id"`
		*Alias
	}{
		MeowID: fmt.Sprint(t.MeowID),
		Alias:  (*Alias)(t),
	})
}

func (t *Timeline) UnmarshalJSON(data []byte) (err error) {
	type Alias Timeline
	obj := &struct {
		MeowID string `json:"meow_id"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err = json.Unmarshal(data, &obj); err != nil {
		return err
	}
	t.MeowID, err = strconv.ParseUint(obj.MeowID, 10, 64)
	return err
}
