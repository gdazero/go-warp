/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Communication

import(
    "./DT"
    "fmt"
)

const MAXBUFFER = 10000
const MAXALLOCN = 1
var(
    Chanptr *[]chan DT.Message
    lock chan int = make(chan int)
    allocations int = 0
)


func AllocateChans(nChan int) {

    if allocations >= MAXALLOCN { return }

    ch := make ([]chan DT.Message, nChan)	// this is to make the array

    for i:=0; i<nChan; i++ {
        ch[i]=make(chan DT.Message, MAXBUFFER)	// this is to make the chans
    }
    Chanptr = &ch
    allocations++
}


/* Send a message to destination */
func Send(msg *DT.Message) {
    (*Chanptr)[msg.Receiver] <- *msg
}


func Receive(recvid DT.Pid) *DT.Message {
    var msg DT.Message
    var ret *DT.Message
    var ok bool = false

    msg,ok = <- (*Chanptr)[recvid]
    if ok {
        ret = &msg
    } else {
        ret = nil
    }

    return ret
}

/* blocking receive */
func BlockingReceive(recvid DT.Pid) *DT.Message {
    select {
        case msg := <- (*Chanptr)[recvid]:
            return &msg
    }
    fmt.Println("GO-WARP: receive error!")
    return nil
}


func Sync() {
    <- lock
}

func Unlock() {
    lock <- 0
}
