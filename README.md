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

The `amt` program accesses a SQLite database and looks up the requested acronym held in a table called '*ACRONYMS*'.

## Database Location

By default the database used to store the acronyms will be located in the same
directory as the programs executable.

However this can be overridden if preferred, and the location of the database
can be stored in an environment variable called *ACRODB*. You should set this
to the path and preferred database file name of your acronyms database.
Examples of how to set this for different operating systems are shown below.

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

### Using amt

When you run `amt` from the command line it outputs the following
details on the screen, assuming you have a valid database already
setup that can be used with the program.


```

			Acronym Management Tool 'amt'
			¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Database: Sybil.db   permissions: -rw-rw-r--   size: 2,027,520 bytes

Database connection status:  √
SQLite3 Database Version:  3.14.0
Current record count is:  17,137
Last acronym entered was:  'SNI'

Running 'amt' version 0.5.5

 - Built with Go Complier 'gc' on Golang version 'go1.7.1'
 - Author's web site: http://www.wiremoons.com/
 - Source code for amt: https://github.com/wiremoons/amt/


Usage of ./amt:

        Flag               Description                                        Default Value
        ¯¯¯¯               ¯¯¯¯¯¯¯¯¯¯¯                                        ¯¯¯¯¯¯¯¯¯¯¯¯¯
        -d                 show debug output                                  false
        -f <filename>      provide filename and path to SQLite database       optional
        -h                 display help for this program                      false
        -n                 add a new acronym record                           optional
        -s <acronym>       provide acronym to search for                      optional
        -r <acronym id>    provide acronym id to remove                       optional
        -v                 display program version                            false
        -w                 search for any similar matches                     false


All is well
```

If you perform a search for an acronym that exists in the database,
the following would be an example output, where more than one record
matches the search. The command used for the output below was `amt -s
sni`:

```

			Acronym Management Tool 'amt'
			¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯
Database: Sybil.db   permissions: -rw-rw-r--   size: 2,027,520 bytes

Database connection status:  √
SQLite3 Database Version:  3.14.0
Current record count is:  17,137
Last acronym entered was:  'SNI'


SEARCH FOR AN ACRONYM RECORD
¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯¯

Searching for:  'sni'  across 17,137 records - please wait...

Matching results are:

ID: 14306
ACRONYM: 'SNI' is: Serious Network Incident.
DESCRIPTION: 
SOURCE: General ICT

ID: 14307
ACRONYM: 'SNI' is: Switch Network Interconnect.
DESCRIPTION: The Switch Network Interconnect (SNI) Zone contains the Point of Interconnect (POI) with the Other Licensed Operators (OLO)
SOURCE: General ICT

ID: 14308
ACRONYM: 'SNI' is: Service Node Interface.
DESCRIPTION: 
SOURCE: General ICT

ID: 14309
ACRONYM: 'SNI' is: Subscriber Network Interface.
DESCRIPTION: 
SOURCE: General ICT

ID: 17137
ACRONYM: 'SNI' is: Server Name Indication.
DESCRIPTION: Name-based virtual hosting allows multiple DNS hostnames to be hosted by a single server (usually a web server) on the same IP address. To achieve this the server uses a hostname presented by the client as part of the protocol (for HTTP the name is presented in the host header). However, when using HTTPS the TLS handshake happens before the server sees any HTTP headers. Therefore, it is not possible for the server to use the information in the HTTP host header to decide which certificate to present and as such only names covered by the same certificate can be served from the same IP address.
SOURCE: General ICT


All is well
```



## License

This program is licensed under the "MIT License" see
http://opensource.org/licenses/mit for more details.

