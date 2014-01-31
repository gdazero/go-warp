/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Heap

/*
 * HEAP MANAGEMENT
 */

import (
    "strconv"
    "fmt"
    "os"
    "./Const"
    "./DT"
)

const EVARRSIZE = 2000

type Node struct {
    time DT.Time;
    events *[]DT.Event
}

type EventHeap []Node

var NPRINT int = 0

/*
 * Position 0 of an eventheap contains some information on the events
 * contained in other positions. 
 * 
 * If heap[0].time == -1 then the heap is empty
 */


/* event heap initialization */
func InitializeHeap() EventHeap{
    var heap = make([]Node, 1, Const.HEAPSIZE)
    heap[0] = Node{-1, nil}	// the heap is empty

    return heap
}


/* returns true if the heap has no events */
func (heap *EventHeap) IsEmpty() bool{
    var ret bool = false
    
    if (*heap)[0].time == -1 { 
        ret = true 
    }

    return ret
}


/* Insert an event in the heap, that remains balanced */
func (heap *EventHeap) Insert(evptr *DT.Event) bool {
    var ret bool = true
    var nodepos, pos, fpos int
    var evArr *[]DT.Event
    var nod, father Node

    length := len(*heap)
    capacity := cap(*heap)
    pos = length

    nodepos = heap.isPresent((*evptr).Time)

    if nodepos > 0 {				// another event with timestamp equal to *evptr is already present
        evArr = (*heap)[nodepos].events

        if len(*evArr) >= cap(*evArr) {
            ret=false
        } else {
            (*evArr) = (*evArr)[0:len(*evArr)+1]
            (*evArr)[len(*evArr)-1] = *evptr	// the event is placed in the last position
        }
    } else if length >= capacity { 
        ret=false             
    } else {					// an event with timestamp equal to *evptr is not present

        arr := make([]DT.Event, 1, EVARRSIZE)
        arr[0] = *evptr
        nod = Node{evptr.Time, &arr}

        (*heap) = (*heap)[0:length+1]

        (*heap)[pos] = nod

        Loop: for pos>1 {

            fpos = pos/2
            father = (*heap)[fpos]

            if father.time > nod.time {
                (*heap)[fpos] = nod
                (*heap)[pos] = father
                pos = fpos 
            } else { break Loop }
    
        }

        (*heap)[0].time = 1
    }

    return ret
}


/* returns the minimum time in the heap */
func (heap *EventHeap) GetMinTime() DT.Time {
    if heap.IsEmpty() { return Const.NOTIME }
    return (* (*heap)[1].events )[0].Time
}


/* extracts the first event in the heap, that remains balanced */
func (heap *EventHeap) ExtractHead() *DT.Event {
    var head DT.Event
    if heap.IsEmpty() { return nil }

    head = (* (*heap)[1].events)[len(*(*heap)[1].events)-1]
    if !heap.Delete(&head) {
        fmt.Println("GO-WARP: extracthead")
        os.Exit(1)
        return nil
    }
    return &head
}


/* deletes an event from the heap, that remains balanced */
func (heap *EventHeap) Delete(evptr *DT.Event) bool{
    var nodepos, pos, minpos int
    var tmp, min Node
    var evArr *[]DT.Event

    nodepos = -1
    Loop1: for i:=1; i<len(*heap); i++ {
        if (*heap)[i].time == evptr.Time {
            nodepos = i
            break Loop1
        }
    }
    if nodepos == -1 {
        return false
    }

    if len(*(*heap)[nodepos].events) > 1 {
        id := evptr.Id
        evArr =(*heap)[nodepos].events
        pos = -1

        Loop2: for i:=0; i<len(*(*heap)[nodepos].events); i++ {
            if (*evArr)[i].Id == id {
                pos = i
                break Loop2
            }
        }
        if pos == -1 {
            return false
        } else {
            for i:=pos+1; i<len(*(*heap)[nodepos].events); i++ {
                (*evArr)[i-1] = (*evArr)[i]
            }
            *evArr = (*evArr)[0:len(*evArr)-1]
        }

    } else { 	// the element to be eliminated is the only one with that timestamp
	  
        (*heap)[nodepos] = (*heap)[len(*heap)-1]

        (*heap) = (*heap)[0:len(*heap)-1]

        if len(*heap) == 1 {
            (*heap)[0].time = -1	// empty heap
        }
        if nodepos == len(*heap) {
            return true
        }

        if (*heap)[nodepos].time < (*heap)[nodepos/2].time {
            Loop3: for nodepos>1 {
                fpos := nodepos/2
                father := (*heap)[fpos]
                if father.time > (*heap)[nodepos].time {
                    (*heap)[fpos] = (*heap)[nodepos]
                    (*heap)[nodepos] = father
                    nodepos = fpos
                } else { break Loop3 }
            }

        } else {
            Loop4: for {
                desc := heap.descendants(nodepos)

                if desc == 0 {
                    break Loop4
                } else if desc == 1 {
                    if (*heap)[nodepos].time > (*heap)[nodepos*2].time {
                        tmp = (*heap)[nodepos]
                        (*heap)[nodepos] = (*heap)[nodepos*2]
                        (*heap)[nodepos*2] = tmp
                        nodepos = nodepos*2
                    } else { 
                        break Loop4
                    }
                } else {

                    if (*heap)[nodepos*2].time < (*heap)[nodepos*2+1].time {
                        min = (*heap)[nodepos*2]
                        minpos = nodepos*2
                    } else {
                        min = (*heap)[nodepos*2+1]
                        minpos = nodepos*2+1
                    }

                    if min.time > (*heap)[nodepos].time {
                        break Loop3
                    } else {
                        (*heap)[minpos] = (*heap)[nodepos]
                        (*heap)[nodepos] = min
                        nodepos = minpos
                    }
                }
            }
        }
    }

    return true
}


