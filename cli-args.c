/* Acronym Management Tool (amt): cli-args.c */

#include "cli-args.h"

/**-------- FUNCTION: getCLIArgs
**
**  Function called when program starts. Used to parse command line
**  options provided by the user. Uses the POSIX compliant getopts()
**
*/
void getCLIArgs(int argc, char **argv)
{
    opterr = 0;
    int c = 0;
    int index = 0;

    while ((c = getopt (argc, argv, "dhvns:")) != -1)
    {
        switch (c)
        {
            // debugging output was requested
            case 'd':
                debug = 1;
                break;
                // provide help summary to user was requested
            case 'h':
                help = 1;
                break;
                // the version of the application requested
            case 'v':
                printf("%s is version: %s\n",argv[0],appversion);
                exit(EXIT_SUCCESS);
                // request to add a new acronym record
            case 'n':
                newrec = 1;
                break;
                // request to search for an acronym
            case 's':
                findme = strdup(optarg);
                if ( strlen(findme) <= 0 || findme == NULL )
                {
                    fprintf(stderr,"ERROR: for -s option please provide "
                            "an acronym to search for");
                    exit(EXIT_FAILURE);
                }
                break;
                // ERROR HANDLING BELOW
                //
                // command line option given - but is missing the required data argument for it
            case ':':
                fprintf(stderr,"ERROR: '%s' option '-%c' requires an argument\n",argv[0], optopt);
                break;
                // invalid option provided on command line - also 'default' as the switch fall-thru
            case '?':
            default:
                /* invalid option */
                fprintf(stderr,"ERROR: '%s' option '-%c' is invalid or missing input data\n",
                        argv[0], optopt);
                break;
        }
    }
    // if debugging requested - display extra getopt() info
    if (debug)
    {
        printf ("DEBUG: optargs() values:\n"
                "\tdebug = %s\n"
                "\tHelp requested = %s\n"
                "\tNew record input requested = %s\n"
                "\tSearch requested for string = %s\n\n",
                debug ? "true" : "false",
                help ? "true" : "false",
                newrec ? "true" : "false",
                findme);
        for (index = optind; index < argc; index++)
        {
            printf ("\tInvalid option argument(s) seen: %s\n", argv[index]);
        }
    }
}


