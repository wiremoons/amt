#.............................................................................
#
#   Makefile Template.
#       Originally Created: June 2013 by Simon Rowe <simon@wiremoons.com>
#.............................................................................
#
#	Makefile for amt.c
#
#   gcc -Wall -std=gnu11 -m64 -g -o amt amt.c cli-args.c main.c sqlite3.c -Lpthread -ldl
#
## CHANGE xxxx FOR YOUR NEW SOURCE FILE NAME & OUTPUT FILE-NAME
SRC=amt.c cli-args.c main.c sqlite3.c
OUTNAME=amt
#
#  NOTE:
#  The settings below assume Microsoft Windows using MinGW Compiler as a
#  starting point. Also assumes 32bit MinGW compiler. This done so:
#     1. we start with the lowest required settings to work on Windows
#     2. The more comprehensive tools available on Linux/MacOSX/*BSD etc
#        systems can dynamically figure out actual 'best' requirements, 
#        if it turns out 'make' is not actually being run from Windows.
#  The above then negates Windows being required to have lots of other
#  tools installed to achieve what other operating systems can do more
#  easily. This also means the default Windows system only needs a simple
#  MinGW and 'make' install/set-up be useful - nothing else required!  
#
# 
## set compiler to use:
CC=gcc
# set default build to 32 bit
ARCH=32
# set default CFLAGS to CFLAGS_32 - will changed below as req!
CFLAGS=$(CFLAGS_$(ARCH))
# set default to add '.exe' to end of compiled output files  - will changed below as req!
EXE_END=.exe
#
## +++ BASELINE CFLAGS +++
# 32 bit default is with debugging, warnings, and use GNU C11 standard
CFLAGS_32=-g -Wall -m32 -pg -std=gnu11
# 64 bit default which provides same - but 64bit build instead
CFLAGS_64=-g -Wall -m64 -pg -std=gnu11
#
## +++ OPTIMSED CFLAGS +++
# normal optimisations for 32bit & 64bit - any computer
N-CFLAGS_32=-m32 -mfpmath=sse -flto -funroll-loops -Wall -std=gnu11
N-CFLAGS_64=-m64 -flto -funroll-loops -Wall -std=gnu11
# max optimisations for 32bit & 64bit - output only for computer compiled on!!
OPT-CFLAGS_32=-m32 -mfpmath=sse -flto -march=native -funroll-loops -Wall -std=gnu11
OPT-CFLAGS_64=-flto -funroll-loops -march=native -Wall -m64 -std=gnu11
#
#  Below is a list of CFLAGS and their purpose:
#
#	-g				: enables debugging
#	-Wall 			: enables compiler warnings
#	-m32 			: compiles a 32bit app
#	-m64 			: compiles a 64bit app
#	-mwindows 		: use the Windows libraries in the code + stop console 
#						window from opening too (ie for Windows code)
#	-mconsole 		: states is a Windows console application
#	-shared-libgcc 	: used by TDG-Minggw64 to used dll instead of static compile on Windows
#	-O2 or -Os		: specify optimisation level for the compiler or use: -mtune=native
#	-march=native  	: optimise for your computer only. To review run:
#						gcc -march=native -E -v - </dev/null 2>&1 | sed -n 's/.* -v - //p'
#	-mtune=native	: native optimisation level. To review run:
#						gcc -mtune=native -E -v - </dev/null 2>&1 | sed -n 's/.* -v - //p'
#	-mfpmath=sse   	: optimisation useful for 32bit compiles with floating point 
#						calcs are included (NB not req. in 64bit)
#	-funroll-loops 	: optimisation by enabling loop unrolling
#	-pg  			: Support application profiling. See: http://www.thegeekstuff.com/2012/08/gprof-tutorial/
#	-std=gnu99  	: use additional GNU with standard C99 features in code
#	-std=c99  		: use standard C99 features in code
# 	-std=c11  		: use standard C11 features in code for iso9899:2011 from 08 Dec 2011
# 	-std=gnu11  	: use additional GNU with standard C11 features in code for iso9899:2011
#	-ldl 			: includes the reference to the library that has the symbols for loading 
#						dynamic libraries (such as dlopen).
#	-Wpedantic  	: issues all the warnings demanded by strict ISO C
#   -Ofast          : includes -ffast-math option of GCC not supported by SQlite3
#  
#  NB: '-march' is specific to current computer. However '-mtune' includes optimisations 
#		for current computer, and runs on others too. Choose to suit for you own use case.
#
#
## +++ LIBRARY FLAGS: LIBFLAGS +++
LIBFLAGS=
#
#  Below is a list of LIBFLAGS and their purpose - add above as needed:
#
#	-DWIN32_LEAN_AND_MEAN		: Used to reduce the size of the header files excluding such ietms crypto, DDE, RPC, Windows Shell, and Winsock
#	-DWIN32_LEAN_AND_MEAN=1		: as above - using full when specifically: #include <winsock2.h>
#	-DUSE_MINGW_ANSI_STDIO=1	: Used by MinGW
#	-lsqlite3					: include sqlite3 library
#	-lpthread					: included pThreads library
#	-lcurses					: include original Curses library
#	-lncurses					: include new NCurses
#	-lpdcurses					: includes PDCurses (Windows) library
#	-lreadline					: inlcudes GNU Readline library support
#
#
## +++ SET DEFAULTS FOR WINDOWS ENVIRONMENT +++
RM = del
#
#
## +++ DETECT ENVIRONMENT +++
#
# command gets 'Linux' 'Darwin' 'FreeBSD' from command line on these OS environments
uname_S := $(shell sh -c 'uname -s 2>/dev/null || echo not')
#
# Use make to see if we are on Linux and if 64 or 32 bits?
ifeq ($(uname_S),Linux)
	# using Linux - change 'del' to 'rm'
	RM = rm
	# check if '64' or '32' bit architecture?
	ARCH := $(shell getconf LONG_BIT)
	# set new default $CFLAGS to match what we have discovered
	CFLAGS=$(CFLAGS_$(ARCH))
	# unset this as Linux does not need '.exe' tagged to outputs
	EXE_END=
	# using linux so change pdcurses to ncurses lib (if required!)
	# and add any other Linux specific library requirements too
	LIBFLAGS=-lpthread -ldl
