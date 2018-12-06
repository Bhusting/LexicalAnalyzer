package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/scanner"
)

func main() {
	stream, err := ioutil.ReadFile(os.Args[1])
	fmt.Println("Processing Input File: " + os.Args[1])
	if err != nil {
		fmt.Println("Error")
	}

	fileStr := string(stream)

	var s scanner.Scanner
	s.Init(strings.NewReader(fileStr))

	count := 0

	file, err := os.Create("test02.txt")

	if err != nil {
		log.Fatal(err)
	}

	var tokens []Token
	var temp Token
	var tempPos int
	update := false
	var errPos []int
	lexError := false

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		var buffer = s.TokenText()

		if buffer == "$" {
			s.Scan()
			file.WriteString("ID[STRING]: " + s.TokenText() + "\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer + s.TokenText()
			} else {
				tempPos = count
			}
			temp = Token{s.TokenText(), "string", "", "ID"}
		} else if buffer == "#" {
			s.Scan()
			file.WriteString("ID[INT]: " + s.TokenText() + "\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer + s.TokenText()
			} else {
				tempPos = count
			}
			temp = Token{s.TokenText(), "int", "", "ID"}
		} else if buffer == "%" {
			s.Scan()
			file.WriteString("ID[REAL]: " + s.TokenText() + "\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer + s.TokenText()
			} else {
				tempPos = count
			}
			temp = Token{s.TokenText(), "real", "", "ID"}
		} else if buffer == "<" {
			s.Scan()
			if s.TokenText() == "=" {
				file.WriteString("ASSIGN\n")
				update = true
				temp = Token{"<=", "", "", "Assign"}
			}
		} else if buffer[0] == '"' {
			buffer = buffer[1:]
			i := len(buffer) - 1
			buffer = buffer[:i]
			file.WriteString("STRING: " + buffer + "\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer
			}
			temp = Token{buffer, "string", buffer, "STRING"}
		} else if n, err := strconv.ParseFloat(buffer, 64); err == nil {
			if n == float64(int64(n)) {
				file.WriteString("INT_CONST: " + buffer + "\n")
				if update {
					tokens[tempPos].Value = tokens[tempPos].Value + buffer
				}
				temp = Token{buffer, "int", buffer, "INT_CONST"}
			} else {
				file.WriteString("REAL_CONST: " + buffer + "\n")
				if update {
					tokens[tempPos].Value = tokens[tempPos].Value + buffer
				}
				temp = Token{buffer, "real", buffer, "REAL_CONST"}
			}
		} else if buffer == ":" {
			file.WriteString("COLON\n")
			update = false
			temp = Token{s.TokenText(), "", "", "COLON"}
		} else if buffer == "(" {
			file.WriteString("LPAREN\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer
			}
			temp = Token{s.TokenText(), "", "", "LPAREN"}
		} else if buffer == ")" {
			file.WriteString("RPAREN\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer
			}
			temp = Token{s.TokenText(), "", "", "RPAREN"}
		} else if buffer == "+" {
			file.WriteString("PLUS\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer
			}
			temp = Token{s.TokenText(), "", "", "PLUS"}
		} else if buffer == "-" {
			file.WriteString("MINUS\n")
			temp = Token{s.TokenText(), "", "", "MINUS"}
		} else if buffer == "/" {
			file.WriteString("DIVISION\n")
			temp = Token{s.TokenText(), "", "", "DIVISION"}
		} else if buffer == "*" {
			file.WriteString("TIMES\n")
			temp = Token{s.TokenText(), "", "", "TIMES"}
		} else if buffer == "^" {
			file.WriteString("POWER\n")
			if update {
				tokens[tempPos].Value = tokens[tempPos].Value + buffer
			}
			temp = Token{s.TokenText(), "", "", "POWER"}
		} else if buffer == "WRITE" {
			file.WriteString("WRITE\n")
			temp = Token{s.TokenText(), "", "", "WRITE"}
		} else if buffer == "BEGIN" {
			file.WriteString("BEGIN\n")
			temp = Token{s.TokenText(), "", "", "BEGIN"}
		} else if buffer == "END" {
			file.WriteString("END\n")
			temp = Token{s.TokenText(), "", "", "END"}
		} else {
			lexError = true
			errPos = append(errPos, count)
			file.WriteString("Lexical Error: Unkown Token\n")
			temp = Token{"UNKOWN TOKEN", "UNKOWN TOKEN", "UNKOWN TOKEN", "UNKOWN TOKEN"}
		}
		tokens = append(tokens, temp)
		count++
	}

	PrintTokens(tokens)

	fmt.Print(count)
	fmt.Println(" Tokens produced")
	if lexError {
		for i := range errPos {
			fmt.Println("Lexical Error found at " + strconv.Itoa(errPos[i]))
		}
	}
	fmt.Println("Results in Output File: test.txt")
	fmt.Println("Extra Credit in Output File: testExtra.txt")
	// Close the file
	file.Close()
}

func PrintTokens(tokens []Token) {
	file, err := os.Create("testExtra.txt")

	if err != nil {
		log.Fatal(err)
	}

	for i := range tokens {
		file.WriteString("Token " + strconv.Itoa(i+1) + "\n")
		file.WriteString("\t" + "Lexeme: " + tokens[i].Lexeme + "\n")
		file.WriteString("\t" + "Type: " + tokens[i].Type + "\n")
		file.WriteString("\t" + "Value: " + tokens[i].Value + "\n")
		file.WriteString("\t" + "Token: " + tokens[i].Token + "\n\n")
	}
	file.Close()
}

type Token struct {
	Lexeme string
	Type   string
	Value  string
	Token  string
}
