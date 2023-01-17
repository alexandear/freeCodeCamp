package main

import (
	"bufio"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

const passwordNotInDatabase = "PASSWORD NOT IN DATABASE"

func CrackSHA1Hash(passwordHash string, ifUseSalts bool) string {
	passwordFile, err := os.Open("top-10000-passwords.txt")
	if err != nil {
		log.Fatalf("Failed to open file with passwords: %v\n", err)
	}
	defer func() { _ = passwordFile.Close() }()

	if !ifUseSalts {
		if password := passwordFromReader(passwordFile, passwordHash, nil); password != nil {
			return string(password)
		}

		return passwordNotInDatabase
	}

	saltFile, err := os.Open("known-salts.txt")
	if err != nil {
		log.Fatalf("Failed to open file with salts: %v\n", err)
	}
	defer func() { _ = saltFile.Close() }()

	salts := make([][]byte, 0, 20)
	for scanner := bufio.NewReader(saltFile); ; {
		salt, _, err := scanner.ReadLine()
		if errors.Is(io.EOF, err) {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read salt line: %v\n", err)
		}
		salts = append(salts, salt)
	}

	if password := passwordFromReader(passwordFile, passwordHash, salts); password != nil {
		return string(password)
	}

	return passwordNotInDatabase
}

func passwordFromReader(passwordReader io.Reader, passwordHash string, salts [][]byte) []byte {
	for scanner := bufio.NewReader(passwordReader); ; {
		password, _, err := scanner.ReadLine()
		if errors.Is(io.EOF, err) {
			return nil
		}
		if err != nil {
			log.Fatalf("Failed to read password line: %v\n", err)
		}

		if salts == nil {
			if passwordHash == hashSHA1(password) {
				return password
			}
			continue
		}

		for _, salt := range salts {
			sp := append(append([]byte{}, salt...), password...)
			if passwordHash == hashSHA1(sp) {
				return password
			}

			ps := append(append([]byte{}, password...), salt...)
			if passwordHash == hashSHA1(ps) {
				return password
			}
		}
	}
}

func hashSHA1(password []byte) string {
	h := sha1.New()
	h.Write(password)
	return fmt.Sprintf("%x", h.Sum(nil))
}
