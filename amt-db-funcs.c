/* Acronym Management Tool (amt): amt.c */

#include "amt-db-funcs.h"

#include <stdlib.h>     /* getenv */
#include <stdio.h>      /* printf */
#include <unistd.h>     /* strdup access and FILE */

/*
 * Run SQL query to obtain current number of acronyms in the database.
 */
int recCount(void)
{
	int totalrec = 0;
	rc = sqlite3_prepare_v2(db,"select count(*) from ACRONYMS",-1, &stmt, NULL);
	if ( rc != SQLITE_OK) {
		exit(-1);
	}

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		totalrec = sqlite3_column_int(stmt,0);
	}

	sqlite3_finalize(stmt);
	return(totalrec);
}

/*
 * Check for a valid database file to open
 */
void check4DB(void)
{

    dbfile = getenv("ACRODB");
    if (dbfile) {
	    printf(" - Database location: %s\n", dbfile);

	    if (access(dbfile, F_OK | R_OK) == -1) {
		    fprintf(stderr,"\n\nERROR: The database file '%s'"
			    " is missing or is not accessible\n\n"
			    , dbfile);
		    exit(EXIT_FAILURE);
		}
    } else {
	    printf("\tWARNING: No database specified using 'ACRODB' "
		   "environment variable\n");
	    exit(EXIT_FAILURE);
    }

/* TODO if neither of the above - check current directory we are
running in - or then offer to create a new db? otherwise exit prog
here */

}

/*
 * ADDING NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * Note: To abort the input of a new record - press 'Ctrl + c'
 * 
 * Enter the new acronym: KSLOC
 * Enter the expanded version of the new acronym: Thousands of Source Line Of Code
 *     Enter any description for the new acronym: The count in thousand of line of
 *     source code that makes up an application, lines of code excluding blank lines
 *     and coments.
 * Enter any source for the new acronym: General ICT
 * Continue to add new acronym:
 * 	ACRONYM: KSLOC
 * 	EXPANDED: Thousands of Source Line Of Code
 * 	DESCRIPTION: The count in thousand of line of source code that makes up an
 * 	application, lines of code excluding blank lines and comments.
 * SOURCE: General ICT
 * Continue? [y/n]: y
 * 1 records added to the database
 */ 
