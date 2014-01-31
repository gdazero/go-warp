/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Const

const(
/* as sender / receiver of a message */
    SERVERID = -1	// server identifier

/* in the field type of an event */
    ANTIMSG = -4 	// the message is an anti-message

/* in the field time of an event */
    ACK = -5		// the message is an acknowledgement
    GVTEVAL = -6
    ABORTMSG = -7
    RBMSG = -8

/* different types of ack messages */
    MINE = -1		// the sender of the ACK assumes the responsibility for the message in the GVT evaluation
    YOURS = -2		// the sender of the ACK doesn't assume the responsibility for the message until the successive GVT evaluation
    
/* significant constants */
    FEWFREEPLACES = 4000 	// the free space in an array is too low
    LISTLEN = 5000		// the max length of a queue
    HEAPSIZE = 500     		// heap size
    TOOLARGE = 500

/* possible message colors */
    WHITE = 1
    BLACK = 2
    NOTACOLOR = -1

/* logs creation constants */
    LOGDIR = "logs/"
    PERM = 0666
    GVTLOG = -3

/* server managment constants */
    NOTIME = -2
    ERR = -1
)

/* possible process states */
const(
    LPNOTSTART = iota
    LPRUNNING = iota
    LPIDLE = iota
    LPSTOPPED = iota
    LPEVALGVT = iota
)
