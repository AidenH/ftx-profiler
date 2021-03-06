package main

import (
	"fmt"
	"log"
	"os"
)

func HandleOsArgs() error {

	args := os.Args

	for i, item := range args {

		if item == "-t" || item == "--tape-only" {
			State.ProfileTrue = false
			log.Println("tape-only enabled")

		} else if item == "-p" || item == "--profile-only" {
			State.TapeTrue = false
			log.Println("profile-only enabled")

		} else if item == "-v" || item == "--volume-counts" {
			State.VolumeCounts = true
			log.Println("starting with profile volume counts enabled")

		} else if item == "-l" || item == "--load-session" {
			filename := args[i+1]

			if err := VolRead(filename); err != nil {
				return err
			}

			FileWrite(fmt.Sprint("volume data saved to: ", filename))

		}
	}

	return nil
}
