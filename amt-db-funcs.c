/* Acronym Management Tool (amt): amt.c */

#include "amt-db-funcs.h"

#include <errno.h>             /* strerror */
#include <libgen.h>            /* basename and dirname */
#include <locale.h>            /* number output formatting with commas */
#include <malloc.h>            /* free for use with strdup and malloc */
#include <readline/history.h>  /* realine history support */
#include <readline/readline.h> /* readline support for text entry */
#include <stdio.h>             /* printf and asprintf*/
#include <stdlib.h>            /* getenv */
#include <string.h>            /* strlen strdup */
#include <sys/stat.h>          /* stat */
#include <sys/types.h>         /* stat */
#include <time.h>              /* stat file modification time */
#include <unistd.h>            /* strdup access stat and FILE */

/*
 * Run SQL query to obtain current number of acronyms in the database.
 */
int get_rec_count(void)
{
        int totalrec = 0;
        rc = sqlite3_prepare_v2(db, "select count(*) from ACRONYMS", -1, &stmt,
                                NULL);

        if (rc != SQLITE_OK) {
                perror("\nERROR: unable to access the SQLite database to "
                       "perform a record count\n");
                exit(EXIT_FAILURE);
        }

        while (sqlite3_step(stmt) == SQLITE_ROW) {
                totalrec = sqlite3_column_int(stmt, 0);
        }

        sqlite3_finalize(stmt);
        return (totalrec);
}

/*
 * Check for a valid database file to open
 */
void check4DB(char *prog_name)
{
        bool run_ok;

        /* get database file from environment variable ARCODB first */
        dbfile = getenv("ACRODB");
        /* if the environment variable exists - check if its valid */
        if (dbfile != NULL) {
                run_ok = check_db_access();
                if (run_ok) {
                        return;
                }
        }

        /* nothing is set in environment variable ARCODB - so database might
         * be found in the application directory instead */

        /* tmp copy needed here as each call to dirname() below can change the
         * string being used in the call - so need one string copy for each
         * successful call we need to make. This is a 'feature' of dirname() */
        char *tmp_dirname = strdup(prog_name);

        size_t new_dbfile_sz = (sizeof(char) * (strlen(dirname(tmp_dirname)) +
                                                strlen("/acronyms.db") + 1));

        char *new_dbfile = malloc(new_dbfile_sz);

        if (new_dbfile == NULL) {
                perror("\nERROR: unable to allocate memory with "
                       "malloc() for 'new_dbfile' and path\n");
                exit(EXIT_FAILURE);
        }

        int x = snprintf(new_dbfile, new_dbfile_sz, "%s%s", dirname(prog_name),
                         "/acronyms.db");

        if (x == -1) {
                perror("\nERROR: unable to allocate memory with "
                       "snprintf() for 'new_dbfile' and path\n");
                exit(EXIT_FAILURE);
        }

        if ((dbfile = strdup(new_dbfile)) == NULL) {
                perror("\nERROR: unable to allocate memory with "
                       "strdup() for 'new_dbfile' to 'dbfile' copy\n");
                exit(EXIT_FAILURE);
        }

        printf("\nnew_dbfile: '%s' and dbfile: '%s'\n", new_dbfile, dbfile);

        if (new_dbfile != NULL) {
                free(new_dbfile);
        }

        /* now recheck if the new_dbfile is suitable for use? */
        run_ok = check_db_access();
        if (run_ok) {
                return;
        }

        /* run out of options to find a suitable database - exit */
        printf("\n\tWARNING: No suitable database file can be located - "
               "program will exit\n");
        exit(EXIT_FAILURE);
}

/*
 * Check the filename and path given for the acronym database and see if it is
 * accessable. This file and patch is stored in the global variable: 'dbfile''
 *
 */
