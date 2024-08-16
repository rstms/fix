package util

import (
	"fmt"
	"testing"
)

func TestConfirmYesLower(t *testing.T) {
	answer, err := Confirm("hit 'y'")
	switch {
	case err != nil:
		t.Errorf("Confirm failed: %v", err)
	case answer != true:
		t.Errorf("Confirm returned false; want true")
	}
}

func TestConfirmYesUpper(t *testing.T) {
	answer, err := Confirm("hit 'Y'")
	switch {
	case err != nil:
		t.Errorf("Confirm failed: %v", err)
	case answer != true:
		t.Errorf("Confirm returned false; want true")
	}
}

func TestConfirmYesDefault(t *testing.T) {

	answer, err := Confirm("hit enter")
	switch {
	case err != nil:
		t.Errorf("Confirm failed: %v", err)
	case answer != true:
		t.Errorf("Confirm returned false; want true")
	}
}

func TestConfirmNoUpper(t *testing.T) {
	fmt.Println("type 'N'")
	answer, err := Confirm("hit 'N'")
	switch {
	case err != nil:
		t.Errorf("Confirm failed: %v", err)
	case answer != false:
		t.Errorf("Confirm returned true; want false")
	}
}

func TestConfirmNoLower(t *testing.T) {
	fmt.Println("type 'n'")
	answer, err := Confirm("hit 'n'")
	switch {
	case err != nil:
		t.Errorf("Confirm failed: %v", err)
	case answer != false:
		t.Errorf("Confirm returned true; want false")
	}
}
