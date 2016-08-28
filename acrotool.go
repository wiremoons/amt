/*

acrotool - program to access an SQLite database and lookup acronyms

author:	 Simon Rowe <simon@wiremoons.com>
license: open-source released under The MIT License (MIT).

The program accesses a SQLite database and looks up the requested acronym held in a table called 'ACRONYMS'.

04 Sep 2014: version 0.1.0 initial outline code written.
29 Sep 2014: version 0.2.0 add database integration, command line
			  params - basic functionality working now.
02 Oct 2014: version 0.3.0 add ability to enter new records
11 Jul 2015: version 0.4.0 show source list on add new record
28 Aug 2016: version 0.5.0 changed to schematic versioning, reformated code and 			  tidy up, changed to MIT license.
28 Aug 2016: version 0.5.1 added ability to view and select existing acronym
			  source entries from enabled
28 Aug 2016: version 0.5.2 added the display of last acronym entered into the
			  database for user reference. Updated to mattn/go-sqlite3 latest
			  version so now running SQLite3 3.14.0 from 3.8.5. Added new app
			  startup message to include SQLite version.

*/

package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
)

// SET GLOBAL VARIABLES

// set the version of the app here
var appversion = "0.5.2"

// below are the flag variables used for command line args
var dbName string
var searchTerm string
var wildLookUp bool
var debugSwitch bool
var helpMe bool
var addNew bool

// create a global db handle - so can be used across functions
var db *sql.DB

// init always runs before applications main() function and is used here to
// set-up the required 'flag' variables from the command line parameters
// provided by the user when they run the app.
func init() {
	// flag types available are: IntVar; StringVar; BoolVar
	// flag parameters are: variable, cmd line flag, initial value, description
	// description is used by flag.Usage() on error or for help output
	flag.StringVar(&dbName, "i", "", "\tUSE: '-i <database_name>' name and path to the SQLite database to use")
	flag.StringVar(&searchTerm, "s", "", "\tUSE: '-s <acronym>' acronym that is to be searched for in the database [MANDATORY]")
	flag.BoolVar(&wildLookUp, "w", false, "\tUSE: '-w=true' to search for any similar matches to the acronym provided")
	flag.BoolVar(&debugSwitch, "d", false, "\tUSE: '-d=true' to include additional debug output when run")
	flag.BoolVar(&helpMe, "h", false, "\tUSE: '-h=true' to provide more detailed help on using this program")
	flag.BoolVar(&addNew, "n", false, "\tUSE: '-n=true' to add a new acronym record")
}

