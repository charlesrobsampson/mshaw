package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"golang.design/x/clipboard"
)

type (
	letter struct {
		Letter string `json:"letter"`
		Name   string `json:"name"`
		Code   string `json:"code"`
		Sound  string `json:"sound"`
		Type   string `json:"type"`
	}

	Shavian map[string]letter
)

//go:embed shavian.json
var shavianJson json.RawMessage

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		spitShav()
	} else {
		switch args[0] {
		case "ls", "list":
			shavBytes, err := json.MarshalIndent(shavianJson, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Println(string(shavBytes))
		case "f", "find":
			args = args[1:]
			if len(args) < 2 {
				fmt.Println("shaw f [name|code] <value>")
				os.Exit(1)
			}
			var shav Shavian
			if args[0] == "name" {
				shav = buildShavianByName()
			} else if args[0] == "code" {
				shav = buildShavian()
			} else {
				fmt.Println("shaw f [name|code] <value>")
				os.Exit(1)
			}
			letter, ok := shav[args[1]]
			output, err := json.MarshalIndent(letter, "", "  ")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if ok {
				fmt.Println(string(output))
			} else {
				fmt.Println("no shavian found for", args[1])
			}
		case "h", "help":
			fmt.Println("shaw ([list|find|help])")
		default:
			fmt.Println("shaw ([list|find|help])")
		}
	}

}

func spitShav() {
	outstring := ""
	lines := getMultiLine("Enter the 2-digit codes for the shavian (50 - 7F):")
	shav := buildShavian()
	for _, line := range lines {
		words := strings.Fields(line)
		for _, word := range words {
			twoCharStrings := breakIntoTwoCharStrings(word)
			for _, str := range twoCharStrings {
				shavChar, ok := shav[strings.ToUpper(str)]
				if ok {
					outstring += shavChar.Letter
				} else {
					outstring += str
				}
			}
			outstring += " "
		}
		outstring += "\n"
	}
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	clipboard.Write(clipboard.FmtText, []byte(outstring))
	fmt.Println()
	fmt.Println(outstring)
	fmt.Println("copied to clipboard")
}

func getMultiLine(prompt string) []string {
	scn := bufio.NewScanner(os.Stdin)
	var lines []string
	fmt.Println(prompt)
	for scn.Scan() {
		line := scn.Text()
		if len(line) == 0 {
			break
		}
		if line[len(line)-1] != '\\' {
			if len(line) > 0 {
				line = strings.TrimSpace(line)
				lines = append(lines, line)
			}
			break
		} else {
			line = line[:len(line)-1]
			line = strings.TrimSpace(line)
			lines = append(lines, line)
		}
	}

	if err := scn.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return lines
}

func buildShavian() Shavian {
	shav := make(Shavian)
	letters := []letter{}
	err := json.Unmarshal(shavianJson, &letters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, l := range letters {
		shav[l.Code] = l
	}
	return shav
}

func buildShavianByName() Shavian {
	shav := make(Shavian)
	letters := []letter{}
	err := json.Unmarshal(shavianJson, &letters)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, l := range letters {
		shav[l.Name] = l
	}
	return shav
}

func breakIntoTwoCharStrings(word string) []string {
	var result []string
	for i := 0; i < len(word); i += 2 {
		if i+2 <= len(word) {
			result = append(result, word[i:i+2])
		} else {
			result = append(result, word[i:])
		}
	}
	return result
}