bool check_db_access(void)
{
        if (dbfile == NULL || strlen(dbfile) == 0) {
                fprintf(stderr, "ERROR: The database file '%s'"
                                " is an empty string\n",
                        dbfile);
                return (false);
        }

        if (access(dbfile, F_OK | R_OK) == -1) {
                fprintf(stderr, "ERROR: The database file '%s'"
                                " is missing or is not accessible\n",
                        dbfile);
                return (false);
        }

        printf(" - Database location: %s\n", dbfile);

        struct stat sb;
        int check;

        check = stat(dbfile, &sb);

        if (check) {
                perror("ERROR: call to 'stat' for database file "
                       "failed\n");
                return (false);
        }

        printf(" - Database size: %'ld bytes\n", sb.st_size);
        printf(" - Database last modified: %s\n", ctime(&sb.st_mtime));

        return (true);
}

/*
 * GET NAME OF LAST ACRONYM ENTERED
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;
 *
 */
char *get_last_acronym(void)
{
        char *acronym_name;

        rc = sqlite3_prepare_v2(
            db, "SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;", -1,
            &stmt, NULL);
        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL prepare error: %s\n", sqlite3_errmsg(db));
                exit(-1);
        }

        while (sqlite3_step(stmt) == SQLITE_ROW) {
                acronym_name =
                    strdup((const char *)sqlite3_column_text(stmt, 0));
        }

        sqlite3_finalize(stmt);

        if (acronym_name == NULL) {
                fprintf(stderr, "ERROR: last acronym lookup return NULL\n");
        }

        return (acronym_name);
}

/*
 * SEARCH FOR A NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select rowid,Acronym,Definition,
 * Description,Source from ACRONYMS
 * where Acronym like ? COLLATE NOCASE ORDER BY Source;
 *
 */

int do_acronym_search(char *findme)
{
        printf("\nSearching for: '%s' in database...\n\n", findme);

        rc = sqlite3_prepare_v2(db,
                                "select rowid,Acronym,Definition,Description,"
                                "Source from ACRONYMS where Acronym like ? "
                                "COLLATE NOCASE ORDER BY Source;",
                                -1, &stmt, NULL);

        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL prepare error: %s\n", sqlite3_errmsg(db));
                exit(EXIT_FAILURE);
        }

        rc =
            sqlite3_bind_text(stmt, 1, (const char *)findme, -1, SQLITE_STATIC);

        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL bind error: %s\n", sqlite3_errmsg(db));
                exit(EXIT_FAILURE);
        }

        int search_rec_count = 0;
        while (sqlite3_step(stmt) == SQLITE_ROW) {
                printf("ID:          %s\n",
                       (const char *)sqlite3_column_text(stmt, 0));
                printf("ACRONYM:     '%s' is: %s.\n",
                       (const char *)sqlite3_column_text(stmt, 1),
                       (const char *)sqlite3_column_text(stmt, 2));
                printf("DESCRIPTION: %s\n",
                       (const char *)sqlite3_column_text(stmt, 3));
                printf("SOURCE:      %s\n\n",
                       (const char *)sqlite3_column_text(stmt, 4));
                search_rec_count++;
        }

        sqlite3_finalize(stmt);

        return search_rec_count;
}

/*
 * ADDING A NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * insert into ACRONYMS(Acronym,Definition,Description,Source)
 * values(?,?,?,?);
 *
 */

