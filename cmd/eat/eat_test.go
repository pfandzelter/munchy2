package main

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/pfandzelter/munchy2/pkg/kaiserstueck"
	"github.com/pfandzelter/munchy2/pkg/personalkantine"
	"github.com/pfandzelter/munchy2/pkg/singh"
	"github.com/pfandzelter/munchy2/pkg/stw"
)

const testLog string = "test.log"

func testCanteen(name string, specDiet bool, c mensa) error {

	fl, err := c.GetFood(time.Now())

	if err != nil {
		log.Print(err)
		return err
	}

	if len(fl) == 0 {
		log.Printf("No food found for %s", name)
		return fmt.Errorf("No food found for %s", name)
	}

	// append output to test.log
	f, err := os.OpenFile(testLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Print(err)
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Found %d items for %s", len(fl), name)
	fmt.Fprintf(f, "%s: %+v\n", name, fl)
	return nil
}

func TestHauptmensa(t *testing.T) {
	err := testCanteen("Hauptmensa", true, stw.New(321))

	if err != nil {
		t.Error(err)
	}
}

func TestVeggie(t *testing.T) {
	err := testCanteen("Pasteria Veggie 2.0", true, stw.New(631))

	if err != nil {
		t.Error(err)
	}
}

func TestArchitektur(t *testing.T) {
	err := testCanteen("Pastaria Architektur", true, stw.New(540))

	if err != nil {
		t.Error(err)
	}
}

func TestMarchstr(t *testing.T) {
	err := testCanteen("Marchstr", true, stw.New(538))

	if err != nil {
		t.Error(err)
	}
}

func TestKaiserstueck(t *testing.T) {
	err := testCanteen("Kaiserstück", false, kaiserstueck.New())

	if err != nil {
		t.Error(err)
	}
}

func TestPersonalkantine(t *testing.T) {
	err := testCanteen("Personalkantine", true, personalkantine.New())

	if err != nil {
		t.Error(err)
	}
}

func TestSingh(t *testing.T) {
	err := testCanteen("Mathe Café", true, singh.New())

	if err != nil {
		t.Error(err)
	}
}

func TestAll(t *testing.T) {
	TestHauptmensa(t)
	TestVeggie(t)
	TestKaiserstueck(t)
	TestPersonalkantine(t)
	TestSingh(t)
}
