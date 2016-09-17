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
var appversion = "0.5.4"
var appname string

// below are the flag variables used for command line args
var dbName string
var searchTerm string
var wildLookUp bool
var debugSwitch bool
var helpMe bool
var addNew bool
var showVer bool

// used to keep track of database record count
var recCount int64

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
	flag.StringVar(&dbName, "f", "", "\tprovide SQLite database `filename` and path")
	flag.StringVar(&searchTerm, "s", "", "\t`acronym` to search for")
	flag.BoolVar(&wildLookUp, "w", false, "\tsearch for any similar matches")
	flag.BoolVar(&debugSwitch, "d", false, "\tshow debug output")
	flag.BoolVar(&helpMe, "h", false, "\tdisplay help for this program")
	flag.BoolVar(&showVer, "v", false, "\tdisplay program version")
	flag.BoolVar(&addNew, "n", false, "\tadd a new acronym record")
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
		log.Println("\t\tShow the applications version:", strconv.FormatBool(addNew))
	}

	// override Go standard flag.Usage function to get better
	// formating and output by using my own function instead
	flag.Usage = func() {
		myUsage()
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

	// see if the user wants to search for a acronym record in the database
	if len(searchTerm) > 0 {
		searchRecord()
	}

	// No specific application options given - show command line usage
	// to help users in case they are stuck
	fmt.Printf("\n")
	flag.Usage()

	// END OF MAIN()
	fmt.Printf("\nAll is well\n")

}
