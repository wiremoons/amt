// amt - program to access an SQLite database and lookup acronyms
//
// author:	Simon Rowe <simon@wiremoons.com>
// license: open-source released under The MIT License (MIT).
//
// Package used to manipulate the SQlite database for application 'amt'
//
// Example record of 'ACRONYMS' table in SQLite database for
// reference:
//
//   rowid 			: hidden internal sqlite record id
//   Acronym 		: 21CN
//   Definition 	: 21st Century Network
//   Description 	: A new BT network
//   Source 		: DFTS

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

// checkDB is used to verify if a valid database file name and path
// has been provided by the user.
//
// The database file name can be provided to the program via the
// command line or via an environment variable named: ACRODB. The
// function checks ensure the database file name exists, obtains its
// size on disk and checks it file permissions. These are items are
// output to Stdout by the function. If there are no errors the
// function returns.
//
// If the function fails for any reason the program is ended with the
// following exit codes:
//
//	-3 : no data base file provided or found
//	-4 : file is not a database file or it cannot be accessed
//
func checkDB() {
	// check if user has specified the location of the database to
	// use - either via command line or environment variable?
	if dbName == "" {
		// nothing provided via command line...
		if debugSwitch {
			log.Print("DEBUG: No database name provided via ")
			log.Println("command line input - check environment instead...")
			log.Println("DEBUG: Environment variable $ACRODB is:", os.Getenv("ACRODB"))
		}

		// get contents of environment variable $ACRODB as no filename
		// given on command line
		dbName = os.Getenv("ACRODB")

		// check if a database name and path was provided via the
		// environment variable
		if dbName == "" {

			if debugSwitch {
				log.Println("DEBUG: No database name provided via environment variable ACRODB")
				log.Println("DEBUG: Exit program")
			}

			// no database name provided via environment variable so
			// inform user and exit
			//
			// TODO : offer to create on instead here...?
			flag.Usage()
			log.Fatalln("FATAL ERROR: please provide the name of a database containing your acronyms\nrun 'amt --help' for more assistance\n")
		}
	}

	// dbName is not empty if we got here
	if debugSwitch {
		log.Printf("DEBUG: database provided is: %s", dbName)
		log.Printf("DEBUG: Checking file stats for: '%s'\n", dbName)
	}

	// check 'dbName' is valid file with os.Stats()
	fi, err := os.Stat(dbName)
	if err == nil {
		mode := fi.Mode()

		if debugSwitch {
			log.Printf("DEBUG: checking is '%s' is a regular file with os.Stat() call\n", dbName)
		}

		// check is a regular file
		if mode.IsRegular() {
			// print out some details of the database file:
			fmt.Printf("Database: %s   permissions: %s   size: %s bytes\n\n",
				fi.Name(), fi.Mode(), humanize.Comma(fi.Size()))

			if debugSwitch {
				log.Println("DEBUG: regular file check completed ok - return to main()")
			}
			// success - we are done!
			return
		}
	}

	if debugSwitch {
		log.Print("DEBUG: Exiting program as specified database file ")
		log.Printf("%s is not valid file that can be accessed", dbName)
	}
	// error found with the provided database file
	log.Fatalf("FATAL ERROR: database: '%s' is not a regular file\nError returned: %v\nrun 'amt --help' for more assistance\nABORT\n", dbName, err)

}

// checkCount provides the current total record count in the acronym
// table. The function takes no inputs. checkCount function returns
// the record count as an int64 variable. If an error occurs obtaining
// the record count from the database it will be printed to stderr.
func checkCount() int64 {

	if debugSwitch {
		log.Println("DEBUG: Getting record count function... ")
	}
	// create variable to hold returned database count of records
	var recCount int64
	// query the database to get number of records - result out in
	// variable recCount
	err := db.QueryRow("select count(*) from ACRONYMS;").Scan(&recCount)
	if err != nil {
		log.Printf("ERROR in function 'checkCount()' with SQL QueryRow: %v\n", err)
	}
	if debugSwitch {
		log.Printf("DEBUG: records count in table returned: %d\n", recCount)
	}
	// return the result
	return recCount
}

// lastAcronym obtains the last acronym entered into the acronym
// table. The lastAcronym function takes not inputs. The lastAcronym
// function returns the last acronym entered into the table as a
// string variable. If an error occurs obtaining the last acronym
// entered from the database it will be printed to stderr.
//
// SQL statement run is:
//
//    SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;
func lastAcronym() string {

	if debugSwitch {
		log.Println("DEBUG: Getting last entered acronym... ")
	}
	// create variable to hold returned database query
	var lastEntry string
	// query the database to get last entered acronym - result
	// returned to variable 'lastEntry'
	err := db.QueryRow("SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;").Scan(&lastEntry)
	if err != nil {
		log.Printf("ERROR: in function 'lastAcronym()' with SQL  QueryRow (lastEntry): %v\n", err)
	}

	if debugSwitch {
		log.Printf("DEBUG: last acronym entry in table returned: %s\n", lastEntry)
	}
	// return the result
	return lastEntry
}

