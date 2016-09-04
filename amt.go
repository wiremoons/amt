// amt - program to access an SQLite database and lookup acronyms
//
// author:	Simon Rowe <simon@wiremoons.com>
// license: open-source released under The MIT License (MIT).

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

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

// main is the application start up function for amt
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
			fmt.Println("\nERROR: please ensure you enter the acronym you want to find\nrun 'amt --help' for more assistance\nABORT")
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
	fmt.Printf("\nSearching for:  '%s'  across %s records - please wait...\n", searchTerm, humanize.Comma(recCount))

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
		// variables to hold returned database values - use []byte instead
		// of string to get around NULL values issue error:
		// 		" Scan error on column index 2: unsupported driver ->
		//			Scan pair: <nil> -> *string"
		var acronym, definition, description, source []byte
		err := rows.Scan(&acronym, &definition, &description, &source)
		if err != nil {
			fmt.Printf("ERROR: reading database record: %v", err)
		}
		// print the current row to screen - need string(...) as values
		// are bytes
		fmt.Printf("ACRONYM: '%s' is: %s.\nDESCRIPTION: %s\nSOURCE: %s\n\n",
			string(acronym), string(definition), string(description), string(source))
	}
	// check there were no other error while reading the database rows
	err = rows.Err()
	if err != nil {
		fmt.Printf("ERROR: reading database row returned: %v", err)
	}

	// END OF MAIN()
	fmt.Printf("\nAll is well\n")

}
