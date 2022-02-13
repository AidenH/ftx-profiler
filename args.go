package main

import (
	"log"
	"os"
)

func HandleOsArgs() {

	args := os.Args

	for _, item := range args {

		if item == "tape-only" {
			State.ProfileTrue = false
			log.Println("tape-only enabled")

		} else if item == "profile-only" {
			State.TapeTrue = false
			log.Println("profile-only enabled")

		} else if item == "volume-counts" {
			State.VolumeCounts = true
			log.Println("starting with profile volume counts enabled")

		}
	}
}