// sqlVersion provides the version of SQLite library that is being
// used by the program. The function take no parameters. The
// sqlVersion function returns a string with a version number obtained
// by running the SQLite3 statement:
//
//     SELECT SQLITE_VERSION();
func sqlVersion() string {

	if debugSwitch {
		log.Println("DEBUG: Getting SQLite3 database version of software... ")
	}
	// create variable to hold returned database query
	var dbVer string
	// query the database to get version - result returned to variable
	// 'dbVer'
	err := db.QueryRow("select SQLITE_VERSION();").Scan(&dbVer)
	if err != nil {
		log.Printf("ERROR: in function 'sqlVersion()' with SQL QueryRow (dbVer): %v\n", err)
	}
	if debugSwitch {
		log.Printf("DEBUG: last acronym entry in table returned: %s\n", dbVer)
	}
	// return the result
	return dbVer
}

// getSources provide the current 'sources' held in the acronym table
// getSources function takes no parameters. The getSources functions
// returns a string contain a list of distinct 'source' records such
// as "General ICT"
func getSources() string {

	if debugSwitch {
		log.Print("DEBUG: Getting source list function... ")
	}
	// create variable to hold returned database source list
	var sourceList []string
	// query the database to extract distinct 'source' records -
	// result out in variable 'sourceList'
	rows, err := db.Query("select distinct(source) from acronyms;")
	if err != nil {
		log.Printf("ERROR: in function 'getSources()' with SQL QueryRow (sourceList): %v\n", err)
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
		log.Fatalf("\n\nFATAL ERROR: The source # you entered '%d' is greater than choices of '0' to '%d' offered, or less than zero\n\n", idxFinal, (len(sourceList) - 1))
	}
	// return the result
	return string(sourceList[idxFinal])
}

// addRecord function adds a new record to the acronym table held in
// the SQLite database It does not take any parameters. It does not
// return any information, and exits the program on completion. The
// applcation will exit of there is an error attempting to insert the
// new record into the database.
//
// The SQL insert statement used is:
//
//    insert into ACRONYMS(Acronym, Definition, Description, Source)
//    values(?,?,?,?)
func addRecord() {

	if debugSwitch {
		log.Printf("DEBUG: Adding new record function... \n")
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
		_, err := db.Exec("insert into ACRONYMS(Acronym, Definition, Description, Source) values(?,?,?,?)",
			acronym, definition, description, source)
		if err != nil {
			log.Fatalf("FATAL ERROR inserting new acronym record: %v\n", err)
		}
		// get new database record count post insert
		newInsertCount := checkCount()
		// inform user of difference in database record counts -
		// should be 1
		fmt.Printf("SUCCESS: %d record added to the database\n",
			(newInsertCount - preInsertCount))
		// inform user of database record counts
		fmt.Printf("\nDatabase record count is: %s  [was: %s]\n",
			humanize.Comma(newInsertCount), humanize.Comma(preInsertCount))
	}

	// function complete
	return
}

// searchRecord function obtains a string from the users and search
// for it in the SQLite acronyms database. It does not take any
// parameters. It does not return any information, and exits the
// program on completion. The applcation will exit of there is an
// error.
//
// The SQL insert statement used is:
//
//    select rowid,Acronym,Definition,Description,Source from ACRONYMS where
//    Acronym like ? ORDER BY Source;
func searchRecord() {
	// start search for an acronym - update user's screen
	fmt.Printf("\n\nSEARCH FOR ACRONYM\n¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n")
	//
	// check we have a term to search for in the acronyms database:
	if debugSwitch {
		fmt.Printf("DEBUG: checking for a search term ... ")
	}
	if debugSwitch {
		fmt.Printf("search term provided: %s\n", searchTerm)
	}
	// update user that the database is open and acronym we will
	// search for in how many records:
	fmt.Printf("\nSearching for:  '%s'  across %s records - please wait...\n",
		searchTerm, humanize.Comma(recCount))

	// flush any output to the screen
	os.Stdout.Sync()

	// run a SQL query to find any matching acronyms to that provided
	// by the user
	rows, err := db.Query("select rowid,Acronym,Definition,Description,Source from ACRONYMS where Acronym like ? ORDER BY Source;", searchTerm)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Printf("\nMatching results are:\n\n")
	for rows.Next() {
		// variables to hold returned database values - use []byte
		// instead of string to get around NULL values issue error
		// which states:
		//
		// " Scan error on column index 2: unsupported driver -> Scan
		// pair: <nil> -> *string"
		var rowid, acronym, definition, description, source []byte
		err := rows.Scan(&rowid, &acronym, &definition, &description, &source)
		if err != nil {
			fmt.Printf("ERROR: reading database record: %v", err)
		}
		// print the current row to screen - need string(...) as
		// values are bytes
		fmt.Printf("ID: %s\nACRONYM: '%s' is: %s.\nDESCRIPTION: %s\nSOURCE: %s\n\n",
			string(rowid), string(acronym), string(definition), string(description), string(source))
	}
	// check there were no other error while reading the database rows
	err = rows.Err()
	if err != nil {
		fmt.Printf("ERROR: reading database row returned: %v", err)
	}
	// function complete ok
	return
}
