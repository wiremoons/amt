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

	// confirm that debug mode is enabled and display other command
	// line flags and their current status for confirmation
	if debugSwitch {
		log.Println("DEBUG: Debug mode enabled")
		log.Printf("DEBUG: Number of command line arguments set by user is: %d", flag.NFlag())
		log.Printf("DEBUG: Command line argument settings are:")
		log.Println("\t\tDatabase name to use via command line:", dbName)
		log.Println("\t\tAcronym to search for:", searchTerm)
		log.Println("\t\tLook for similar matches:", strconv.FormatBool(wildLookUp))
		log.Println("\t\tDisplay additional debug output when run:", strconv.FormatBool(debugSwitch))
		log.Println("\t\tDisplay additional help information:", strconv.FormatBool(helpMe))
		log.Println("\t\tAdd a new acronym record:", strconv.FormatBool(addNew))
		log.Println("\t\tShow the applications version:", strconv.FormatBool(addNew))
	}

	// a function that will run at the end of the program
	defer func() {
		// END OF MAIN()
		fmt.Printf("\nAll is well\n")
	}()

	// override Go standard flag.Usage function to get better
	// formating and output by using my own function instead
	flag.Usage = func() {
		if debugSwitch {
			log.Println("DEBUG: Running flag.Usage override function")
		}
		myUsage()
	}

	// print out start up banner
	if debugSwitch {
		log.Println("DEBUG: Calling 'printBanner()'")
	}
	printBanner()

	// check if a valid database file is available on the system
	if debugSwitch {
		log.Println("DEBUG: Calling 'checkDB()'")
	}
	err = checkDB()
	if err != nil {
		log.Fatal(err)
	}

	// open the database and retrive initial and print to screen
	if debugSwitch {
		log.Println("DEBUG: Calling 'openDB()'")
	}
	err = openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if debugSwitch {
		log.Println("DEBUG: Start 'switch'...")
	}

	switch {
	// display help information for the program
	case helpMe:
		flag.Usage()
		versionInfo()
		break
	// display version information for the program
	case showVer:
		versionInfo()
		break
	// see if the user want to add a new record via the -n command line switch
	case addNew:
		addRecord()
		break
	// see if the user wants to search for a acronym record in the database
	case len(searchTerm) > 0:
		searchRecord()
		break
	// No specific application cli options given - show command line
	// usage to help user in case they are stuck
	default:
		if debugSwitch {
			log.Println("DEBUG: Default switch statement called")
		}
		versionInfo()
		flag.Usage()
	}

	// PROGRAM END
}