// main is the application start up function for acrotool
func main() {
	// print out start up banner
	printBanner()
	// get the command line args passed to the program
	flag.Parse()
	// confirm debug mode is enabled
	if debugSwitch {
		fmt.Println("DEBUG: Debug mode enabled")
	}
	// if debug is enabled - confirm the command line parameters received
	if debugSwitch {
		fmt.Println("DEBUG: Command Line Arguments provided are:")
		fmt.Println("\tDatabase name to use via command line:", dbName)
		fmt.Println("\tAcronym to search for:", searchTerm)
		fmt.Println("\tLook for similar matches:", strconv.FormatBool(wildLookUp))
		fmt.Println("\tDisplay additional debug output when run:", strconv.FormatBool(debugSwitch))
		fmt.Println("\tDisplay additional help information:", strconv.FormatBool(helpMe))
		fmt.Println("\tAdd a new acronym record:", strconv.FormatBool(addNew))
	}

	// check if a valid database file has been provided - either via the
	// environment variable $ACRODB or via the command line from the user
	checkDB()

	// open the database - or abort if fails
	if debugSwitch {
		fmt.Printf("DEBUG: Opening database: '%s' ... ", dbName)
	}
	// declare err as db is global var so already exists
	// otherwise get: "panic: runtime error: invalid memory address or nil pointer de-reference"
	var err error
	// get global handle to database
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		if debugSwitch {
			fmt.Printf("DEBUG: FAILED to open %s with error: %v - will exit application\n", dbName, err)
		}
		log.Fatal(err)
		if debugSwitch {
			fmt.Println("DEBUG: Exit program")
		}
		os.Exit(-6)
	}
	defer db.Close()

	// check connection to database is ok
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database connection status:  √")

	// get the SQLite database version we are comppiled with
	fmt.Printf("SQLite3 Database Version:  %s\n", sqlVersion())
	// get current record count for future use
	recCount := checkCount()
	fmt.Printf("Current record count is:  %s\n", humanize.Comma(recCount))
	// show last acronym entered in the database for info
	fmt.Printf("Last acronym entered was:  '%s'\n", lastAcronym())

	// see if the user want to add a new record via the -n command line switch
	if addNew {
		addRecord()
	}

	// ok - must want to search for an acronym
	fmt.Printf("\n\nSEARCH FOR ACRONYM\n¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n")
	//
	// check we have a term to search for in the acronyms database:
	if debugSwitch {
		fmt.Printf("DEBUG: checking for a search term ... ")
	}
	//
	// check if there a value in searchTerm provided via the command line...?
	if searchTerm == "" {
		if debugSwitch {
			fmt.Println("\nDEBUG: no search term found - asking the user for one")
		}

		// no search term found on command line - prompt the user for one:
		searchTerm = getInput("Enter an acronym to find: ")

		// check if searchTerm is populated now...
		if searchTerm == "" {
			fmt.Println("\nERROR: please ensure you enter the acronym you want to find\nrun 'acrotool --help' for more assistance\nABORT")
			if debugSwitch {
				fmt.Println("DEBUG: Exit program")
			}
			os.Exit(-7)
		}
	}
	if debugSwitch {
		fmt.Printf("search term provided: %s\n", searchTerm)
	}

	// update user that the database is open and acronym we will search for in how many records:
	fmt.Printf("\nDatabase status: OPEN - \tSearching for:  '%s'  across %s records - please wait ...\n", searchTerm, humanize.Comma(recCount))

	// flush any output to the screen
	os.Stdout.Sync()

	// Example record:
	//   rowid 			: hidden internal sqlite record id
	//   Acronym 		: 21CN
	//   Definition 	: 21st Century Network
	//   Description 	: A new BT network
	//   Source 		: DFTS

	// Example SQL queries
	// Last inserted acronym record:
	//		SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;
	// Search for acronym:
	//		"select Acronym,Definition,Description,Source from ACRONYMS where Acronym like ? ORDER BY Source;", searchTerm

	// run a SQL query to find any matching acronyms to that provided by the user
	rows, err := db.Query("select Acronym,Definition,Description,Source from ACRONYMS where Acronym like ? ORDER BY Source;", searchTerm)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("\nMatching results are:\n\n")
	for rows.Next() {
		// variables to hold returned database values - use []byte instead of string to get around NULL values issue
		// error:  " Scan error on column index 2: unsupported driver -> Scan pair: <nil> -> *string"
		var acronym, definition, description, source []byte
		err := rows.Scan(&acronym, &definition, &description, &source)
		if err != nil {
			fmt.Printf("ERROR: reading returned database result: %v", err)
		}
		// print the current row to screen - need string(...) as values are bytes
		fmt.Printf("ACRONYM: '%s' is: %s.\nDESCRIPTION: %s\nUSED BY: %s\n\n",
			string(acronym), string(definition), string(description), string(source))
	}
	// check there were no other error while reading the database rows
	err = rows.Err()
	if err != nil {
		fmt.Printf("ERROR: row reading error found: %v", err)
	}

	// END OF MAIN()
	fmt.Printf("\nAll is well\n")

}

// printBanner is used to print out program banner which displays:
// application name and application version
func printBanner() {
	fmt.Printf("\n\t\t\tAcronym Search - version: %s\n\n", appversion)
}

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

