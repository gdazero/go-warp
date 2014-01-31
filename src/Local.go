/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Local

/*
 * This package contains all the default variables and structures,
 * each LP must import this package and initialize these variables 
 * 
 */

import(
    "./DT"
    "./Heap"
    "fmt"
    "os"
    list "container/list"
)

type LocalData struct {
    N_PROCESSED int
    SimTime DT.Time
    Gvt DT.Time
    IndexLP DT.Pid
    FutureEvents Heap.EventHeap
    ProcessedEvents *list.List
    MsgSent *list.List
    AntiMsg2Annihilate *list.List
    OutgoingMsg *list.List
    Acked *list.List
    Pending bool
    GvtFlag bool
}


/*
 * every LP must perform an initialize() operation, that creates
 * the needed structures and variables
 */
func Initialize(i DT.Pid) *LocalData {
    var d LocalData = *new(LocalData)

    /* initialize all LP variables */
    d.N_PROCESSED = 0
    d.IndexLP = i
    d.SimTime = 0
    d.Gvt = 0
    d.Pending = true
    d.GvtFlag = false
    d.FutureEvents = Heap.InitializeHeap()
    d.ProcessedEvents = DT.NewList()
    d.MsgSent = DT.NewList()
    d.AntiMsg2Annihilate = DT.NewList()
    d.OutgoingMsg = DT.NewList()
    d.Acked = DT.NewList()

    return &d
}


func (l *LocalData) NewEvent(ev *DT.Event) {
    if !l.FutureEvents.Insert(ev) {
        fmt.Println("GO-WARP, ERROR: EVENT NOT INSERTED") 
        l.FutureEvents.Print()
        os.Exit(1)
    }
}
