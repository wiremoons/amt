/*

	Acronym Management Tool (amt)
	
	amt is program to manage acronyms held in an SQLite database
	
	author:	simon rowe <simon@wiremoons.com>
	license:	open-source released under "New BSD License"
	
	Program to access a SQLite database and look up a requested acronym
	that maybe held in a table called 'ACRONYMS'.
	
	Also shall allow the creation of new acronym records, alterations of
	exisiting, and deletion of records no longer required.
	
	created: 20 Jan 2016 - version: 0.1 written - initial outline code written
	
	
	Database used is SQLite. The location of the database can be stored
	in an environment variable called ACRODB. To set this you can
	specific the path as so:
	
	bash:
	export ACRODB=/home/simon/work/acrotool/Sybil.db
	
	powershell:
	$env:ACRODB += "C:\Users\Simon\Work\Scratch\Sybil\Sybil.db"
	
	cmd.exe / DOS:
	set ACRODB=C:\Users\Simon\Work\Scratch\Sybil\Sybil.db
	
	The application uses the SQLite amalgamation source code files, so
	ensure they are included in the same directory as this programs
	source code and then compile with:
	
	gcc -Wall amt.c sqlite3.c -o amt.exe
	
*/

#include "sqlite3.h"    // SQLite
#include <stdlib.h>     // getenv
#include <stdio.h>      //
#include <unistd.h>     // strdup access
#include <string.h>     // strlen

/*
   GLOBAL VARIABLES
   */
// path and acronyms database filename
char *dbfile="";
// handle to the database
sqlite3 *db = NULL;
// returned result codes from calling SQLite functions
int rc=0;
// set the version of the app here
char appversion[] = "0.1";
// control debug outputs 0 == off | 1 == on
int debug = 0;
// control help outputs request 0 == off | 1 == on
int help = 0;
// string request on command line for acronym search
char *findme;
// request to add a new record 0 == off | 1 == on
int newrec = 0;

/**-------- FUNCTION: printstart

Function: print application start banner

*/
void printstart()
{
    printf("\n");
    printf("\t\tAcronym Management Tool\n"
            "\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n");
    printf("Summary:\n"
            " - App version: %s complied with SQLite version: %s\n",
            appversion, SQLITE_VERSION);
}

/**-------- FUNCTION: getCLIArgs

  Function called when program starts. Used to parse command line
  options provided by the user. Uses the POSIX compliant getopts()

*/
void getCLIArgs(int argc, char **argv)
{
    opterr = 0;
    int c = 0;
    int index = 0;

    while ((c = getopt (argc, argv, "dhvns:")) != -1)
    {
        switch (c)
        {
            // debugging output was requested
            case 'd':
                debug = 1;
                break;
                // provide help summary to user was requested
            case 'h':
                help = 1;
                break;
                // the version of the application requested
            case 'v':
                printf("%s is version: %s\n",argv[0],appversion);
                exit(EXIT_SUCCESS);
                // request to add a new acronym record
            case 'n':
                newrec = 1;
                break;
                // request to search for an acronym
            case 's':
                findme = strdup(optarg);
                if ( strlen(findme) <= 0 || findme == NULL )
                {
                    fprintf(stderr,"ERROR: for -s option please provide "
                            "an acronym to search for");
                    exit(EXIT_FAILURE);
                }
                break;
                // ERROR HANDLING BELOW
                //
                // command line option given - but is missing the required data argument for it
            case ':':
                fprintf(stderr,"ERROR: '%s' option '-%c' requires an argument\n",argv[0], optopt);
                break;
                // invalid option provided on command line - also 'default' as the switch fall-thru
            case '?':
            default:
                /* invalid option */
                fprintf(stderr,"ERROR: '%s' option '-%c' is invalid or missing input data\n",
                        argv[0], optopt);
                break;
        }
    }
    // if debugging requested - display extra getopt() info
    if (debug)
    {
        printf ("DEBUG: optargs() values:\n"
                "\tdebug = %s\n"
                "\tHelp requested = %s\n"
                "\tNew record input requested = %s\n"
                "\tSearch requested for string = %s\n\n",
                debug ? "true" : "false",
                help ? "true" : "false",
                newrec ? "true" : "false",
                findme);
        for (index = optind; index < argc; index++)
        {
            printf ("\tInvalid option argument(s) seen: %s\n", argv[index]);
        }
    }
}


/**-------- FUNCTION: exitCleanup

  function called when program exits
  Used via registration with 'atexit()' in main()
  run any final checks and db close down here

*/
void exitCleanup()
{
    // check if a database handle was created and assigned yet
    if (db == NULL)
    {
        printf("\nNo SQLite database shutdown required\n\nAll is well\n");
        exit(EXIT_SUCCESS);
    }
    //
    // db handle exists - so close the database connection
    rc = sqlite3_close_v2(db);
    // if did not close properly
    if (rc != SQLITE_OK)
    {
        fprintf(stderr,"\nWARNING: error '%s' when trying to close the database\n",
                sqlite3_errstr(rc));
        exit(EXIT_FAILURE);
    }
    // close down and exit
    sqlite3_shutdown();
    printf("\nCompleted SQLite database shutdown\n\nAll is well\n");
    exit(EXIT_SUCCESS);
}


/**-------- FUNCTION: showHelp

  Show on screen a summary of the command line switches available in the
  program.

*/
void showHelp()
{
    printf("\n"
            "Help Summary:\n"
            "The following command line switches can be used:\n\n"
            "  -d\tDebug - include addtional debug outputs when run\n"
            "  -h\tHelp - Show this help information\n"
            "  -v\tVersion - display the version of the program\n"
            "  -n\tNew - add a new acronym record to the database\n"
            "  -s ?\tSearch - find an acronym where ? == acronym to search for\n");
}

/**-------- FUNCTION: checkDB

Function: check for a valid database file to open

*/
void checkDB()
{

    /* check if acronyms database file was supplied on the command line */

    /* if ( ! dbfile = "") */
    /* { */

    /* } */

    /* obtain the acronyms database file from the environment */
    dbfile = getenv("ACRODB");
    if (dbfile)
    {
        printf(" - Database location: %s\n", dbfile);
        /* check database file is valid and accessable */
        if (access(dbfile, F_OK | R_OK) == -1)
        {
            fprintf(stderr,"\n\nERROR: The database file '%s'"
                    " is missing or is not accessable\n\n", dbfile);
            exit(EXIT_FAILURE);
        }
    } else {
        printf("\tWARNING: No database specified using 'ACRODB' environment variable\n");
        exit(EXIT_FAILURE);
    }
    // if neither of the above - check current directory we are running
    // in - or then ofer to create a new db? otherwise exit prog here

}

/* MAIN ENTRY POINT FOR APPLICATION */

int main(int argc, char **argv)
{
    // register our atexit() function
    atexit(exitCleanup);

    // get any command line arguments provided by the user
    // and then process using getopts() via function below
    getCLIArgs(argc,argv);

    /* Print application startup banner to the screen */
    printstart();

    // Check it it was just help output the user requested?
    if (help)
    {
        if (debug) { printf("\nDEBUG: User request help output\n"); }
        showHelp();
        return EXIT_SUCCESS;
    }

    // check we have a database file to open
    checkDB();

    /* Initialise and then open the database */
    sqlite3_initialize();
    // open the database in read & write mode
    rc = sqlite3_open_v2(dbfile, &db, SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE, NULL);
    /* check it opened OK - if not exit */
    if (rc != SQLITE_OK)
    {
        exit(EXIT_FAILURE);
    }

    /* perform db ops here */
    return (EXIT_SUCCESS);
}

