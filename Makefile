#.............................................................................
#
#   Makefile Template.
#       Originally Created: June 2013 by Simon Rowe <simon@wiremoons.com>
#.............................................................................
#
#	Makefile for amt
#
#   gcc -Wall -std=gnu11 -m64 -g -o amt amt.c cli-args.c main.c sqlite3.c -Lpthread -ldl
#
## CHANGE BELOW TO SOURCE FILE NAMES & OUTPUT FILE-NAME
SRC=amt-db-funcs.c cli-args.c main.c sqlite3.c
OUTNAME=amt
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
# add any LIBFLAGS
LIBFLAGS=-lpthread -ldl
#
## +++ SET DEFAULTS FOR WINDOWS ENVIRONMENT +++
RM = del
## +++ DETECT ENVIRONMENT +++
#
# find 'Linux' 'Darwin' 'FreeBSD' from command line on these OS environments
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

# clean up riles
clean:
	$(RM) $(OUTNAME)$(EXE_END)

val:
	 $(shell sh -c 'valgrind --leak-check=full --show-leak-kinds=all ./$(OUTNAME)$(EXE_END)')

# used by Emacs for 'flymake'
check-syntax:
	gcc -o nul -S ${CHK_SOURCES}
