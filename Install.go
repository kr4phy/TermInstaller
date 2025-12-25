package main

import "errors"

func Install() error {
	if false {
		return errors.New("errInstallFailed")
	} else {
		return nil
	}
}
