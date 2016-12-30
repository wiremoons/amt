/* Acronym Management Tool (amt): main.c */

#include "main.h"

#include <stdlib.h> /* getenv */
#include <stdio.h>  /* printf */
#include <unistd.h> /* strdup access */
#include <string.h> /* strlen */
#include <malloc.h> /* free for use with strdup */
#include <locale.h> /* number output formatting with commas */

int main(int argc, char **argv)
{
    atexit(exit_cleanup);
    setlocale(LC_NUMERIC, "");
    
    get_cli_args(argc,argv);

    print_start_screen();

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

    int totalrec = recCount();
    printf(" - Record count is: %'d\n",totalrec);

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

    if (findme != NULL) {
	    free(findme);
    }
    
    exit(EXIT_SUCCESS);
}


void print_start_screen(void)
{
    printf(
    "\n"
    "\t\tAcronym Management Tool\n"
    "\t\t¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯\n"
    "Summary:\n"
    " - App version: %s complied with SQLite version: %s\n",
    appversion, SQLITE_VERSION);
}


void show_help(void)
{
    printf(
   "\n"
   "Help Summary:\n"
   "The following command line switches can be used:\n\n"
   "  -d\tDebug - include additional debug outputs when run\n"
   "  -h\tHelp - Show this help information\n"
   "  -v\tVersion - display the version of the program\n"
   "  -n\tNew - add a new acronym record to the database\n"
   "  -s ?\tSearch - find an acronym where ? == acronym to search for\n"
   );
}

