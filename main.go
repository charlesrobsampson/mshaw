package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

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
		letters := []letter{}
		err := json.Unmarshal(shavianJson, &letters)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		printCodeTable(letters)
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
		case "t", "table":
			letters := []letter{}
			err := json.Unmarshal(shavianJson, &letters)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			printShavianTable(letters)
		case "u", "unicode":
			fmt.Println("unicode prefix: U+104")
			letters := []letter{}
			err := json.Unmarshal(shavianJson, &letters)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			printCodeTable(letters)
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

func printShavianTable(shav []letter) {
	sort.Slice(shav, func(i, j int) bool {
		return shav[i].Code < shav[j].Code
	})
	// Initialize tabwriter with os.Stdout as the output, along with some formatting options
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "Letter\tName\tSound\tCode\tType\t") // The header
	fmt.Fprintln(w, "------\t----\t-----\t----\t----\t") // Header underline

	for _, letter := range shav {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", letter.Letter, letter.Name, letter.Sound, letter.Code, letter.Type)
	}

	w.Flush() // Flush writes to the underlying io.Writer
}

func printCodeTable(shav []letter) {
	sort.Slice(shav, func(i, j int) bool {
		return shav[i].Code < shav[j].Code
	})

	codes := map[string]map[string]string{
		"5": {},
		"6": {},
		"7": {},
	}
	for _, letter := range shav {
		parts := strings.Split(letter.Code, "")
		codes[parts[0]][parts[1]] = letter.Letter
	}

	// Initialize tabwriter with os.Stdout as the output, along with some formatting options
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, " \t 0\t 1\t 2\t 3\t 4\t 5\t 6\t 7\t 8\t 9\t A\t B\t C\t D\t E\t F\t") // The header
	fmt.Fprintln(w, " \t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t -\t") // Header underline

	fmt.Fprintf(w, "5\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", codes["5"]["0"], codes["5"]["1"], codes["5"]["2"], codes["5"]["3"], codes["5"]["4"], codes["5"]["5"], codes["5"]["6"], codes["5"]["7"], codes["5"]["8"], codes["5"]["9"], codes["5"]["A"], codes["5"]["B"], codes["5"]["C"], codes["5"]["D"], codes["5"]["E"], codes["5"]["F"])
	fmt.Fprintf(w, "6\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", codes["6"]["0"], codes["6"]["1"], codes["6"]["2"], codes["6"]["3"], codes["6"]["4"], codes["6"]["5"], codes["6"]["6"], codes["6"]["7"], codes["6"]["8"], codes["6"]["9"], codes["6"]["A"], codes["6"]["B"], codes["6"]["C"], codes["6"]["D"], codes["6"]["E"], codes["6"]["F"])
	fmt.Fprintf(w, "7\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n", codes["7"]["0"], codes["7"]["1"], codes["7"]["2"], codes["7"]["3"], codes["7"]["4"], codes["7"]["5"], codes["7"]["6"], codes["7"]["7"], codes["7"]["8"], codes["7"]["9"], codes["7"]["A"], codes["7"]["B"], codes["7"]["C"], codes["7"]["D"], codes["7"]["E"], codes["7"]["F"])
	fmt.Println()

	w.Flush() // Flush writes to the underlying io.Writer
}
