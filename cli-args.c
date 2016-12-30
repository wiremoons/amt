/* Acronym Management Tool (amt): cli-args.c */

#include "cli-args.h"

#include <stdlib.h> /* getenv */
#include <stdio.h>  /* printf */
#include <malloc.h> /* free for use with strdup */
#include <unistd.h> /* strdup access */
#include <string.h> /* strlen */

/*
 * Used to parse command line options provided by the user.
 * Uses the POSIX compliant getopts()
 */
void get_cli_args(int argc, char **argv)
{
    opterr = 0;
    int c = 0;

    while ((c = getopt (argc, argv, "hvns:")) != -1) {
        switch (c) {
        case 'h':
		help = 1;
		break;
        case 'v':
                printf("%s is version: %s\n",argv[0],appversion);
                exit(EXIT_SUCCESS);
        case 'n':
                newrec = 1;
                break;
        case 's':
                findme = strdup(optarg);
                if ( strlen(findme) <= 0 || findme == NULL ) {
                    fprintf(stderr,"ERROR: for -s option please provide "
                            "an acronym to search for");
                    exit(EXIT_FAILURE);
                }
                break;
        case ':':
                fprintf(stderr,
			"ERROR: '%s' option '-%c' requires an argument\n"
			,argv[0], optopt);
                break;
       case '?':
       default:
                fprintf(stderr,
			"ERROR: '%s' option '-%c' is invalid or missing input data\n",
                        argv[0], optopt);
                break;
        }
    }
}