/* heap debug */
func (heap *EventHeap) PrintNodes() {
    if heap.IsEmpty() {
        fmt.Println("GO-WARP, empty heap")
    } else {

        fmt.Println("GO-WARP, nodes in the heap:")

        for i:=1; i<len(*heap); i++ {
            fmt.Printf("%d,   ", (*heap)[i].time)
        }

        fmt.Println()
    }
}


/* heap debug */
func (heap *EventHeap) Print() {
    if heap.IsEmpty() {
        fmt.Println("GO-WARP, empty heap")
    } else {

        fmt.Println("GO-WARP, nodes in the heap:")

        for i:=1; i<len(*heap); i++ {
            evArr := (*heap)[i].events

            fmt.Printf("Events with time %d: ",(*heap)[i].time)
            for j:=0; j<len(*(*heap)[i].events); j++ {
                fmt.Printf("%d,  ", (*evArr)[j].Id)
            }
            fmt.Println()
        }
        fmt.Println()
    }
}


func (heap *EventHeap) descendants(pos int) int{
    var num int
    if len(*heap) <= 2*pos {
        num = 0
    } else if len(*heap) == 2*pos+1 {
        num = 1
    } else { 
        num = 2
    }
    return num
}


/*
 *  searches for an elememnt with a given time, if found then returns the 
 *  position otherwise returns -1 
 */
func (heap *EventHeap) isPresent(t DT.Time) int {
    var ret int = -1
    for i:=1; i<len(*heap); i++ {
        if t == (*heap)[i].time {
            ret = i
            break;
        }
    }

    return ret
}


/* 
 * searches, deletes and returns an event using its identifier
 */
func (heap *EventHeap) DeleteExternId(ev *DT.Event) DT.Event {
    var ret DT.Event = DT.Event{Const.ERR,Const.ERR,DT.Info{0,0,0}}

    Loop: for i:=1;i<len(*heap);i++ {
        for j:=0;j<len(*(*heap)[i].events);j++ {
            if (*(*heap)[i].events)[j].Id == ev.Id {
                ret = (*(*heap)[i].events)[j]
                if heap.Delete(&(*(*heap)[i].events)[j]) {
                    break Loop
                }
            }
        }
    }
    return ret
}


func (heap *EventHeap) GetCopy() EventHeap {
    var ret EventHeap = InitializeHeap()

    l := len(*heap)
    ret = ret[0:l]
    for i:=0; i<l; i++ {
        ret[i] = (*heap)[i]
        if (*heap)[i].events != nil {
            lea := len(*(*heap)[i].events)
            ea := make([]DT.Event,EVARRSIZE)
            ea = ea[0:lea]
            for j:=0;j<lea;j++ {
                ea[j] = (*(*heap)[i].events)[j]
            }
            ret[i].events = &ea
        }
    }
    return ret
}


func (heap *EventHeap) GetString() string {
    var s string = ""
      
    if (*heap)[0].time == -1 {
        s = "The heap is emtpy"
    } else {

        s += "Heap nodes:\n"

        Loop: for i:=1; i<len(*heap); i++ {
            s += "Events "+ strconv.Itoa(int((*heap)[i].time))+":"

            evArr := (*heap)[i].events
            if evArr == nil {
                s += "(*heap)[i].events = NIL !!!\n"
                continue Loop
            }
            for j:=0; j<len(*(*heap)[i].events); j++ {
                s += strconv.Itoa(int((*evArr)[j].Id)) + "   "
            }
            s += "\n"
        }
        s += "\n"
    }
    return s
}


func (heap *EventHeap) DeleteMatching(f func(ev *DT.Event, a ...interface{}) bool, t ...interface{}) {
    var matching []*DT.Event = make([]*DT.Event, Const.HEAPSIZE)
    var index int = 0

    for i:=1;i<len(*heap);i++ {
        for j:=0;j<len(*(*heap)[i].events);j++ {
            if f(&(*(*heap)[i].events)[j],t[0],t[1]) {
                matching[index] = &(*(*heap)[i].events)[j]
                index++
            }
        }
    }
    Loop: for i:=0;i<index;i++ {
        if matching[i] == nil { break Loop }
        heap.Delete(matching[i])
    }
}
