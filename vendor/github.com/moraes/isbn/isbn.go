// Copyright 2015 Rodrigo Moraes. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package isbn provides functions to validate ISBN strings, calculate ISBN
check digits and convert ISBN-10 to ISBN-13.
*/
package isbn

import (
	"fmt"
	"strconv"
)

// sum10 returns the weighted sum of the provided ISBN-10 string. It is used
// to calculate the ISBN-10 check digit or to validate an ISBN-10.
//
// The provided string must have a length of 9 or 10 and no formatting
// characters (spaces or hyphens).
func sum10(isbn string) (int, error) {
	s := 0
	w := 10
	for k, v := range isbn {
		if k == 9 && v == 88 {
			// Handle "X" as the digit.
			s += 10
		} else {
			n, err := strconv.Atoi(string(v))
			if err != nil {
				return -1, fmt.Errorf("Failed to convert ISBN-10 character to int: %s", string(v))
			}
			s += n * w
		}
		w--
	}
	return s, nil
}

// sum13 returns the weighted sum of the provided ISBN-13 string. It is used
// to calculate the ISBN-13 check digit or to validate an ISBN-13.
//
// The provided string must have a length of 12 or 13 and no formatting
// characters (spaces or hyphens).
func sum13(isbn string) (int, error) {
	s := 0
	w := 1
	for _, v := range isbn {
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return -1, fmt.Errorf("Failed to convert ISBN-13 character to int: %s", string(v))
		}
		s += n * w
		if w == 1 {
			w = 3
		} else {
			w = 1
		}
	}
	return s, nil
}

// CheckDigit10 returns the check digit for an ISBN-10.
//
// The provided string must have a length of 9 or 10 and no formatting
// characters (spaces or hyphens). For a 10-length string, the last character
// (the digit) is ignored since that is what is being (re)calculated.
func CheckDigit10(isbn10 string) (string, error) {
	if len(isbn10) != 9 && len(isbn10) != 10 {
		return "", fmt.Errorf("A string of length 9 or 10 is required to calculate the ISBN-10 check digit. Provided was: %s", isbn10)
	}
	s, err := sum10(isbn10[:9])
	if err != nil {
		return "", err
	}
	d := (11 - (s % 11)) % 11
	if d == 10 {
		return "X", nil
	}
	return strconv.Itoa(d), nil
}

// CheckDigit13 returns the check digit for an ISBN-13.
//
// The provided string must have a length of 12 or 13 and no formatting
// characters (spaces or hyphens). For a 13-length string, the last character
// (the digit) is ignored since that is what is being (re)calculated.
func CheckDigit13(isbn13 string) (string, error) {
	if len(isbn13) != 12 && len(isbn13) != 13 {
		return "", fmt.Errorf("A string of length 12 or 13 is required to calculate the ISBN-13 check digit. Provided was: %s", isbn13)
	}
	s, err := sum13(isbn13[:12])
	if err != nil {
		return "", err
	}
	d := 10 - (s % 10)
	if d == 10 {
		return "0", nil
	}
	return strconv.Itoa(d), nil
}

// Validate returns true if the provided string is a valid ISBN-10 or ISBN-13.
//
// The provided string must have a length of 10 or 13 and no formatting
// characters (spaces or hyphens).
func Validate(isbn string) bool {
	switch len(isbn) {
	case 10:
		return Validate10(isbn)
	case 13:
		return Validate13(isbn)
	}
	return false
}

// Validate10 returns true if the provided string is a valid ISBN-10.
//
// The provided string must have a length of 10 and no formatting
// characters (spaces or hyphens).
func Validate10(isbn10 string) bool {
	if len(isbn10) == 10 {
		s, _ := sum10(isbn10)
		return s%11 == 0
	}
	return false
}

// Validate13 returns true if the provided string is a valid ISBN-13.
//
// The provided string must have a length of 13 and no formatting
// characters (spaces or hyphens).
func Validate13(isbn13 string) bool {
	if len(isbn13) == 13 {
		s, _ := sum13(isbn13)
		return s%10 == 0
	}
	return false
}

// To13 converts an ISBN-10 to an ISBN-13.
//
// The provided string must have a length of 9 or 10 and no formatting
// characters (spaces or hyphens).
func To13(isbn10 string) (string, error) {
	if len(isbn10) != 9 && len(isbn10) != 10 {
		return "", fmt.Errorf("A string of length 9 or 10 is required to convert an ISBN-10 to an ISBN-13. Provided was: %s", isbn10)
	}
	isbn13 := "978" + isbn10[:9]
	d, err := CheckDigit13(isbn13)
	if err != nil {
		return "", err
	}
	return isbn13 + d, nil
}
