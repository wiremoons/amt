/*
 * Acronym Management Tool (amt): main.h
 *
 * amt is program to manage acronyms held in an SQLite database
 *
 * author:     simon rowe <simon@wiremoons.com>
 * license:    open-source released under "MIT License"
 * source:     https://github.com/wiremoons/amt
 * 
 * Program to access a SQLite database and look up a requested acronym that maybe
 * held in a table called 'ACRONYMS'.
 * 
 * Also supports the creation of new acronym records, alterations of existing,
 * and deletion of records no longer required.
 * 
 * created: 20 Jan 2016 - version: 0.1 written - initial outline code written
 * 
 * The application uses the SQLite amalgamation source code files, so ensure they
 * are included in the same directory as this programs source code and then
 * compile with:
 * 
 * gcc -Wall main.c cli-args.c amt.c sqlite3.c -o amt.exe
 *
 */

#ifndef MAIN_H_ /* Include guard */
#define MAIN_H_ 

#include "cli-args.h"	    /* manages the command line args from user */
#include "amt-db-funcs.h"   /* manages the database access for the application */
#include "sqlite3.h"	    /* SQLite header */

#include <stdlib.h>     /* to allow NULL to be used for globals var declarations */

/*
 *   APPLICATION GLOBAL VARIABLES
 */

char *dbfile="";	    /* path and name of acronyms database filename */
sqlite3 *db=NULL;	    /* handle to the database */
int rc=0;		    /* returned result codes from calling SQLite functions */
const char *data=NULL;	    /* data returned from SQL stmt run */
sqlite3_stmt *stmt=NULL;    /* preprepared SQL query statement */
char appversion[]="0.2.1";  /* set the version of the app here */
int help=0;		    /* control help outputs request 0 == off | 1 == on */
char *findme=NULL;	    /* string request on command line for acronym search */
int recordid=-1;	    /* database record id (rowid) used to delete records */
int newrec=0;		    /* request to add a new record 0 == off | 1 == on */

/* FUNCTION DECLARATIONS FOR main.c */

void exit_cleanup(void);
void show_help(void);
void print_start_screen(char *prog_name);
	
#endif // MAIN_H_ 
