/* Acronym Management Tool (amt): amt.h */

#ifndef AMT_H_ /* Include guard */
#define AMT_H_ 

#include <stdlib.h>     /* getenv */
#include <stdio.h>      /* */
#include <unistd.h>     /* strdup access and FILE */


#include "sqlite3.h"    /* SQLite header */

extern char *dbfile;
extern sqlite3 *db;
extern int rc;
extern sqlite3_stmt *stmt;
extern int debug;
extern const char *data;

/*
*   FUNCTION DECLARATIONS
*/

int recCount(void);						/* get current acronym record count */
void checkDB(void);						/* ensure database is accessible */

#endif // AMT_H_ 
