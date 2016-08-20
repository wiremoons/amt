/* Acronym Management Tool (amt): cli-args.h */

#ifndef CLI_ARGS_H_ /* Include guard */
#define CLI_ARGS_H_ 

#include <stdlib.h>     /* getenv */
#include <stdio.h>      /* */
#include <malloc.h>     /* free for use with strdup */
#include <unistd.h>     /* strdup access */
#include <string.h>     /* strlen */

extern int argc;
extern char **argv;
extern char *findme;
extern int debug;
extern int newrec;
extern int help;
extern char appversion[];

/*
 *   FUNCTION DECLARATIONS FOR cli-args.c
 */

void getCLIArgs(int argc, char **argv);

#endif // AMT_H_ 

