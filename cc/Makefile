C = g++

.DEFAULT_GOAL=all

all: strip_packing

strip_packing: *.cc *.h
	$C -O2 *.cc -o strip_packing `pkg-config --cflags --libs gtk+-2.0`