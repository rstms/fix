package util

import (
    "fmt"
    "os"
    //"github.com/eiannone/keyboard"
    "golang.org/x/term"
)

const DELAY = 100

func getKey() (rune, error) {
    var ret rune
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        return ret, fmt.Errorf("Error setting stdin to raw mode: %v\n", err)
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
    /*
    err := keyboard.Open()
    if err != nil {
	return false, err
    }

    defer func() {
        fmt.Printf("closing keyboard...")
        _ = keyboard.Close()
        fmt.Printf("keboard closed")
    }()
    
    for {
        fmt.Println("GetKey...")
	key, _, err := keyboard.GetKey()
	if err != nil {
            return false, err
        }
        fmt.Printf("GetKey returned key=%v err=%v\n", key, err)
    */
    key, err := getKey()
    if err != nil {
        return false, err
    }

    ret := false
    switch key {
        case 'y','Y', 0:
	    fmt.Println("yes")
            ret = true
	default:
	    fmt.Println("no")
	}
    return ret, nil
}
