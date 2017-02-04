/* Acronym Management Tool (amt): amt.c */

#include "amt-db-funcs.h"

#include <stdlib.h>		/* getenv */
#include <stdio.h>		/* printf */
#include <unistd.h>		/* strdup access stat and FILE */
#include <string.h>		/* strlen strdup */
#include <malloc.h>		/* free for use with strdup */
#include <locale.h>		/* number output formatting with commas */
#include <sys/types.h>		/* stat */
#include <sys/stat.h>		/* stat */
#include <time.h>		/* stat file modification time */
#include <readline/readline.h>	/* readline support for text entry */
#include <readline/history.h>	/* realine history support */

/*
 * Run SQL query to obtain current number of acronyms in the database.
 */
int get_rec_count(void)
{
	int totalrec = 0;
	rc = sqlite3_prepare_v2(db,"select count(*) from ACRONYMS"
				,-1,&stmt,NULL);

	if ( rc != SQLITE_OK) {
		exit(-1);
	}

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		totalrec = sqlite3_column_int(stmt,0);
	}

	sqlite3_finalize(stmt);
	return(totalrec);
}

/*
 * Check for a valid database file to open
 */
void check4DB(void)
{
	dbfile = getenv("ACRODB");
	if (dbfile) {
		printf(" - Database location: %s\n",dbfile);

		if (access(dbfile, F_OK | R_OK) == -1) {
			fprintf(stderr,"\n\nERROR: The database file '%s'"
				" is missing or is not accessible\n\n"
				,dbfile);
			exit(EXIT_FAILURE);
		}

		struct stat sb;
		int check;
    
		check = stat (dbfile, &sb);
    
		if (check) {
			perror("\nERROR: call to 'stat' for database file failed\n");
			exit(EXIT_FAILURE);
		}
    
		printf(" - Database size: %'ld bytes\n",sb.st_size);
		printf(" - Database last modified: %s\n",ctime(&sb.st_mtime));
	    
	} else {
		printf("\tWARNING: No database specified using 'ACRODB' "
		       "environment variable\n");
		exit(EXIT_FAILURE);
	}

/* TODO if neither of the above - check current directory we are
   running in - or then offer to create a new db? otherwise exit prog
   here */

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
	
	rc = sqlite3_prepare_v2(db,"SELECT Acronym FROM acronyms Order by rowid DESC LIMIT 1;"
				,-1,&stmt,NULL);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL prepare error: %s\n",sqlite3_errmsg(db));
		exit(-1);
	}

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		acronym_name = strdup((const char*)sqlite3_column_text(stmt,0));
	}

	sqlite3_finalize(stmt);

	if (acronym_name == NULL) {
		fprintf(stderr,"ERROR: last acronym lookup return NULL\n");
	}

	return(acronym_name);
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
	printf("\nSearching for: '%s' in database...\n\n",findme);

	rc = sqlite3_prepare_v2(db,"select rowid,Acronym,Definition,Description,"
				"Source from ACRONYMS where Acronym like ? "
				"COLLATE NOCASE ORDER BY Source;"
				,-1,&stmt,NULL);

	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL prepare error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	rc = sqlite3_bind_text(stmt,1,(const char*)findme,-1,SQLITE_STATIC);

	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL bind error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	int search_rec_count = 0;
	while(sqlite3_step(stmt) == SQLITE_ROW) {
		printf("ID:          %s\n"
		       ,(const char*)sqlite3_column_text(stmt,0));
		printf("ACRONYM:     '%s' is: %s.\n"
		       ,(const char*)sqlite3_column_text(stmt,1)
		       ,(const char*)sqlite3_column_text(stmt,2));
		printf("DESCRIPTION: %s\n"
		       ,(const char*)sqlite3_column_text(stmt,3));
		printf("SOURCE:      %s\n\n"
		       ,(const char*)sqlite3_column_text(stmt,4));
		search_rec_count++;
	}

	sqlite3_finalize(stmt);
	
	return search_rec_count;
}


/*
 * ADDING A NEW RECORD
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * insert into ACRONYMS(Acronym,Definition,Description,Source) values(?,?,?,?);
 * 
 */ 

