[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)

# What is 'amt'?

A small command line application called '`amt`' (which is an acronym for
'*Acronym Management Tool*') which can be used to store, look up, change, or
delete acronyms that are kept in a SQLite database.

The program is small, fast, and is free software. It is used on a daily basis by
the author, running it on both Linux and Windows operating systems. It should
also compile and run on BSD Unix too, and Mac OS X as well, although this has
not been tested.


## Usage Examples

Running `amt` without any parameters, but with a database already setup will
output the following information:

```
		Acronym Management Tool
		¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Summary:
 - 'amt' version is: 0.4.4 complied with SQLite version: 3.17.0
 - Database location: /home/simon/Work/Scratch/Sybil/Sybil.db
 - Database size: 1,851,392 bytes
 - Database last modified: Sun Mar 19 07:57:30 2017

 - Current record count is: 17,268
 - Newest acronym is: ASF

Completed SQLite database shutdown

All is well
```

Running `amt -h` displays the help screen which will output the following
information:

```
		Acronym Management Tool
		¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Summary:
 - 'amt' version is: 0.4.4 complied with SQLite version: 3.17.0

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

A C language compiler will be needed to build the application. There are also a
few dependencies that will need to be met for a successful build of `amt`.
These steps are explained below.

### Dependencies

In order to compile `amt` successfully, it requires a few development C code
libraries to support its functionality. These libraries are described below for
reference.

#### SQLite

The application uses SQLite, however it **includes copies** of the `sqlite3.c` and
`sqlite3.h` source code files within this applications own code distribution.
These are are compiled directly into the application when it is built. The two
includes files are obtained from [http://www.sqlite.org/amalgamation.html](The
SQLite Amalgamation) source download option.

More information on SQLite can be found here: http://www.sqlite.org/

#### Readline

The application uses GNU Readline, and requires these development libraries to
be installed on the system that `amt` is being built on. On most Linux systems
the Readline library can be installed from the distributions software
repositories. The libraries will need to be available before `amt` will compile
successfully. The GNU Readline library is used to provide some of the
applications functionality during text entry.

More information on GNU Readline Library can be found here:
http://cnswww.cns.cwru.edu/php/chet/readline/readline.html

### Install a C Compiler and Supporting Libraries

To install the required libraries and compiler tools on various systems, use the
following commands before attempting to compile `amt`:

- Ubuntu Linux: `sudo apt install build-essential libreadline6 libreadline6-dev`
- Fedora (Workstation) Linux: `sudo dnf install readline-devel `

On Windows you can use a C compiler such as MinGW (or equivalent) to build the
application. In order to build the Windows version for testing and personal use,
this is done on Fedora Workstation as a cross compile.

### Building 'amt'

Before trying to build `amt` make sure you have the required dependencies install
- see above for more information.

An example command to compile `amt` with GCC compiler on a 64bit Linux system is shown below:
```
gcc -g -Wall -m64 -std=gnu11 -o amt amt-db-funcs.c cli-args.c main.c sqlite3.c -lpthread -ldl -lreadline
```

Alternatively, the program can be compiled using the provided '*Makefile*'. For
example to compile a optimised (non debug) version of `amt`, run:
```
make opt
```


## Database Location

By default the SQLite database used to store the acronyms can be located in the
same directory as the programs executable.

However this can be overridden, and the preferred location can be specified by
an environment variable called ***ACRODB***. You should set this to the path and
preferred database file name of your SQLite acronyms database. Examples of how
to set this for different operating systems are shown below.

On Linux and similar operating systems when using bash shell, add this line to
your `.bashrc` configuration file, located in your home directory (ie
*~/.bashrc*), just amend the path and database file name to suit your own needs:

```
export acrodb=$HOME/work/acrotool/sybil.db
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

The following license apply to the `amt` source code and resulting built
application, as described below.

#### License for 'amt'

This program `amt` is licensed under the **MIT License** see
http://opensource.org/licenses/mit for more details.

#### License for 'SQLite'

The SQLite database code used in this application is licensed as **Public
Domain**, see http://www.sqlite.org/copyright.html for more details.

#### License for 'GNU Readline'

The GNU Readline Library used in this application is licensed under the **GNU
GPL v3**, see https://www.gnu.org/licenses/ for more details. 

### License Limitations

**PLEASE NOTE**: Any *distributed compiled binary* versions of `amt` will
include GNU Readline functionality, and these specific binary versions will
therefore be 'infected' by the GPL license. In these circumstances, if the
binary is onward distributed (ie shared with another party), it should be
licensed as 'GNU GPL v3'. For most users of other GPL licensed software (such as
GNU Linux distributions), this will not be an issue or concern either.

The `amt` source code remains completely free under a the **MIT License**
however. This GNU GPL license *developer limitation of freedom* can be
circumvented by perhaps using an alternative to the GNU Readline library, such
as:

- [http://thrysoee.dk/editline/](EditLine) 
- [https://github.com/antirez/linenoise](Linenoise)

These two example alternative libraries are **not** GNU GPL licensed libraries.
These options provide an alternative *choice*, should you want to invoke it.