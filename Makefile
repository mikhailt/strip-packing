include $(GOROOT)/src/Make.inc

TARG=strip_packing
GOFILES=\
	main.go\
	kp_algo.go

include $(GOROOT)/src/Make.cmd