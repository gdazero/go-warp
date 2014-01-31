/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package DT /* Data Types */

import(
    "./Const"
    "fmt"
    list "container/list"
)

const EMPTYPLACE = 2

type Pid int16
type Time int32
type Info struct {    //the simulation modeler can use this field for data transfer
    From int
    To int
    Flag int32
}

type Message struct{
    Sender Pid
    Receiver Pid
    Ev Event
}
type TimedMessage struct{
    M Message
    T Time
}

type Event struct { 
    Id int32   // high-order 16 bits = LP info, low-order 16 bits = event info
    Time Time
    Type Info
}

/* interface useful as Elem of a List */
type Elem interface{
    GetTime() Time
    IsEqual(e Elem) bool
}


/* list.List management functions */
func NewList() *list.List {
    return list.New()
}


func Insert(e Elem, L *list.List) int {
    if L.Len() == 0 {
        L.PushFront(e)
        return L.Len()
    }

    front := L.Front()
    if e.GetTime() < front.Value.(Elem).GetTime() {
        L.InsertBefore(e,front)
        return L.Len()
    }

    el := L.Back()
    Loop: for {
        if el.Value.(Elem).GetTime() > e.GetTime() {
            el = el.Prev()
        }else {
            break Loop
        }
    }
    L.InsertAfter(e, el)

    return L.Len()
}


/* delete all the elements with time <= t */
func DeleteBefore(t Time, L *list.List) {
    if L.Len() == 0 {
        return
    }
    back := L.Back()
    if back.Value.(Elem).GetTime() <= t {
        L = L.Init()
        return
    }
    Loop: for {
        el := L.Front()
        if el.Value.(Elem).GetTime() <= t {
            L.Remove(el)
        } else {
            break Loop
        }
    }
}


/* delete all the elements with time >= t */
func DeleteAfter(t Time, L *list.List) {
    if L.Len() == 0 {
        return
    }
    front := L.Front()
    if front.Value.(Elem).GetTime() >= t {
        L = L.Init()
        return
    }
    Loop: for {
        el := L.Back()
        if el.Value.(Elem).GetTime() >= t {
            L.Remove(el)
        } else {
            break Loop
        }
    }
}


func Delete(e Elem, L *list.List) bool {
    ret := false

    if L.Len() == 0 {
        return ret
    }
    back := L.Back()
    if e.GetTime() > back.Value.(Elem).GetTime() {
        return ret
    }

    el := L.Front()
    Loop: for i:=0;el!=nil;i++ {
        elt := el.Value.(Elem).GetTime()
        if elt > e.GetTime() {
            break Loop
        } else if e.IsEqual(el.Value.(Elem)) {
            L.Remove(el)
            ret = true
            break Loop
        }
        el = el.Next()
    }
    return ret
}


func GetMinTime(L *list.List) Time {
    if L.Len() == 0 {
        return Const.NOTIME
    }
    front := L.Front()
    return front.Value.(Elem).GetTime()
}


func Print(L *list.List) {
    fmt.Println("GO-WARP, list content...")
    el := L.Front()
    Loop: for el!=nil {
        fmt.Print(el.Value)
        el = el.Next()
    }
    fmt.Println()
}


/* Checks if an Element with time t is in the list */
func IsPresent(t Time, L *list.List) bool {
    ret := false
    el := L.Front()
    Loop: for el!=nil{
        elt := el.Value.(Elem).GetTime()
        if elt < t {
            el = el.Next()
        } else if elt == t {
            ret = true
            break Loop
        } else {
            break Loop
        }
    }
    return ret
}


func GetNearest(t Time, L *list.List) Elem {
    el := L.Front()
    past := el
    Loop: for el!=nil {
        if el.Value.(Elem).GetTime() <= t {
            past = el
            el = el.Next()
        } else {
            break Loop
        }
    }
    return past.Value.(Elem)
}


func CreateEvent(id int32, t Time, info Info) *Event {
    var ev *Event = new(Event)
    *ev = Event{id, t, info}
    return ev
}


func CreateMessage(sendr Pid, recvr Pid, e Event) *Message {
    var msgptr *Message = new(Message)

    *msgptr = Message{sendr, recvr, e}
    return msgptr
}


/* type Event implements Elem interface */
func (ev Event) GetTime() Time {
    return ev.Time
}


func (ev Event) IsEqual(e Elem) bool {
    ret := false
    ev1 := e.(Event)
    if ev.Id==ev1.Id && ev.Time==ev1.Time {
        ret = true
    }
    return ret
}


/* type TimedMessage implements Elem interface */
func (tm TimedMessage) GetTime() Time {
    return tm.T
}


func (tm TimedMessage) IsEqual(e Elem) bool {
    ret := false
    tm1 := e.(TimedMessage)
    if tm.M.Receiver==tm1.M.Receiver && tm.M.Sender==tm1.M.Sender && tm.M.Ev.Id==tm1.M.Ev.Id {
        ret = true
    }
    return ret
}
