[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)

# Summary for `amt`

A small application called '`amt`' which is an acronym for 'Acronym Management
Tool' that can be used to store, look up, and change or delete acronyms that
are held in a SQLite database.

## About

`amt` (short for '*Acronym Management Tool*') is used to manage a list of
acronyms that are held in a local SQLite database table.

The program can search for acronyms, add or delete acronyms, and amend
existing acronyms which are all stored in the SQLite database.


## Usage Examples

Running `amt` without any parameters, but with a database already setup will
output the following information:

```
		Acronym Management Tool
		¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Summary:
 - 'amt' version is: 0.4.2 complied with SQLite version: 3.15.2
 - Database location: /home/simon/Work/Scratch/Sybil/Sybil.db
 - Database size: 2,076,672 bytes
 - Database last modified: Sat Feb 11 09:06:55 2017

 - Current record count is: 17,239
 - Newest acronym is: CLOS

Completed SQLite database shutdown

All is well
```

Running `amt -h` displays the help screen which will output the following
information:

```
		Acronym Management Tool
		¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Summary:
 - 'amt' version is: 0.4.2 complied with SQLite version: 3.15.2

Help Summary:
The following command line switches can be used:

  -d ?      Delete : remove an acronym where ? == ID of record to delete
  -h        Help   : show this help information
  -n        New    : add a new acronym record to the database
  -s ?      Search : find an acronym where ? == acronym to search for

No SQLite database shutdown required

All is well
```


## Building the Application

A C language compiler will be needed to build the application.

Compile with on a 64bit Linux system with:
```
gcc -g -Wall -m64 -std=gnu11 -o amt amt-db-funcs.c cli-args.c main.c sqlite3.c -lpthread -ldl -lreadline
```
 or use the provided 'Makefile', and run one of the options such as `make opt`.

## Dependencies

#### SQLite

The application uses SQLite, and includes copies of the `sqlite3.c` and
`sqlite3.h` code files within this applications own code distribution. These are
are compiled directly into the application when it is built.

More information on SQLite can be found here: http://www.sqlite.org/

#### Readline

The application uses GNU Readline, and requires these development libraries to be
installed on the system that `amt` is being built on. On most Linux systems the
Readline library can be installed from the distributions software repositories.

More information on GNU Readline Library can be found here:
http://cnswww.cns.cwru.edu/php/chet/readline/readline.html


## Database Location

By default the SQLite database used to store the acronyms will be located in the same
directory as the programs executable.

However this can be overridden if preferred, and the location the database
should be stored in, can be specified by an environment variable called
*ACRODB*. You should set this to the path and preferred database file name of
your acronyms database. Examples of how to set this for different operating
systems are shown below.

On Linux and similar operating systems when using bash shell:

```
export acrodb=/home/simon/work/acrotool/sybil.db
```

on Windows or Linux when using Microsoft Powershell:

```
$env:acrodb += "c:\users\simon\work\scratch\sybil\sybil.db"
```

on Windows when using a cmd.exe console:

```
set acrodb=c:\users\simon\work\scratch\sybil\sybil.db
```

or Windows to add persistently to your environment run the following in a
cmd.exe console:

```
setx acrodb=c:\users\simon\work\scratch\sybil\sybil.db
```

## Licenses

This program `amt` is licensed under the **MIT License** see
http://opensource.org/licenses/mit for more details.

The GNU Readline Library used in this application is licensed under the **GNU
GPL v3**, see http://www.gnu.org/licenses/ for more details.

The SQLite database code used in this application is licensed a **Public
Domain**, see http://www.sqlite.org/copyright.html for more details.

