/* Acronym Management Tool (amt): main.h */

#ifndef MAIN_H_ /* Include guard */
#define MAIN_H_ 

/*
 *
 * Acronym Management Tool (amt)
 *
 * amt is program to manage acronyms held in an SQLite database
 *
 * author:     simon rowe <simon@wiremoons.com>
 * license:	open-source released under "MIT License"
 * 
 * Program to access a SQLite database and look up a requested acronym that maybe
 * held in a table called 'ACRONYMS'.
 * 
 * Also shall allow the creation of new acronym records, alterations of existing,
 * and deletion of records no longer required.
 * 
 * created: 20 Jan 2016 - version: 0.1 written - initial outline code written
 * 
 * 
 * The application uses the SQLite amalgamation source code files, so ensure they
 * are included in the same directory as this programs source code and then
 * compile with:
 * 
 * gcc -Wall main.c cli-args.c amt.c sqlite3.c -o amt.exe
 *
 */


#include <stdlib.h>     /* getenv */
#include <stdio.h>      /* printf */
#include <unistd.h>     /* strdup access */
#include <string.h>     /* strlen */
#include <malloc.h>     /* free for use with strdup */
#include <locale.h>     /* number output formatting with commas */

#include "cli-args.h"   /* manages the command line args from user */
#include "amt.h"        /* manages the database access for the application */
#include "sqlite3.h"    /* SQLite header */


/*
 *   APPLICATION GLOBAL VARIABLES
 */

/* path and name of acronyms database filename */
char *dbfile="";
/* handle to the database */
sqlite3 *db=NULL;
/* returned result codes from calling SQLite functions */
int rc=0;
/* data returned from SQL stmt run */
const char *data=NULL;
/* preprepared SQL query statement */
sqlite3_stmt *stmt=NULL;
/* set the version of the app here */
char appversion[]="0.1";
/* control debug outputs 0 == off | 1 == on */
int debug=0;
/* control help outputs request 0 == off | 1 == on */
int help=0;
/* string request on command line for acronym search */
char *findme;
/* request to add a new record 0 == off | 1 == on */
int newrec=0;

/*
 *   FUNCTION DECLARATIONS FOR main.c
 */

void exitCleanup(void);
void showHelp(void);
void printstart(void);

#endif // MAIN_H_ 
