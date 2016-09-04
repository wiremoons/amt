package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

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
			fmt.Println("ERROR: please provide the name of a database containing your acronyms\nrun 'amt --help' for more assistance")
			flag.Usage()
			if debugSwitch {
				fmt.Println("DEBUG: Exit program")
			}
			os.Exit(-3)
		}
	}

	if debugSwitch {
		fmt.Printf("DEBUG: database provided is: %s", dbName)
		fmt.Printf("DEBUG: Checking file stats for: '%s'\n", dbName)
	}

	// check 'dbName' is valid file with os.Stats()
	fi, err := os.Stat(dbName)

	if err == nil {
		mode := fi.Mode()

		if debugSwitch {
			fmt.Printf("DEBUG: checking is '%s' is a regular file with os.Stat() call\n", dbName)
		}

		// check is a regular file
		if mode.IsRegular() {
			// print out some details of the database file:
			fmt.Printf("Database: %s   permissions: %s   size: %s bytes\n\n", fi.Name(), fi.Mode(), humanize.Comma(fi.Size()))

			if debugSwitch {
				fmt.Println("DEBUG: regular file check completed ok - return to main()")
			}

			// we are done!
			return
		}

		// error found with the provided database file
		fmt.Printf("ERROR: database: '%s' is not a regular file\nrun 'amt --help' for more assistance\nABORT\n", dbName)
		if debugSwitch {
			fmt.Println("DEBUG: Exit program")
		}
		os.Exit(-4)

		// os.Stat() error occurred - so update user and exit
	} else {
		fmt.Printf("ERROR: unable to verify database: '%s' as error returned: %v\nrun 'amt --help' for more assistance\nABORT\n", dbName, err)
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
	var sourceList []string
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
