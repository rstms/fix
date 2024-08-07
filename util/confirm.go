package util

import (
    "fmt"
    "github.com/eiannone/keyboard"
    "time"
)

const DELAY = 100

func Confirm(prompt string) (bool, error) {
    fmt.Printf("%s? [Y/n] ", prompt)
    err := keyboard.Open()
    if err != nil {
	return false, err
    }
    defer keyboard.Close()
    for {
	key, _, err := keyboard.GetKey()
	if err != nil {
	    time.Sleep(DELAY * time.Millisecond)
	    continue
	}
	switch key {
	case 'y','Y', 0:
	    fmt.Println("yes")
	    return true, nil
	default:
	    fmt.Println("no")
	    return false, nil
	}
    }
    return false, nil
}
