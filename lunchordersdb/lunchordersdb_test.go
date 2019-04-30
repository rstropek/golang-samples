package lunchordersdb

import (
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/dgraph-io/badger"
)

func runInDatabase(body func(LunchDB) error) error {
	err := os.RemoveAll("./db")
	defer os.RemoveAll("./db")

	db, err := Open("./db")
	if err != nil {
		return err
	}

	defer db.Close()

	if body != nil {
		return body(db)
	}

	return nil
}

func handle(err error, t *testing.T) {
	if err != nil {
		t.Errorf("Error while accessing DB: %s", err.Error())
	}
}

func expect(expected error, err error, t *testing.T) {
	switch {
	case err == nil:
		t.Errorf("Expected %s, got nil", expected.Error())
	case err != expected:
		t.Errorf("Expected %s, got %s", expected.Error(), err.Error())
	}
}

func TestOpenClose(t *testing.T) {
	handle(runInDatabase(func(db LunchDB) error { return nil }), t)
}

func TestAddGetItem(t *testing.T) {
	handle(runInDatabase(func(db LunchDB) error {
		t := NewTable(1)
		p := Person{Firstname: "Foo", Lastname: "Bar"}
		db.AddItem(t, &p)
		err := db.GetItem(t, p.ID, &p)
		switch {
		case err != nil:
			return err
		case p.Firstname != "Foo" || p.Lastname != "Bar":
			return errors.New("Did not return proper result")
		default:
			return nil
		}
	}), t)
}

func TestAddIterateItem(t *testing.T) {
	handle(runInDatabase(func(db LunchDB) error {
		t := NewTable(1)
		db.AddItem(t, &Person{Firstname: "Foo", Lastname: "Bar"})
		db.AddItem(t, &Person{Firstname: "John", Lastname: "Doe"})
		counter := 0
		err := db.IterateItems(t, func(value []byte) error {
			counter++
			var person Person
			err := json.Unmarshal(value, &person)
			if err != nil {
				return err
			}

			if person.Firstname != "Foo" && person.Firstname != "John" {
				return errors.New("Invalid data")
			}

			return nil
		})
		switch {
		case err != nil:
			return err
		case counter != 2:
			return errors.New("Wrong number of results")
		default:
			return nil
		}
	}), t)
}
func TestGetItemNotFound(t *testing.T) {
	expect(badger.ErrKeyNotFound, runInDatabase(func(db LunchDB) error {
		var p Person
		err := db.GetItem(NewTable(1), 1, &p)
		return err
	}), t)
}

func TestDeletePerson(t *testing.T) {
	handle(runInDatabase(func(db LunchDB) error {
		p := Person{Firstname: "Foo", Lastname: "Bar"}
		t := NewTable(1)
		db.AddItem(t, &p)
		err := db.DeleteItem(t, p.ID)
		if err != nil {
			return err
		}

		err = db.GetItem(t, p.ID, &p)
		if err != badger.ErrKeyNotFound {
			return errors.New("Expected key not found after delete, didn't get error")
		}

		return nil
	}), t)
}

func TestUpdatePerson(t *testing.T) {
	handle(runInDatabase(func(db LunchDB) error {
		p := Person{Firstname: "Foo", Lastname: "Bar"}
		t := NewTable(1)
		db.AddItem(t, &p)
		p.Firstname = "John"
		p.Lastname = "Doe"
		return db.UpdateItem(t, p)
	}), t)
}

func TestUpdatePersonNotFound(t *testing.T) {
	expect(badger.ErrKeyNotFound, runInDatabase(func(db LunchDB) error {
		return db.UpdateItem(NewTable(1), Person{})
	}), t)
}
