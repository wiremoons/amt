/* Acronym Management Tool (amt): amt.c */

#include "amt-db-funcs.h"

#include <stdlib.h>		/* getenv */
#include <stdio.h>		/* printf */
#include <unistd.h>		/* strdup access stat and FILE */
#include <string.h>		/* strlen strdup */
#include <malloc.h>		/* free for use with strdup */
#include <locale.h>		/* number output formatting with commas */
#include <sys/types.h>		/* stat */
#include <sys/stat.h>		/* stat */
#include <time.h>		/* stat file modification time */

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

		struct stat sb;
		int check;
    
		check = stat (dbfile, &sb);
    
		if (check) {
			perror("\nERROR: call to 'stat' for database file failed\n");
			exit(EXIT_FAILURE);
		}
    
		printf(" - Database size: %'ld bytes\n",sb.st_size);
		printf(" - Database last modified: %s\n",ctime(&sb.st_mtime));
	    
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
 * Obtain the last acronym entered into the database
 */
char *get_last_acronym()
{
	char *acronym_name;
	
	rc = sqlite3_prepare_v2(db,"SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;",-1, &stmt, NULL);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL error: %s\n", sqlite3_errmsg(db));
		exit(-1);
	}

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		acronym_name = strdup((const char*)sqlite3_column_text(stmt,0));
	}

	sqlite3_finalize(stmt);

	if (acronym_name == NULL) {
		fprintf(stderr,"ERROR: last acronym lookup return NULL\n");
	}

	return(acronym_name);
}


/*
 * SEARCH FOR A NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select rowid,Acronym,Definition,Description,Source from ACRONYMS where Acronym like ? ORDER BY Source;
 */

int do_acronym_search(char *findme)
{
	printf("\nSearchning for: '%s' in database...\n\n",findme);

	rc = sqlite3_prepare_v2(db,"select rowid,Acronym,Definition,Description,Source from ACRONYMS where Acronym like ? ORDER BY Source;",-1, &stmt, NULL);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL error: %s\n", sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	sqlite3_bind_text(stmt,1,(const char*)findme,-1,SQLITE_STATIC);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL error: %s\n", sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	int search_rec_count = 0;
	while(sqlite3_step(stmt) == SQLITE_ROW) {
		printf("ID:          %s\n", (const char*)sqlite3_column_text(stmt,0));
		printf("ACRONYM:     '%s' is: %s.\n", (const char*)sqlite3_column_text(stmt,1),(const char*)sqlite3_column_text(stmt,2));
		/* printf("DEFINITION:  %s\n", (const char*)sqlite3_column_text(stmt,2)); */
		printf("DESCRIPTION: %s\n", (const char*)sqlite3_column_text(stmt,3));
		printf("SOURCE:      %s\n\n", (const char*)sqlite3_column_text(stmt,4));
		search_rec_count++;
	}

	sqlite3_finalize(stmt);
	
	return search_rec_count;
}

/*
 * DELETE A RECORD BASE ROWID
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select rowid,Acronym,Definition,Description,Source from ACRONYMS where rowid = ?;
 *
 * delete from ACRONYMS where rowid = ?;
 */

/*
 * CHECKING SQLITE VERSION
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select SQLITE_VERSION();
 */


/*
 * GETTING LIST OF ACRONYM SOURCES
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select distinct(source) from acronyms;
 */


/*
 * ADDING NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * insert into ACRONYMS(Acronym, Definition, Description, Source) values(?,?,?,?)
 * 
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
