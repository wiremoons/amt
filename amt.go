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
	"path/filepath"
	"strconv"

	"github.com/dustin/go-humanize"
	_ "github.com/mattn/go-sqlite3"
)

// SET GLOBAL VARIABLES

// set the version of the app here
var appversion = "0.5.2"
var appname string

// below are the flag variables used for command line args
var dbName string
var searchTerm string
var wildLookUp bool
var debugSwitch bool
var helpMe bool
var addNew bool
var showVer bool

// used to hold any errors
var err error

// create a global db handle - so can be used across functions
var db *sql.DB

// init always runs before applications main() function and is used here to
// set-up the required 'flag' variables from the command line parameters
// provided by the user when they run the app.
func init() {
	// flag types available are: IntVar; StringVar; BoolVar
	// flag parameters are: variable, cmd line flag, initial value, description
	// description is used by flag.Usage() on error or for help output
	flag.StringVar(&dbName, "f", "", "\tUSE: '-f <database_name>' name and path to the SQLite database to use")
	flag.StringVar(&searchTerm, "s", "", "\tUSE: '-s <acronym>' acronym that is to be searched for in the database")
	flag.BoolVar(&wildLookUp, "w", false, "\tUSE: '-w' to search for any similar matches to the acronym provided")
	flag.BoolVar(&debugSwitch, "d", false, "\tUSE: '-d' to include additional debug output when run")
	flag.BoolVar(&helpMe, "h", false, "\tUSE: '-h' to provide help on using this program")
	flag.BoolVar(&showVer, "v", false, "\tUSE: '-v' display the version information for the program")
	flag.BoolVar(&addNew, "n", false, "\tUSE: '-n' to add a new acronym record")
	// get the command line args passed to the program
	flag.Parse()
	// get the name of the application as called from the command line
	appname = filepath.Base(os.Args[0])
}

// main is the application start up function for amt
func main() {

	// confirm if debug mode is enabled
	if debugSwitch {
		log.Println("DEBUG: Debug mode enabled")
		log.Printf("DEBUG: Debug mode enabled")
	}

	// print out start up banner
	printBanner()

	// if debug is enabled - confirm the command line parameters received
	if debugSwitch {
		log.Println("DEBUG: Command Line Arguments provided are:")
		log.Println("\t\tDatabase name to use via command line:", dbName)
		log.Println("\t\tAcronym to search for:", searchTerm)
		log.Println("\t\tLook for similar matches:", strconv.FormatBool(wildLookUp))
		log.Println("\t\tDisplay additional debug output when run:", strconv.FormatBool(debugSwitch))
		log.Println("\t\tDisplay additional help information:", strconv.FormatBool(helpMe))
		log.Println("\t\tAdd a new acronym record:", strconv.FormatBool(addNew))
	}

	// check if command line help was request?
	if helpMe {
		flag.Usage()
		versionInfo()
		os.Exit(0)
	}

	// check if command line application version was request?
	if showVer {
		versionInfo()
		os.Exit(0)
	}

	// check if a valid database file has been provided - either via the
	// environment variable $ACRODB or via the command line from the user
	checkDB()

	// open the database - or abort if fails
	if debugSwitch {
		fmt.Printf("DEBUG: Opening database: '%s' ... ", dbName)
	}

	// get handle to database file
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {

		if debugSwitch {
			log.Printf("DEBUG: FAILED to open %s with error: %v - will exit application\n", dbName, err)
			log.Println("DEBUG: Exit program with call to 'log.Fatal()'")
		}

		log.Fatalf("FATAL ERROR: unable to get handle to SQLite database file: %s\nError is: %v\n", dbName, err)
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
