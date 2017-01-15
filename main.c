/* Acronym Management Tool (amt): main.c */

#include "main.h"

#include <stdlib.h> /* getenv */
#include <stdio.h>  /* printf */
#include <unistd.h> /* strdup access */
#include <string.h> /* strlen strdup */
#include <malloc.h> /* free for use with strdup */
#include <locale.h> /* number output formatting with commas */

int main(int argc, char **argv)
{
    atexit(exit_cleanup);
    setlocale(LC_NUMERIC, "");

    char *prog_name=strdup(argv[0]);
    if (prog_name == NULL) {
	    fprintf(stderr,"ERROR: unable to set program name\n");
    }

    get_cli_args(argc,argv);

    print_start_screen(prog_name);

    if (help) {
        show_help();
        return EXIT_SUCCESS;
    }

    check4DB();
    
    sqlite3_initialize();
    rc = sqlite3_open_v2(dbfile, &db, SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE, NULL);
    if (rc != SQLITE_OK) {
        exit(EXIT_FAILURE);
    }

    int totalrec = get_rec_count();
    printf(" - Current record count is: %'d\n",totalrec);

    char *lastacro = get_last_acronym();
    printf(" - Newest acronym is: %s\n",lastacro);
    if (lastacro != NULL ) {
	    free(lastacro);
    }

    /* perform a database acronym search */
    if ( findme != NULL ) {
	    int rec_match = 0;
	    rec_match = do_acronym_search(findme);
	    printf("\nDatabase search found '%'d' matching records\n",rec_match);
    }

    /* add a new acronym record */
    if (newrec) {
	    int add_worked = new_acronym();
	    if (add_worked) {
		    printf("\nADD DONE");
	    }
    }

    /* delete an acronym record */
    if (recordid >= 0 ) {
	    int del_worked = del_acro_rec(recordid);
	    if (del_worked) {
		    printf("\nDELETE DONE");
	    }
    }
    
    return (EXIT_SUCCESS);
}


void exit_cleanup(void)
{
    if (db == NULL) {
        printf("\nNo SQLite database shutdown required\n\nAll is well\n");
        exit(EXIT_SUCCESS);
    }

    rc = sqlite3_close_v2(db);
    if (rc != SQLITE_OK) {
        fprintf(stderr,"\nWARNING: error '%s' when trying to close the database\n",
                sqlite3_errstr(rc));
        exit(EXIT_FAILURE);
    }

    sqlite3_shutdown();
    printf("\nCompleted SQLite database shutdown\n\nAll is well\n");

    /* free any global varables below */
    if (findme != NULL) {
	    free(findme);
    }

    exit(EXIT_SUCCESS);
}


void print_start_screen(char *prog_name)
{
    printf(
    "\n"
    "\t\tAcronym Management Tool\n"
    "\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n"
    "Summary:\n"
    " - '%s' version is: %s complied with SQLite version: %s\n",
    prog_name,appversion, SQLITE_VERSION);

    /* done with this now */
    if (prog_name != NULL) {
	    free(prog_name);
    }

}


void show_help(void)
{
    printf(
   "\n"
   "Help Summary:\n"
   "The following command line switches can be used:\n"
   "\n"
   "  -d ?      Delete : remove an acronym where ? == ID of record to delete\n"
   "  -h        Help   : show this help information\n"
   "  -n        New    : add a new acronym record to the database\n"
   "  -s ?      Search : find an acronym where ? == acronym to search for\n"
   );
}