endif
# Use make to see if we are on MAc OS X and if 64 or 32 bits?
ifeq ($(uname_S),Darwin)
	# using Mac OS X - change 'del' to 'rm'
	RM = rm
	# check if '64' or '32' bit architecture?
	ARCH := $(shell getconf LONG_BIT)
	# set new default $CFLAGS to match what we have discovered
	CFLAGS=$(CFLAGS_$(ARCH))
	# unset this as Mac OS X does not need '.exe' tagged to outputs
	EXE_END=
	# using macosx so change pdcurses to ncurses lib (if required)
	# and add any other Mac OS X specific library requirements too
	LIBFLAGS=
endif
#
#
## +++ DEFAULT MAKE OUTPUT +++ :
$(OUTNAME): $(SRC)
	$(CC) $(CFLAGS) -o $(OUTNAME)$(EXE_END) $(SRC) $(LIBFLAGS)

# if: 'make norm' use the other CFLAGS based on $ARCH defined
norm: $(SRC)
	$(CC) $(N-CFLAGS_$(ARCH)) -o $(OUTNAME)$(EXE_END) $(SRC) $(LIBFLAGS)

# if: 'make opt' use the other CFLAGS based on $ARCH defined
opt: $(SRC)
	$(CC) $(OPT-CFLAGS_$(ARCH)) -o $(OUTNAME)$(EXE_END) $(SRC) $(LIBFLAGS)

# if: 'make as' use to output an assembly source code version 
as: $(SRC)
	$(CC) -S -std=gnu11 -c $(SRC) $(LIBFLAGS)


clean:
	$(RM) $(OUTNAME)$(EXE_END)

val:
	 $(shell sh -c 'valgrind --leak-check=full --show-leak-kinds=all ./$(OUTNAME)$(EXE_END)')

# used by Emacs for 'flymake'
check-syntax:
	gcc -o nul -S ${CHK_SOURCES}
