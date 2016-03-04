# Summary for `amt`

A small application called 'amt' which is an abbreviation for 'Acronym
Management Tool' that can be used to store, look up, and change or delete
acronyms that are held in a SQLite database.

## About

`amt` (short for '*Acronym Management Tool*') is used to manage a list of
ancronyms that are held in a local SQLite database table.

The program can search for existing acronyms, add or delete acronyms, and
amend existsing acronyms which are stored in the SQLite database.

## Building the Application

A c compiler will be needed to build the application.

Compile with: `gcc -Wall --std=gnu11 amt.c sqlite3.c -o amt` or use the
provided 'Makefile', and run one of the options such as `make opt`.

## License

The program is licensed under the "MIT License" see
http://opensource.org/licenses/MIT for more details.
