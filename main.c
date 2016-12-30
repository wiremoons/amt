/* Acronym Management Tool (amt): main.c */

#include "main.h"

/* MAIN ENTRY POINT FOR APPLICATION */

int main(int argc, char **argv)
{

    /* register our atexit() function */
    atexit(exitCleanup);

    /* use locale to format numbers output */
    setlocale(LC_NUMERIC, "");
	
    /* get any command line arguments provided by the user and then
       process using getopts() via function below */
    getCLIArgs(argc,argv);

    /* Print application startup banner to the screen */
    printstart();

    /* Check if help output was requested? */
    if (help)
    {
        if (debug) { printf("\nDEBUG: User request help output\n"); }
        showHelp();
        return EXIT_SUCCESS;
    }

    /* check we have a database file to open */
    checkDB();

    /* Initialise and then open the database */
    sqlite3_initialize();
    /* open the database in read & write mode */
    rc = sqlite3_open_v2(dbfile, &db, SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE, NULL);
    /* check it opened OK - if not exit */
    if (rc != SQLITE_OK)
    {
        exit(EXIT_FAILURE);
    }
    /* Get number of records in database and display on screen */
    int totalrec = recCount();
    printf(" - Record count is: %'d\n",totalrec);


    /* perform db ops here */

	/* program exit ok */
    return (EXIT_SUCCESS);
}

/*
** FUNCTION: exitCleanup
**
** function called when program exits Used via registration with
** 'atexit()' in main() run any final checks and db close down here
**
*/
void exitCleanup(void)
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
    // free up an memory we allocated:
    if (findme != NULL) free(findme);
    exit(EXIT_SUCCESS);
}


/* FUNCTION: printstart */

/* Function to display basic information when application is started */


void printstart(void)
{
    printf("\n");
    printf("\t\tAcronym Management Tool\n"
		   "\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n");
    printf("Summary:\n"
		   " - App version: %s complied with SQLite version: %s\n",
		   appversion, SQLITE_VERSION);
}

/**-------- FUNCTION: showHelp
 ** 
 ** Show on screen a summary of the command line switches available in the
 ** program.
 ** 
 */
void showHelp(void)
{
    printf("\n"
		   "Help Summary:\n"
		   "The following command line switches can be used:\n\n"
		   "  -d\tDebug - include additional debug outputs when run\n"
		   "  -h\tHelp - Show this help information\n"
		   "  -v\tVersion - display the version of the program\n"
		   "  -n\tNew - add a new acronym record to the database\n"
		   "  -s ?\tSearch - find an acronym where ? == acronym to search for\n");
}

