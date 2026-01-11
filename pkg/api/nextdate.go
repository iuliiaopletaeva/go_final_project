package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	var interval int

	if dstart == "" {
		return "", errors.New("dstart cannot be empty")
	}

	if repeat == "" {
		return "", errors.New("repeat cannot be empty")
	}

	date, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", errors.New("cannot parse date of start")
	}

	repeatParts := strings.Split(repeat, " ")

	dateType := string(repeatParts[0])

	if dateType == "d" {
		if len(repeatParts) < 2 {
			return "", errors.New("interval for days is not specified")
		}

		interval, err = strconv.Atoi(repeatParts[1])
		if err != nil {
			return "", errors.New("cannot convert days count to int")
		}

		if interval == 0 {
			return "", errors.New("interval for days is not specified")
		}

		if interval > 400 {
			return "", errors.New("interval cannot be more than 400 days")
		}

		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}

	}

	if dateType == "y" {
		if len(repeatParts) > 1 {
			return "", errors.New("for year interval cannot be numbers")
		}

		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
	}

	if repeatParts[0] != "d" && repeatParts[0] != "y" {
		return "", errors.New("repeat interval is not corrected")
	}

	return date.Format(dateFormat), nil
}

func afterNow(date, now time.Time) bool {
	return date.Format(dateFormat) > now.Format(dateFormat)
}