// checkDB is used to verify if a valid database file name and path has been
// provided by the user.
//
// The database file name can be provided to the program via the command line
// or via an environment variable named: ACRODB.
// The function checks ensure the database file name exists, obtains its size
// on disk and checks it file permissions. If there are no errors the function // returns. These are items are output to Stdout by the function.
//
// If the function fails for any reason the program is ended with the following
// exit codes:
//
//	-3 : no data base file provided or found
//	-4 : file is not a database file or it cannot be accessed
//
func checkDB() {
	// check if user has specified the location of the database to use - either via command line or environment variable?
	if dbName == "" {
		// nothing provided via command line...
		if debugSwitch {
			fmt.Println("DEBUG: No database name provided via command line input - check environment instead...")
			fmt.Println("DEBUG: Environment variable $ACRODB is:", os.Getenv("ACRODB"))
		}
		// get the content of environment variable $ACRODB - if exists...
		dbName = os.Getenv("ACRODB")
		// check if a database name and path was provided via the environment variable
		if dbName == "" {
			// no database name provided via environment variable either - tell user and exit
			if debugSwitch {
				fmt.Println("DEBUG: No database name provided via environment variable ACRODB")
			}
			fmt.Println("ERROR: please provide the name of a database containing your acronyms\nrun 'acrotool --help' for more assistance")
			flag.Usage()
			if debugSwitch {
				fmt.Println("DEBUG: Exit program")
			}
			os.Exit(-3)
		}
	}
	// ok - we have a dbName provided - make sure it is valid file - get os.Stats() for the provided filename: dbName
	if debugSwitch {
		fmt.Printf("DEBUG: database provided is: %s", dbName)
		fmt.Printf("DEBUG: Checking file stats for: '%s'\n", dbName)
	}
	fi, err := os.Stat(dbName)
	// if no error from os.Stat() call
	if err == nil {
		mode := fi.Mode()
		// check is a regular file?
		if debugSwitch {
			fmt.Printf("DEBUG: checking is '%s' is a regular file with os.Stat() call\n", dbName)
		}
		if mode.IsRegular() {
			// print out some details of the database file:
			fmt.Printf("Database: %s   permissions: %s   size: %s bytes\n\n", fi.Name(), fi.Mode(), humanize.Comma(fi.Size()))
			// we are done!
			if debugSwitch {
				fmt.Println("DEBUG: regular file check completed ok - return to main()")
			}
			return
		} else {
			fmt.Printf("ERROR: database: '%s' is not a regular file\nrun 'acrotool --help' for more assistance\nABORT\n", dbName)
			if debugSwitch {
				fmt.Println("DEBUG: Exit program")
			}
			os.Exit(-4)
		}
		// os.Stat() error occurred - so update user and exit
	} else {
		fmt.Printf("ERROR: unable to verify database: '%s' as error returned: %v\nrun 'acrotool --help' for more assistance\nABORT\n", dbName, err)
		if debugSwitch {
			fmt.Println("DEBUG: Exit program as os.Stat() failed")
		}
		os.Exit(-4)
	}
	// complete
}

// checkCount provides the current record count in the acronym table.
// The function takes not inputs. The function returns the record count as an
// int64 variable. If an error occurs obtaining the record count from the
// database it will be printed to Stdout.
func checkCount() int64 {
	if debugSwitch {
		fmt.Print("DEBUG: Getting record count function... ")
	}
	// create variable to hold returned database count of records
	var recCount int64
	// query the database to get number of records - result out in variable recCount
	err := db.QueryRow("select count(*) from ACRONYMS;").Scan(&recCount)
	if err != nil {
		fmt.Printf("QueryRow: %v\n", err)
	}
	if debugSwitch {
		fmt.Printf("DEBUG: records count in table returned: %d\n", recCount)
	}
	// return the result
	return recCount
}

// lastAcronym obtains the acronym entered into the acronym table.
// The function takes not inputs. The function returns the last acronym as a
// string variable. If an error occurs obtaining the acronym from the
// database it will be printed to Stdout.
//
// SQL statement run is:
// 		SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;
func lastAcronym() string {
	if debugSwitch {
		fmt.Print("DEBUG: Getting last entered acronym... ")
	}
	// create variable to hold returned database count of records
	var lastEntry string
	// query the database to get last entered acronym - result out in
	// variable 'lastEntry'
	err := db.QueryRow("SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;").Scan(&lastEntry)
	if err != nil {
		fmt.Printf("QueryRow (lastEntry): %v\n", err)
	}
	if debugSwitch {
		fmt.Printf("DEBUG: last acronym entry in table returned: %s\n", lastEntry)
	}
	// return the result
	return lastEntry
}

// sqlVersion provides the version of SQLite that is being used by the
// program. The function take no parameters. The function returns a string
// with a version number obtained by running the SQLite3 statement:
//		select SQLITE_VERSION();
func sqlVersion() string {
	if debugSwitch {
		fmt.Print("DEBUG: Getting SQLite3 database version of software... ")
	}
	// create variable to hold returned database count of records
	var dbVer string
	// query the database to get version - result out in
	// variable 'dbVer'
	err := db.QueryRow("select SQLITE_VERSION();").Scan(&dbVer)
	if err != nil {
		fmt.Printf("QueryRow (dbVer): %v\n", err)
	}
	if debugSwitch {
		fmt.Printf("DEBUG: last acronym entry in table returned: %s\n", dbVer)
	}
	// return the result
	return dbVer
}

