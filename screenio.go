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

// getInput asks user a question and return their answer.
// The question is provided to the function as a string 'question' and
// the users response is returned as a string 'response'.
func getInput(question string) string {
	if debugSwitch {
		fmt.Println("\nDEBUG: in function 'getInput' ...")
	}
	// create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the user the question passed to the function
	fmt.Printf("%s", question)
	// read the users response - terminating their input on newline
	response, _ := reader.ReadString('\n')
	if debugSwitch {
		fmt.Printf("\nDEBUG: user provided input: '%s' \n", response)
	}
	// remove the trailing newline (Unix/Mac) or return and newline (Windows)
	// from the string provided by the user, as the ReadString() keeps any line
	// suffix (\n or \r\n) when it returns.
	// If 'response' doesn't end with either suffix, it is returned unchanged - so no harm done!
	response = strings.TrimSuffix(response, "\n")
	response = strings.TrimSuffix(response, "\r")
	if debugSwitch {
		fmt.Printf("\nDEBUG: user provided input (after TrimSuffix): '%s' \n", response)
	}
	// flush any output to the screen
	os.Stdout.Sync()
	// return the string from the user to the calling function
	return response
}

// checkContinue asks the user if they would like to continue with the
// currently running part of the application.
//
// checkContinue function reads input from the users console to see if
// they provide a a 'y' or 'n' response.
//
// The function returns a bool depending on the users response.
// if the response contains the letter 'y' it returns 'true'. Any other
// response will return 'false'.
func checkContinue() bool {
	// create a new reader from stdin
	reader := bufio.NewReader(os.Stdin)
	// ask the user a question
	fmt.Print("Continue? [y/n]: ")
	// read the users response - terminating their input on newline
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

// printBanner is used to print out program banner which displays:
// application name and application version
func printBanner() {
	fmt.Println("\n\t\t\tAcronym Management Tool 'amt'")
	fmt.Println("\t\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯")
}

// versionInfo function collects details of the program being run and
// displays it on stdout
func versionInfo() {

	// define a template for display on screen with place holders for data
	const appInfoTmpl = `
Running '{{.appname}}' version {{.appversion}}

 - Built with Go Complier '{{.compiler}}' on Golang version '{{.version}}'
 - Author's web site: http://www.wiremoons.com/
 - Source code for {{.appname}}: https://github.com/wiremoons/passgen/

All is well
`
	// build a map with keys set to match the template names used
	// and the data fields to be used in the template as values
	data := map[string]interface{}{
		"appname":    appname,
		"appversion": appversion,
		"compiler":   runtime.Compiler,
		"version":    runtime.Version(),
	}
	//
	t := template.Must(template.New("appinf").Parse(appInfoTmpl))
	if err := t.Execute(os.Stdout, data); err != nil {
		log.Fatalf("FATAL ERROR: in function 'versionInfo()' when building template with err: %v", err)
	}
}
