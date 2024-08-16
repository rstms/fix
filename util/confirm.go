package util

import (
	"fmt"
	"golang.org/x/term"
	"os"
)

const DELAY = 100

func GetKey() (rune, error) {
	var ret rune
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return ret, fmt.Errorf("Failed setting stdin to raw mode: %v\n", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	buf := make([]byte, 1)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return ret, fmt.Errorf("Stdin.Read failed: %v\n", err)
	}
	return rune(buf[0]), nil
}

func Confirm(prompt string) (bool, error) {
	fmt.Printf("%s? [Y/n] ", prompt)
	key, err := GetKey()
	if err != nil {
		return false, err
	}

	ret := false
	switch key {
	case 'y', 'Y', '\r', '\n':
		fmt.Println("yes")
		ret = true
	default:
		fmt.Println("no")
	}
	return ret, nil
}
