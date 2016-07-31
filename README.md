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

## Building the Application

A C language compiler will be needed to build the application.

Compile with: `gcc -Wall --std=gnu11 amt.c sqlite3.c -o amt` or use the
provided 'Makefile', and run one of the options such as `make opt`.

## Dependencies

The application uses SQLite, and includes copies of the `sqlite3.c` and
`sqlite3.h` code files with the this applications own code. These are are
compiled into the application when it is built.

More information on SQLite can be found here: http://www.sqlite.org/

## Database Location

By default the database used to store the acronyms will be located in the same
directory as the programs executable.

However this can be override if prefer, and the location of the database can be
stored in an environment variable called acrodb. you should set this to the
path and preferred database file name of you acronyms database. Example of
setting this are shown below.

On Linux and similar operating systems when using bash:

```
export acrodb=/home/simon/work/acrotool/sybil.db
```

on windows or linux when using microsoft powershell:

```
$env:acrodb += "c:\users\simon\work\scratch\sybil\sybil.db"
```

on windows when using a cmd.exe console:

```
set acrodb=c:\users\simon\work\scratch\sybil\sybil.db
```

or windows to add persistently to your environment run the following in a
cmd.exe console:

```
setx acrodb=c:\users\simon\work\scratch\sybil\sybil.db
```

## license

this program is licensed under the "mit license" see
http://opensource.org/licenses/mit for more details.

