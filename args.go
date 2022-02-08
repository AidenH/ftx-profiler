package main

import (
	"os"
)

func HandleOsArgs() {

	args := os.Args

	for _, item := range args {

		if item == "tape-only" {

			State.ProfileTrue = false

		} else if item == "profile-only" {

			State.TapeTrue = false

		} else if item == "volume-counts" {

			State.VolumeCounts = true

		}
	}
}