// getSources provide the current 'sources' held in the acronym table
// It takes no parameters. It returns a string contain a list
// of distinct 'source' records such as "General ICT"
func getSources() string {

	if debugSwitch {
		fmt.Print("DEBUG: Getting source list function... ")
	}
	// create variable to hold returned database source list
	sourceList := make([]string, 0)
	// query the database to extract distinct 'source' records - result
	// out in variable 'sourceList'
	rows, err := db.Query("select distinct(source) from acronyms;")
	if err != nil {
		// TODO: below should be to stderr or log.Fatal(err) ??
		fmt.Printf("QueryRow: %v\n", err)
	}
	defer rows.Close()

	var srcname string

	for rows.Next() {
		err = rows.Scan(&srcname)
		sourceList = append(sourceList, srcname)
		//fmt.Printf("Source: %s\n", srcname)
	}

	fmt.Printf("\nExisting %d acronym 'source' choices:\n\n", len(sourceList))
	for idx, source := range sourceList {
		fmt.Printf("[%d]: '%s'  ", idx, source)
	}
	fmt.Printf("\n\n")
	// ask user to choose one...
	idxChoice := getInput("Enter a source [#] for the new acronym: ")
	idxFinal, err := strconv.Atoi(idxChoice)
	// error - could not convert to Int so just return the string as is...
	if err != nil {
		return string(idxChoice)
	}
	// check the number entered is not greater or less than is should be..
	if (idxFinal > (len(sourceList) - 1)) || (idxFinal < 0) {
		// error - entered value is out of range warn user and exit
		fmt.Printf("\n\nERROR: The source # you entered '%d' is greater than choices of '0' to '%d' offered, or less than zero\n\n", idxFinal, (len(sourceList) - 1))
		os.Exit(-8)
	}
	// return the result
	return string(sourceList[idxFinal])
}

/*
age := 27
    rows, err := db.Query("SELECT name FROM users WHERE age=?", age)
    if err != nil {
            log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
            var name string
            if err := rows.Scan(&name); err != nil {
                    log.Fatal(err)
            }
            fmt.Printf("%s is %d\n", name, age)
    }
    if err := rows.Err(); err != nil {
            log.Fatal(err)
    }
*/

// addRecord adds a new record to the acronym table held in the SQLite database
// It does not take any parameters. It does not return any information.
func addRecord() {
	if debugSwitch {
		fmt.Printf("DEBUG: Adding new record function... \n")
	}
	// update screen for user
	fmt.Printf("\n\nADDING NEW RECORD\n¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n")
	fmt.Printf("Note: To abort the input of a new record - press 'Ctrl + c'\n\n")
	// get new acronym from user
	acronym := getInput("Enter the new acronym: ")
	// todo: check the acronym does not already exist - check with user to continue...
	definition := getInput("Enter the expanded version of the new acronym: ")
	description := getInput("Enter any description for the new acronym: ")
	// show list of sources currently used and get one from the user
	source := getSources()
	// check the user is happy with what has been collected from them...
	fmt.Printf("\nContinue to add new acronym:\n\tACRONYM: %s\n\tEXPANDED: %s\n\tDESCRIPTION: %s\n\tSOURCE: %s\n", acronym, definition, description, source)

	// get current database record count
	preInsertCount := checkCount()

	// see if user wants to continue with the
	if checkContinue() {
		// ok - add record to the database table
		_, err := db.Exec("insert into ACRONYMS(Acronym, Definition, Description, Source) values(?,?,?,?)", acronym, definition, description, source)
		if err != nil {
			fmt.Printf("ERROR inserting new acronym record: %v\n", err)
			os.Exit(-8)
		}
		// get new database record count post insert
		newInsertCount := checkCount()
		// inform user of difference in database record counts - should be 1
		fmt.Printf("SUCCESS: %d record added to the database\n", (newInsertCount - preInsertCount))
		// inform user of database record counts
		fmt.Printf("\nDatabase record count is: %s  [was: %s]\n", humanize.Comma(newInsertCount), humanize.Comma(preInsertCount))
	}
	// leave the program as record entered ok
	os.Exit(0)
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
