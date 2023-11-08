package helper

import log "github.com/sirupsen/logrus"

// @project photo-studio
// @created 06.02.2023
// @author arthur

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
