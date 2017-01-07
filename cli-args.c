/* Acronym Management Tool (amt): cli-args.c */

#include "cli-args.h"

#include <stdlib.h> /* getenv */
#include <stdio.h>  /* printf */
#include <malloc.h> /* free for use with strdup */
#include <unistd.h> /* strdup access */
#include <string.h> /* strlen */
#include <ctype.h>  /* isdigit */

/*
 * Used to parse command line options provided by the user.
 * Uses the POSIX compliant getopts()
 */
void get_cli_args(int argc, char **argv)
{
    opterr = 0;
    int c = 0;

    while ((c = getopt (argc, argv, "d:hns:")) != -1) {
        switch (c) {
        case 'd':
                recordid = *optarg;
                if ( recordid < 0 || !isdigit(recordid) ) {
                    fprintf(stderr,"\nERROR: for -d option please provide "
                            "an acronym ID for removal. Use search function to locate correct 'ID' first\n");
                    exit(EXIT_FAILURE);
                }
                break;
        case 'h':
		help = 1;
		break;
        case 'n':
                newrec = 1;
                break;
        case 's':
                findme = strdup(optarg);
                if ( strlen(findme) <= 0 || findme == NULL ) {
                    fprintf(stderr,"\nERROR: for -s option please provide "
                            "an acronym to search for\n");
                    exit(EXIT_FAILURE);
                }
                break;
        case ':':
                fprintf(stderr,
			"\nERROR: '%s' option '-%c' requires an argument\n"
			,argv[0], optopt);
		exit(EXIT_FAILURE);
                break;
       case '?':
       default:
                fprintf(stderr,
			"\nERROR: '%s' option '-%c' is invalid or missing input data\n",
                        argv[0], optopt);
		exit(EXIT_FAILURE);
                break;
        }
    }
}

