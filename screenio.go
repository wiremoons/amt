// amt - program to access an SQLite database and lookup acronyms
//
// author:	Simon Rowe <simon@wiremoons.com>
// license: open-source released under The MIT License (MIT).
//
// Package used to display output for application 'amt'

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"text/template"
)

// getInput function asks the user a question and returns their
// answer. The question is provided to the function as a string
// 'question' and the users response is returned by the function as a
// string 'response'.
func getInput(question string) string {
	if debugSwitch {
		fmt.Println("\nDEBUG: in function 'getInput' ...")
	}
	// create a new reader and attached to stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the user the question passed to the function
	fmt.Printf("%s", question)
	// read the user's response - terminating their input on newline
	response, _ := reader.ReadString('\n')
	if debugSwitch {
		fmt.Printf("\nDEBUG: user provided input: '%s' \n", response)
	}
	// remove the trailing newline (Unix/Mac) or both the newline and
	// return (Windows) from the input string provided by the user. As
	// the ReadString() function keeps any line suffix ('\n' or
	// '\r\n)' when it returns, both are required to ensure options
	// are covered. If 'response' doesn't end with either suffix, it
	// is returned unchanged - so no harm done!
	response = strings.TrimSuffix(response, "\n")
	response = strings.TrimSuffix(response, "\r")
	if debugSwitch {
		fmt.Printf("\nDEBUG: user provided input (after TrimSuffix): '%s' \n", response)
	}
	// flush any output to the screen
	_ = os.Stdout.Sync()
	// return the string read from the user to the calling function
	return response
}

// checkContinue function asks the user if they would like to continue
// with the currently running part of the application.
//
// checkContinue function reads input from the users console to see if
// they provide a 'y' or 'n' response.
//
// The function returns a bool depending on the user's response.
// If the response contains the letter 'y' it returns 'true'. Any other
// response will return 'false'.
func checkContinue() bool {
	// create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the user if they wish to continue
	fmt.Print("Continue? [y/n]: ")
	// read the user's response - terminating their input on newline
	response, _ := reader.ReadString('\n')
	// convert the response to lower case - easier to compare
	response = strings.ToLower(response)
	// see if the user input contains 'y'
	if strings.Contains(response, "y") {
		// done here - so return
		return true
	}
	// if above failed - so return false
	return false
}

// printBanner function is used to print out a small program banner
// which displays the application name.
func printBanner() {
	fmt.Println("\n\t\t\tAcronym Management Tool 'amt'")
	fmt.Println("\t\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯")
}

// versionInfo function collects details of the program being run and
// displays it on stdout
func versionInfo() {
	// define a template for display on screen with placeholders for data
	const appInfoTmpl = `
Running '{{.appname}}' version {{.appversion}}

 - Built with Go Compiler '{{.compiler}}' on Golang version '{{.version}}'
 - Author's web site: https://www.wiremoons.com/
 - Source code for {{.appname}}: https://github.com/wiremoons/amt-go/

`
	// build a map with keys set to match the template names used
	// and the data fields to be used in the template as values
	data := map[string]interface{}{
		"appname":    appname,
		"appversion": appversion,
		"compiler":   runtime.Compiler,
		"version":    runtime.Version(),
	}
	// check and build the template so the data field values are added
	// and the final output is displayed. Check for any error, and
	// abort if one is found.
	t := template.Must(template.New("appInfo").Parse(appInfoTmpl))
	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatalf("FATAL ERROR: in function 'versionInfo()' when building template with err: %v", err)
	}
}

// myUsage function replaces the standard flag.Usage() function from Go. The
// function takes no parameters, but outputs the command line flags
// that can be used when running the program.
func myUsage() {
	usageText := `
Usage of ./amt:

        Flag               Description                                        Default Value
        ¯¯¯¯               ¯¯¯¯¯¯¯¯¯¯¯                                        ¯¯¯¯¯¯¯¯¯¯¯¯¯
        -d                 show debug output                                  false
        -f <filename>      provide filename and path to SQLite database       optional
        -h                 display help for this program                      false
        -n                 add a new acronym record                           optional
        -s <acronym>       provide acronym to search for                      optional
        -r <acronym id>    provide acronym id to remove                       optional
        -v                 display program version                            false
        -w                 search for any similar matches                     false`
	fmt.Println(usageText)

}
