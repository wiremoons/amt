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

*/

package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"go-humanize"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// SET GLOBAL VARIABLES

// set the version of the app here
var appversion = "0.5.0"

// below are the flag variables used for command line args
var dbName string
var searchTerm string
var wildLookUp bool
var debugSwitch bool
var helpMe bool
var addNew bool

// create a global db handle - so can be used across functions
var db *sql.DB

// init() function - always runs before main() - used here to set-up required flags variables
// from the command line parameters provided by the user when they run the app
func init() {
	// IntVar; StringVar; BoolVar all required: variable, cmd line flag, initial value, description used by flag.Usage() on error / help
	flag.StringVar(&dbName, "i", "", "\tUSE: '-i <database_name>' name and path to the SQLite database to use")
	flag.StringVar(&searchTerm, "s", "", "\tUSE: '-s <acronym>' acronym that is to be searched for in the database [MANDATORY]")
	flag.BoolVar(&wildLookUp, "w", false, "\tUSE: '-w=true' to search for any similar matches to the acronym provided")
	flag.BoolVar(&debugSwitch, "d", false, "\tUSE: '-d=true' to include additional debug output when run")
	flag.BoolVar(&helpMe, "h", false, "\tUSE: '-h=true' to provide more detailed help on using this program")
	flag.BoolVar(&addNew, "n", false, "\tUSE: '-n=true' to add a new acronym record")
}

//-------------------------------------------------------------------------
// FUNCTION:  MAIN
//-------------------------------------------------------------------------

func main() {
	// print out start up banner
	printBanner()
	//-------------------------------------------------------------------------
	// sort out the command line arguments
	//-------------------------------------------------------------------------
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

	// check if a valid database file has been provided - either via the environment variable $ACRODB
	// or on the command line when the program is run
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
	fmt.Println("Database connection ok")

	// get current record count for future use
	recCount := checkCount()

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
	fmt.Printf("\nDatabase status: OPEN - \tSearching for:  '%s'  in %s records ... please wait ...\n", searchTerm, humanize.Comma(recCount))

	// flush any output to the screen
	os.Stdout.Sync()

	// Example record:
	//   rowid 			: hidden internal sqlite record id
	//   Acronym 		: 21CN
	//   Definition 	: 21st Century Network
	//   Description 	: A new BT network
	//   Source 		: DFTS

	// Example SQL queries
	// Last inserted records:
	//		SELECT * FROM acronyms Order by rowid DESC LIMIT 1;
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

//**********************  APPLICATION FUNCTIONS BELOW *************************

//-------------------------------------------------------------------------
// FUNCTION:  printBanner - print out program banner to show version
//-------------------------------------------------------------------------

func printBanner() {
	fmt.Printf("\n\t\t\tAcronym Search - version: %s\n\n", appversion)
}

//----------------------------------------------------------------------------
// FUNCTION: getInput - ask user a question and return there input
//----------------------------------------------------------------------------
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

//-------------------------------------------------------------------------
// FUNCTION:  checkDB - check if a valid dbName has been provided and it exists
//-------------------------------------------------------------------------

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

//-------------------------------------------------------------------------
// FUNCTION:  checkCount - provide the current record count in the acronym table
//-------------------------------------------------------------------------

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

//-------------------------------------------------------------------------
// FUNCTION:  getSources - provide the current sources in the acronym table
//-------------------------------------------------------------------------

func getSources() string {
	if debugSwitch {
		fmt.Print("DEBUG: Getting source list function... ")
	}
	// create variable to hold returned database source list
	var sourceList []byte
	// query the database to distinct source records - result out in variable sourceList
	err := db.QueryRow("select distinct(source) from acronyms;").Scan(&sourceList)
	if err != nil {
		fmt.Printf("QueryRow: %v\n", err)
	}
	if debugSwitch {
		fmt.Printf("DEBUG: source list in table returned: %s\n", string(sourceList))
	}
	// return the result
	return string(sourceList)
}

//-------------------------------------------------------------------------
// FUNCTION:  addRecord - add a new record to the acronym table
//-------------------------------------------------------------------------

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
	// show list of sources currently used
	fmt.Printf("Source Options: %s\n", getSources())
	source := getInput("Enter any source for the new acronym: ")
	fmt.Printf("Continue to add new acronym:\n\tACRONYM: %s\n\tEXPANDED: %s\n\tDESCRIPTION: %s\n\tSOURCE: %s\n", acronym, definition, description, source)

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
		// inform user of difference in database record counts - shoulld be 1
		fmt.Printf("%d records added to the database\n", (newInsertCount - preInsertCount))
	}

	os.Exit(0)
}

//----------------------------------------------------------------------------
// FUNCTION:  checkContinue - get user input on continue or not
// checkContinue function reads input from the users console to see if
// requesting a 'y' or 'n' response.
// Returns a bool depending on the users response.
// if the response contains the letter 'y' assume true
//----------------------------------------------------------------------------
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
