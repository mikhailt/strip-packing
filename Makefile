include $(GOROOT)/src/Make.inc

TARG=strip_packing
GOFILES=\
	main.go\
	kp_algo.go\
	validate.go\
	2d_algo.go\
	kp_msp_algo.go

include $(GOROOT)/src/Make.cmd