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
		if password := checkPasswordsInFile(passwordFile, passwordHash, nil); password != "" {
			return password
		}

		return passwordNotInDatabase
	}

	saltFile, err := os.Open("known-salts.txt")
	if err != nil {
		log.Fatalf("Failed to open file with salts: %v\n", err)
	}
	defer func() { _ = saltFile.Close() }()

	salts := make([]string, 0, 20)
	for scanner := bufio.NewReader(saltFile); ; {
		salt, _, err := scanner.ReadLine()
		if errors.Is(io.EOF, err) {
			break
		}
		if err != nil {
			log.Fatalf("Failed to read salt line: %v\n", err)
		}
		salts = append(salts, string(salt))
	}

	if password := checkPasswordsInFile(passwordFile, passwordHash, salts); password != "" {
		return password
	}

	return passwordNotInDatabase
}

func checkPasswordsInFile(passwordFile io.Reader, passwordHash string, salts []string) string {
	for scanner := bufio.NewReader(passwordFile); ; {
		line, _, err := scanner.ReadLine()
		if errors.Is(io.EOF, err) {
			return ""
		}
		if err != nil {
			log.Fatalf("Failed to read password line: %v\n", err)
		}
		password := string(line)

		if salts == nil {
			if passwordHash == hashSHA1(password) {
				return password
			}
			continue
		}

		for _, salt := range salts {
			if passwordHash == hashSHA1(salt+password) {
				return password
			}

			if passwordHash == hashSHA1(password+salt) {
				return password
			}
		}
	}
}

func hashSHA1(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", h.Sum(nil))
}
