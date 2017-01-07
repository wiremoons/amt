/* Acronym Management Tool (amt): amt.h */

#ifndef AMT_DB_FUNCS_H_ /* Include guard */
#define AMT_DB_FUNCS_H_ 


#include "sqlite3.h"    /* SQLite header */

extern char *dbfile;
extern sqlite3 *db;
extern int rc;
extern sqlite3_stmt *stmt;
extern const char *data;
extern char *findme;

/*
*   FUNCTION DECLARATIONS
*/

int recCount(void);			/* get current acronym record count */
void check4DB(void);			/* ensure database is accessible */
char *get_last_acronym();		/* get last acronym added to database */
int do_acronym_search(char *findme);	/* search database for 'findme' string */
	
#endif // AMT_DB_FUNCS_H_ 