int new_acronym(void)
{
	int old_rec_cnt = get_rec_count();
	char *source_list = get_acro_src();
	
	printf("\nAdding a new record...\n");
	printf("\nNote: To abort the input of a new record - press 'Ctrl + c'\n\n");

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

		printf("\n%s\n",source_list);
		
		n_acro_src = readline("\nEnter the acronym source: ");
		add_history(n_acro_src);

		printf("\nConfirm entry for:\n\n");
		printf("ACRONYM:     '%s' is: %s.\n",n_acro, n_acro_expd);
		printf("DESCRIPTION: %s\n", n_acro_desc);
		printf("SOURCE:      %s\n\n",n_acro_src);

		complete = readline("Enter record? [ y/n or q ] : ");
		if ( strcasecmp((const char*)complete,"y") == 0 ){
			break;
		}
		if ( strcasecmp((const char*)complete,"q") == 0 ){
			/* clean up memory allocations here */
			exit(EXIT_FAILURE);
		}

	}


	char *sql_ins = NULL;
	sql_ins = sqlite3_mprintf("insert into ACRONYMS"
				  "(Acronym, Definition, Description, Source) "
				  "values(%Q,%Q,%Q,%Q);"
				  ,n_acro,n_acro_expd,n_acro_desc,n_acro_src);
		
	rc = sqlite3_prepare_v2(db,sql_ins,-1,&stmt,NULL);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL prepare error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	rc = sqlite3_exec(db,sql_ins,NULL,NULL,NULL);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL exec error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}
	
	sqlite3_finalize(stmt);
	
	/* free up any allocated memory by sqlite3 */
	if(sql_ins !=NULL) { sqlite3_free(sql_ins); }

	/* free up any allocated memory by readline */
	if(n_acro != NULL) { free(n_acro); }
	if(n_acro_expd != NULL) { free(n_acro_expd); }
	if(n_acro_desc != NULL) { free(n_acro_desc); }
	if(n_acro_src != NULL) { free(n_acro_src); }

	int new_rec_cnt = get_rec_count();
	printf("Inserted '%d' new record. Total database record count is now"
	       " %'d (was %'d).\n"
	       ,(new_rec_cnt - old_rec_cnt),new_rec_cnt,old_rec_cnt);

	return 0;
}

/*
 * DELETE A RECORD BASE ROWID
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select rowid,Acronym,Definition,Description,Source from ACRONYMS where rowid = ?;
 *
 * delete from ACRONYMS where rowid = ?;
 * 
 */
int del_acro_rec(int recordid)
{
	int old_rec_cnt = get_rec_count();


	
	printf("\nDeleting an acronym record...\n");
	printf("\nNote: To abort the delete of a record - press 'Ctrl + c'\n\n");

	printf("\nSearching for record ID: '%d' in database...\n\n",recordid);

	rc = sqlite3_prepare_v2(db,"select rowid,Acronym,Definition,Description,"
				"Source from ACRONYMS where rowid like ?;"
				,-1,&stmt,NULL);

	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL prepare error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	rc = sqlite3_bind_int(stmt,1,recordid);
	if ( rc != SQLITE_OK) {
		fprintf(stderr,"SQL bind error: %s\n",sqlite3_errmsg(db));
		exit(EXIT_FAILURE);
	}

	int delete_rec_count = 0;
	while(sqlite3_step(stmt) == SQLITE_ROW) {
		printf("ID:          %s\n"
		       ,(const char*)sqlite3_column_text(stmt,0));
		printf("ACRONYM:     '%s' is: %s.\n"
		       ,(const char*)sqlite3_column_text(stmt,1)
		       ,(const char*)sqlite3_column_text(stmt,2));
		printf("DESCRIPTION: %s\n"
		       ,(const char*)sqlite3_column_text(stmt,3));
		printf("SOURCE: %s\n"
		       ,(const char*)sqlite3_column_text(stmt,4));
		delete_rec_count++;
	}

	sqlite3_finalize(stmt);

	if ( delete_rec_count > 0 ) {
		char *cont_del = NULL;
		cont_del = readline("\nDelete above record? [ y/n ] : ");
		if ( strcasecmp((const char*)cont_del,"y") == 0 ) {
		
			rc = sqlite3_prepare_v2(db,"delete from ACRONYMS where "
						"rowid = ?;"
						,-1,&stmt,NULL);
			if ( rc != SQLITE_OK) {
				fprintf(stderr,"SQL prepare error: %s\n",sqlite3_errmsg(db));
				exit(EXIT_FAILURE);
			}

			rc = sqlite3_bind_int(stmt,1,recordid);
			if ( rc != SQLITE_OK) {
				fprintf(stderr,"SQL bind error: %s\n",sqlite3_errmsg(db));
				exit(EXIT_FAILURE);
			}

			
			rc = sqlite3_step(stmt);
			if ( rc != SQLITE_DONE) {
				fprintf(stderr,"SQL step error: %s\n",sqlite3_errmsg(db));
				exit(EXIT_FAILURE);
			}
	
			sqlite3_finalize(stmt);
		}
	} else { printf(" » no record ID: '%d' found «\n\n",recordid); }

	int new_rec_cnt = get_rec_count();
	printf("Deleted '%d' record. Total database record count is now"
	       " %'d (was %'d).\n"
	       ,(old_rec_cnt - new_rec_cnt),new_rec_cnt,old_rec_cnt);
	
	return delete_rec_count;
	
}


/*
 * GETTING LIST OF ACRONYM SOURCES
 * ¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
 * select distinct(source) from acronyms;
 */

char *get_acro_src(void)
{
	rc = sqlite3_prepare_v2(db,"select distinct(source) from acronyms;"
				,-1,&stmt,NULL);

	if ( rc != SQLITE_OK) {
		exit(-1);
	}

	while(sqlite3_step(stmt) == SQLITE_ROW) {
		printf("%s\n",sqlite3_column_text(stmt,0));
	}

	sqlite3_finalize(stmt);

	char *returnme = "done";
	return(returnme);	
	
}
