/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Gvt

import(
    "./Const"
    "./DT"
    "./Shared"
    "sync"
)

var(
    lpNum int
    localMin []DT.Time
    gvt DT.Time
    gvtFlag bool
    lock sync.Mutex
)

const MAXTIME = 1 << 31 -1
const EMPTY = -13


func Setup(lpnum int){

    localMin = make([]DT.Time, lpnum)
    lpNum = lpnum
    gvt = 0
    gvtFlag = false
}


func StartEvaluation(lpnum int) {
    if !gvtFlag {
        for i:=0;i<lpnum;i++ {
            localMin[i] = EMPTY
        }
        lpNum = lpnum
        gvtFlag = true
    }
}


/* if true the a GVT calculation is running */
func CheckEvaluation() bool {
    return gvtFlag
}


func SetLocalMin(time DT.Time,pid DT.Pid) {

    if !gvtFlag {
        return
    }

    localMin[pid] = time

    Loop: for i:=0;i<len(localMin);i++ {
        if localMin[i] == EMPTY { return }
    }

    lock.Lock()
    setGVT()
    lock.Unlock()
}


func setGVT() {

    Loop: for i:=0;i<len(localMin);i++ {
        if localMin[i] == EMPTY { return }
    }

    tmpMin := DT.Time(MAXTIME)
    for i:=0;i<len(localMin);i++ {
        if localMin[i] < tmpMin && localMin[i] != Const.NOTIME {
            tmpMin = localMin[i]
        }
    }
    gvt = tmpMin

    for i:=0;i<len(localMin);i++ {
        localMin[i] = EMPTY
    }
    gvtFlag = false
    Shared.N_gvt++
}


func GetGvt() DT.Time {

    if gvtFlag {
        return Const.ERR
    }
    return gvt
}
