/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package State

import(
    "./DT"
)

type State struct {
    SimTime DT.Time
    LpVar DT.LPstate
}


func CreateState(time DT.Time, lpvar DT.LPstate) *State {
    var state *State = new(State)

    *state = State{time,lpvar}
    return state
}


/* type State implements the DT.Element interface */
func (st State) GetTime() DT.Time {
    return st.SimTime
}


func (st State) IsEqual(e DT.Elem) bool {
    ret := false
    st1 := e.(State)
    if st.SimTime==st1.SimTime {
        ret = true
    }
    return ret
}