int new_acronym(void)
{
        int old_rec_cnt = get_rec_count();

        printf("\nAdding a new record...\n");
        printf("\nNote: To abort the input of a new record - press "
               "'Ctrl + "
               "c'\n\n");

        char *complete = NULL;
        char *n_acro = NULL;
        char *n_acro_expd = NULL;
        char *n_acro_desc = NULL;
        char *n_acro_src = NULL;

        while (1) {
                n_acro = readline("Enter the acronym: ");
                add_history(n_acro);
                n_acro_expd = readline("Enter the expanded acronym: ");
                add_history(n_acro_expd);
                n_acro_desc = readline("Enter the acronym description: \n\n");
                add_history(n_acro_desc);

                get_acro_src();
                n_acro_src = readline("\nEnter the acronym source: ");
                add_history(n_acro_src);

                printf("\nConfirm entry for:\n\n");
                printf("ACRONYM:     '%s' is: %s.\n", n_acro, n_acro_expd);
                printf("DESCRIPTION: %s\n", n_acro_desc);
                printf("SOURCE:      %s\n\n", n_acro_src);

                complete = readline("Enter record? [ y/n or q ] : ");
                if (strcasecmp((const char *)complete, "y") == 0) {
                        break;
                }
                if (strcasecmp((const char *)complete, "q") == 0) {
                        /* Clean up readline allocated memory */
                        if (complete != NULL) {
                                free(complete);
                        }
                        if (n_acro != NULL) {
                                free(n_acro);
                        }
                        if (n_acro_expd != NULL) {
                                free(n_acro_expd);
                        }
                        if (n_acro_desc != NULL) {
                                free(n_acro_desc);
                        }
                        if (n_acro_src != NULL) {
                                free(n_acro_src);
                        }
                        rl_clear_history();
                        exit(EXIT_FAILURE);
                }
        }

        char *sql_ins = NULL;
        sql_ins = sqlite3_mprintf("insert into ACRONYMS"
                                  "(Acronym, Definition, Description, Source) "
                                  "values(%Q,%Q,%Q,%Q);",
                                  n_acro, n_acro_expd, n_acro_desc, n_acro_src);

        rc = sqlite3_prepare_v2(db, sql_ins, -1, &stmt, NULL);
        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL prepare error: %s\n", sqlite3_errmsg(db));
                /* Clean up readline allocated memory */
                if (complete != NULL) {
                        free(complete);
                }
                if (n_acro != NULL) {
                        free(n_acro);
                }
                if (n_acro_expd != NULL) {
                        free(n_acro_expd);
                }
                if (n_acro_desc != NULL) {
                        free(n_acro_desc);
                }
                if (n_acro_src != NULL) {
                        free(n_acro_src);
                }
                rl_clear_history();
                exit(EXIT_FAILURE);
        }

        rc = sqlite3_exec(db, sql_ins, NULL, NULL, NULL);
        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL exec error: %s\n", sqlite3_errmsg(db));
                /* Clean up readline allocated memory */
                if (complete != NULL) {
                        free(complete);
                }
                if (n_acro != NULL) {
                        free(n_acro);
                }
                if (n_acro_expd != NULL) {
                        free(n_acro_expd);
                }
                if (n_acro_desc != NULL) {
                        free(n_acro_desc);
                }
                if (n_acro_src != NULL) {
                        free(n_acro_src);
                }
                rl_clear_history();
                exit(EXIT_FAILURE);
        }

        sqlite3_finalize(stmt);

        /* free up any allocated memory by sqlite3 */
        if (sql_ins != NULL) {
                sqlite3_free(sql_ins);
        }

        /* Clean up readline allocated memory */
        if (complete != NULL) {
                free(complete);
        }
        if (n_acro != NULL) {
                free(n_acro);
        }
        if (n_acro_expd != NULL) {
                free(n_acro_expd);
        }
        if (n_acro_desc != NULL) {
                free(n_acro_desc);
        }
        if (n_acro_src != NULL) {
                free(n_acro_src);
        }
        rl_clear_history();

        int new_rec_cnt = get_rec_count();
        printf("Inserted '%d' new record. Total database record count "
               "is now"
               " %'d (was %'d).\n",
               (new_rec_cnt - old_rec_cnt), new_rec_cnt, old_rec_cnt);

        return 0;
}

/*
 * DELETE A RECORD BASE ROWID
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select rowid,Acronym,Definition,Description,Source from ACRONYMS
 * where rowid
 * = ?;
 *
 * delete from ACRONYMS where rowid = ?;
 *
 */
