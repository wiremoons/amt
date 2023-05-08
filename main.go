// amt - program to access an SQLite database and lookup acronyms
//
// author:	Simon Rowe <simon@wiremoons.com>
// license: open-source released under The MIT License (MIT).

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"amt-go/lib"
)

// SET GLOBAL VARIABLES

// set the version of the app here prep var to hold app name
var Appversion = "0.6.0"
var Appname string

// flag() variables used for command line args
var DbName string
var searchTerm string
var wildLookUp bool
var DebugSwitch bool
var helpMe bool
var addNew bool
var showVer bool
var rmid string

// used to keep track of database record count
var RecCount int64

// used to hold any errors
var err error

// init() always runs before the applications main() function and is
// used here to set up the flag() variables from the command line
// parameters - which are provided by the user when they run the app.
func init() {
	// flag types available are: IntVar; StringVar; BoolVar
	// flag parameters are: variable; cmd line flag; initial value; description.
	// 'description' is used by flag.Usage() on error or for help output
	flag.StringVar(&DbName, "f", "", "\tprovide SQLite database `filename` and path")
	flag.StringVar(&searchTerm, "s", "", "\t`acronym` to search for")
	flag.StringVar(&rmid, "r", "", "\t`acronym id` to remove")
	flag.BoolVar(&wildLookUp, "w", false, "\tsearch for any similar matches")
	flag.BoolVar(&DebugSwitch, "d", false, "\tshow debug output")
	flag.BoolVar(&helpMe, "h", false, "\tdisplay help for this program")
	flag.BoolVar(&showVer, "v", false, "\tdisplay program version")
	flag.BoolVar(&addNew, "n", false, "\tadd a new acronym record")
	// get the command line args passed to the program
	flag.Parse()
	// get the name of the application as called from the command line
	Appname = filepath.Base(os.Args[0])
}

// main is the application start up function for amt
func main() {

	// inject needed global variables into out sub-package 'utils'
	lib.DebugSwitch = DebugSwitch
	lib.DbName = DbName
	lib.Appversion = Appversion
	lib.Appname = Appname
	lib.RecCount = RecCount

	// confirm if debug mode is enabled and display other command line
	// flags and their current status
	if DebugSwitch {
		log.Println("DEBUG: Debug mode enabled")
		log.Printf("DEBUG: Number of command line arguments set by user is: %d", flag.NFlag())
		log.Printf("DEBUG: Command line argument settings are:")
		log.Println("\t\tDatabase name to use via command line:", DbName)
		log.Println("\t\tAcronym to search for:", searchTerm)
		log.Println("\t\tAcronym to remove:", rmid)
		log.Println("\t\tLook for similar matches:", strconv.FormatBool(wildLookUp))
		log.Println("\t\tDisplay additional debug output when run:", strconv.FormatBool(DebugSwitch))
		log.Println("\t\tDisplay additional help information:", strconv.FormatBool(helpMe))
		log.Println("\t\tAdd a new acronym record:", strconv.FormatBool(addNew))
		log.Println("\t\tShow the applications version:", strconv.FormatBool(showVer))
	}

	// a function that will run at the end of the program
	defer func() {
		// END OF MAIN()
		fmt.Printf("\nAll is well\n")
	}()

	// override Go standard flag.Usage() function to get better
	// formatting and output by using my own function instead
	flag.Usage = func() {
		if DebugSwitch {
			log.Println("DEBUG: Running flag.Usage override function")
		}
		lib.MyUsage()
	}

	// print out start up banner
	if DebugSwitch {
		log.Println("DEBUG: Calling 'printBanner()'")
	}
	lib.PrintBanner()

	// check if a valid database file is available on the system
	if DebugSwitch {
		log.Println("DEBUG: Calling 'checkDB()'")
	}

	err = lib.CheckDB()
	if err != nil {
		log.Println(err)
		// no database found - offer to create one
		fmt.Printf("\nCreate a new database and add a few example acronyms?")
		if !lib.CheckContinue() {
			// no database available - exit application
			log.Fatal("ERROR: unable to continue without a valid acronym database.\n")
		}
		// user wants a new database - so attempt to create it in the same directory as the
		// program executable using the file named: 'amt-db.db' - set location here then attempt to open it
		DbName = filepath.Join(filepath.Dir(os.Args[0]), "amt-db.db")
	}
	// Setup and open the database ready for use
	if DebugSwitch {
		log.Println("DEBUG: database found - attempting to open with 'OpenDataBase()'")
	}
	err = lib.OpenDataBase()
	if err != nil {
		log.Println(err)
	}

	// attempt to populate the database with some example records if it
	// is empty - ask user first
	if (lib.CheckCount()) == 0 {
		fmt.Println("\nWould you like to add some initial records to your empty acronyms database?")
		if lib.CheckContinue() {
			err = lib.PopNewDB()
			if err != nil {
				// records could not be added - exit application
				log.Fatalf("ERROR: aborting program with error: %v\n", err)
			}
		}
	}

	if DebugSwitch {
		log.Println("DEBUG: Start 'switch'...")
	}

	switch {
	case helpMe:
		if DebugSwitch {
			log.Println("DEBUG: 'helpme' switch statement called")
		}
		flag.Usage()
		fallthrough

	case showVer:
		if DebugSwitch {
			log.Println("DEBUG: 'showVer' switch statement called")
		}
		lib.VersionInfo()

	case addNew:
		if DebugSwitch {
			log.Println("DEBUG: 'addNew' switch statement called")
		}
		lib.AddRecord()

	case len(searchTerm) > 0:
		if DebugSwitch {
			log.Println("DEBUG: search switch statement called")
		}
		lib.SearchRecord(searchTerm)

	case len(rmid) > 0:
		if DebugSwitch {
			log.Println("DEBUG: remove switch statement called")
		}
		_ = lib.RemoveRecord(rmid)

	default:
		if DebugSwitch {
			log.Println("DEBUG: Default switch statement called")
		}
		lib.VersionInfo()
		flag.Usage()
	}

	// PROGRAM END
}