int del_acro_rec(int recordid)
{
        int old_rec_cnt = get_rec_count();
        printf("\nDeleting an acronym record...\n");
        printf("\nNote: To abort the delete of a record - press 'Ctrl "
               "+ c'\n\n");

        printf("\nSearching for record ID: '%d' in database...\n\n", recordid);

        rc = sqlite3_prepare_v2(db,
                                "select rowid,Acronym,Definition,Description,"
                                "Source from ACRONYMS where rowid like ?;",
                                -1, &stmt, NULL);

        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL prepare error: %s\n", sqlite3_errmsg(db));
                exit(EXIT_FAILURE);
        }

        rc = sqlite3_bind_int(stmt, 1, recordid);
        if (rc != SQLITE_OK) {
                fprintf(stderr, "SQL bind error: %s\n", sqlite3_errmsg(db));
                exit(EXIT_FAILURE);
        }

        int delete_rec_count = 0;
        while (sqlite3_step(stmt) == SQLITE_ROW) {
                printf("ID:          %s\n",
                       (const char *)sqlite3_column_text(stmt, 0));
                printf("ACRONYM:     '%s' is: %s.\n",
                       (const char *)sqlite3_column_text(stmt, 1),
                       (const char *)sqlite3_column_text(stmt, 2));
                printf("DESCRIPTION: %s\n",
                       (const char *)sqlite3_column_text(stmt, 3));
                printf("SOURCE: %s\n",
                       (const char *)sqlite3_column_text(stmt, 4));
                delete_rec_count++;
        }

        sqlite3_finalize(stmt);

        if (delete_rec_count > 0) {
                char *cont_del = NULL;
                cont_del = readline("\nDelete above record? [ y/n ] : ");
                if (strcasecmp((const char *)cont_del, "y") == 0) {

                        rc =
                            sqlite3_prepare_v2(db, "delete from ACRONYMS where "
                                                   "rowid = ?;",
                                               -1, &stmt, NULL);
                        if (rc != SQLITE_OK) {
                                fprintf(stderr, "SQL prepare error: %s\n",
                                        sqlite3_errmsg(db));
                                if (cont_del != NULL) {
                                        free(cont_del);
                                }
                                exit(EXIT_FAILURE);
                        }

                        rc = sqlite3_bind_int(stmt, 1, recordid);
                        if (rc != SQLITE_OK) {
                                fprintf(stderr, "SQL bind error: %s\n",
                                        sqlite3_errmsg(db));
                                if (cont_del != NULL) {
                                        free(cont_del);
                                }
                                exit(EXIT_FAILURE);
                        }

                        rc = sqlite3_step(stmt);
                        if (rc != SQLITE_DONE) {
                                fprintf(stderr, "SQL step error: %s\n",
                                        sqlite3_errmsg(db));
                                if (cont_del != NULL) {
                                        free(cont_del);
                                }
                                exit(EXIT_FAILURE);
                        }

                        /* free readline memory allocated */
                        if (cont_del != NULL) {
                                free(cont_del);
                        }
                        sqlite3_finalize(stmt);
                }
        } else {
                printf(" » no record ID: '%d' found «\n\n", recordid);
        }

        int new_rec_cnt = get_rec_count();
        printf("Deleted '%d' record. Total database record count is now"
               " %'d (was %'d).\n",
               (old_rec_cnt - new_rec_cnt), new_rec_cnt, old_rec_cnt);

        return delete_rec_count;
}

/*
 * GETTING LIST OF ACRONYM SOURCES
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select distinct(source) from acronyms;
 */

void get_acro_src(void)
{
        rc = sqlite3_prepare_v2(db, "select distinct(source) from acronyms;",
                                -1, &stmt, NULL);

        if (rc != SQLITE_OK) {
                exit(-1);
        }

        char *acro_src_name;

        printf("\nSelect a source (use ↑ or ↓ ):\n\n");

        while (sqlite3_step(stmt) == SQLITE_ROW) {
                acro_src_name =
                    strdup((const char *)sqlite3_column_text(stmt, 0));
                printf("[ %s ] ", acro_src_name);
                add_history(acro_src_name);

                /* free per loop to stop memory leaks - strdup malloc
                 * above */
                if (acro_src_name != NULL) {
                        free(acro_src_name);
                }
        }
        printf("\n");

        sqlite3_finalize(stmt);
}