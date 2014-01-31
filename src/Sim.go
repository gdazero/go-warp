/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Sim

/*
 * this is the backbone of each Logical Process, here are defined and implemented 
 * all the Time Warp functions needed by the LPs
 */

import(
    "os"
    "fmt"
    "./Communication"
    "./DT"
    "./Local"
    "./Const"
    "./Gvt"
    "./Shared"
)


const TOOFAR = 25     // limited optimism synchronization: sets how far from the GVT a LP can go


func Setup(lpn int, simt DT.Time, f func(ev *DT.Event, l *Local.LocalData)) {
    Communication.AllocateChans(lpn)
    Gvt.Setup(lpn)
    Shared.Setup(lpn, simt, f)
}


/*
 * every LP must perform an initialize() operation, that creates all
 * the needed structures and variables
 */
func Initialize(i DT.Pid) *Local.LocalData{
    var data *Local.LocalData

    data = Local.Initialize(i)
    Shared.State[i] = Const.LPRUNNING

    return data
}


func Simulate(data *Local.LocalData) {

    for {

        if Shared.State[data.IndexLP] == Const.LPSTOPPED {
            Shared.State[data.IndexLP] = Const.LPSTOPPED
            
            return
        }

        receiveAll(data)

        if data.SimTime>=Shared.EndTime {
            goIdle(data)
        }

        manageEvent(data)

        if data.GvtFlag && !Gvt.CheckEvaluation() {
            t := Gvt.GetGvt() 
            if t != Const.ERR {
                setGvt(t,data)
            }
        }

    }
}


/*
 * creates and sends a message to the receiver that contains the event to be
 * noticed. Saves the related anti-message in sender local area
 */
func NoticeEvent(ev *DT.Event, receiver DT.Pid, data *Local.LocalData) {
    var tm DT.TimedMessage
    var msg *DT.Message

    /* creating the message to send */
    msg = DT.CreateMessage(data.IndexLP, receiver, *ev)

    if receiver == data.IndexLP {
        if !data.FutureEvents.Insert(ev) {
            fmt.Println("GO-WARP, ERROR: THE HEAP IS FULL -",len(data.FutureEvents))
            fmt.Println(data.IndexLP,"- GO-WARP, ERROR: EVENT NOT INSERTED!")
            os.Exit(1)
        }
    } else {
        /* sending the message */
        sendMessage(msg,data)
    }

    tm = DT.TimedMessage{*msg,data.SimTime}

    size := DT.Insert(tm, data.MsgSent)
    if size > Const.TOOLARGE && Shared.State[data.IndexLP] != Const.LPEVALGVT {
        ask4NewGvt(data)
    }
}


func receiveAll(data *Local.LocalData) {
    Loop: for {
        msg := Communication.Receive(data.IndexLP)

        if msg == nil {
            break Loop
        }

        manageMessage(data, msg)
    }
}


func manageMessage(data *Local.LocalData, msg *DT.Message) {
    switch msg.Ev.Time {
        case Const.GVTEVAL:
        if Shared.State[data.IndexLP] != Const.LPSTOPPED {
            evaluateLocalMin(data)
        }

        case Const.ABORTMSG:
        Shared.State[data.IndexLP] = Const.LPSTOPPED

        case Const.ACK:
        gotAck(msg,data)

        default:
        sendAck(msg,data)

        if checkAntimsg(&msg.Ev, data) {
            return
        }
        if msg.Ev.Type.Flag == Const.ANTIMSG {		// anti-message
            annihilate(&(msg.Ev), data)
            return
        }
        if msg.Ev.Time < data.SimTime {			// straggler message
            rollback(msg.Ev.Time, data)
        }

        /* finally we can insert the message in the heap */
        if !(data.FutureEvents).Insert(&msg.Ev) {
            fmt.Println(data.IndexLP,"- GO-WARP, ERROR: EVENT NOT INSERTED!")
            os.Exit(1)
        }
    }
}


/* 
 * returns false only if it has failed managing an event (because the heap is empty)
 */
func manageEvent(data *Local.LocalData) bool {
    var ev *DT.Event

    data.N_PROCESSED++

    t := data.FutureEvents.GetMinTime()

    if t >= Shared.EndTime {
        goIdle(data)
        return false
    } else if t > data.SimTime {
        data.SimTime = t
    } else if t == data.SimTime {
        /* OK, DN */
    } else if t == Const.NOTIME {
        /* heap empty */
    } else {
        fmt.Println(data.IndexLP, "- GO-WARP, ERROR: PROCESSING AN EVENT IN THE PAST!")
        os.Exit(1)
    }

    ev = data.FutureEvents.ExtractHead()
    if ev == nil { 
        return false
    }

    Shared.EventManager(ev, data)

    size := DT.Insert(*ev,data.ProcessedEvents)
    if size > Const.TOOLARGE && Shared.State[data.IndexLP] != Const.LPEVALGVT {
        ask4NewGvt(data)
    }
    
    return true
}


func rollback(t DT.Time, data *Local.LocalData){
    data.SimTime = t

    el := data.ProcessedEvents.Back()
    Loop: for el != nil {
        e := el.Value.(DT.Event)
        el = el.Prev()
        if e.Time >= data.SimTime {
            if !data.FutureEvents.Insert(&e) {
                fmt.Println("GO-WARP, ERROR: INSERTING A PROCESSED EVENT!")
            }
            data.N_PROCESSED--
        } else {
            break Loop
        }
    }

    el = data.MsgSent.Back()
    Loop1: for el != nil {
        mp := el.Value.(DT.TimedMessage)
        if mp.GetTime() >= data.SimTime {
            el = el.Prev()
            anti := createAntiMessage(&mp.M)

            if mp.M.Receiver == data.IndexLP {
                annihilate(&(anti.Ev), data)
            } else {
                sendMessage(anti,data)
            }
        } else {
            break Loop1
        }
    }

    DT.DeleteAfter(data.SimTime, data.ProcessedEvents)
    DT.DeleteAfter(data.SimTime, data.MsgSent)

    Shared.N_rollback[data.IndexLP]++

}


func sendMessage(msg *DT.Message, data *Local.LocalData) {
    tm := DT.TimedMessage{*msg,data.SimTime}
    size := DT.Insert(tm, data.OutgoingMsg)

    if size > Const.TOOLARGE {
        if Shared.State[data.IndexLP] != Const.LPEVALGVT {
            ask4NewGvt(data)
        }
    }
    Communication.Send(msg)
}


func goIdle(data *Local.LocalData) {
    if Shared.State[data.IndexLP] == Const.LPSTOPPED { return }

    Shared.State[data.IndexLP] = Const.LPIDLE
    term := checkAllIdle()
    if term { 
        killall(data)
        Shared.State[data.IndexLP] = Const.LPSTOPPED
    } else {
        m := Communication.BlockingReceive(data.IndexLP)	// the process blocks indefinitively
        
        manageMessage(data,m)

        if Shared.State[data.IndexLP] != Const.LPSTOPPED {
            Shared.State[data.IndexLP] = Const.LPRUNNING
        }
    }
}


func createAntiMessage(msg *DT.Message) *DT.Message{
    var e DT.Event
    var m DT.Message

    e = *DT.CreateEvent(-msg.Ev.Id,msg.Ev.Time,DT.Info{0,0,Const.ANTIMSG})
    m = *DT.CreateMessage(msg.Sender, msg.Receiver, e)

    return &m
}


func annihilate(antimsg *DT.Event, data *Local.LocalData) {
    done := false
    done = DT.IsPresent(antimsg.Time,data.ProcessedEvents)
 
    if done && antimsg.Time <= data.SimTime {
        rollback(antimsg.Time,data)
    }

    ev := DT.CreateEvent(-antimsg.Id,0,DT.Info{0,0,0})
    del := data.FutureEvents.DeleteExternId(ev)

    if del.Id == Const.ERR && del.Time == Const.ERR {
        antimsg.Time = data.SimTime	// timestamping the anti-message it will be possible to rollback its reception
            
        DT.Insert(*antimsg, data.AntiMsg2Annihilate)
    }
}


func checkAntimsg(ev *DT.Event, data *Local.LocalData) bool {
    m := DT.CreateMessage(0,0,*ev)
    anti := createAntiMessage(m)
    ret := DT.Delete(anti.Ev, data.AntiMsg2Annihilate)
    return ret
}


func ask4NewGvt(data *Local.LocalData) {
    if Shared.State[data.IndexLP] == Const.LPSTOPPED { return }
    if Gvt.CheckEvaluation() { return }
    Gvt.StartEvaluation(Shared.Lpnum)

    ev := DT.CreateEvent(0,Const.GVTEVAL,DT.Info{0,0,0})
    for i:=0;i<Shared.Lpnum;i++ {
        if Shared.State[i] != Const.LPSTOPPED && data.IndexLP != DT.Pid(i) {
            msg := DT.CreateMessage(data.IndexLP,DT.Pid(i),*ev)
            Communication.Send(msg)
        }
    }
    evaluateLocalMin(data)
}


func evaluateLocalMin(data *Local.LocalData) {
    var mintime DT.Time = 1000000

    Shared.State[data.IndexLP] = Const.LPEVALGVT
 
    /* mintime computation and communication */
     minheap := data.FutureEvents.GetMinTime()
    minout := DT.GetMinTime(data.OutgoingMsg)
    minack := DT.GetMinTime(data.Acked)
    
    if minheap<mintime && minheap!=Const.NOTIME {
        mintime = minheap
    }
    if minout<mintime && minout!=Const.NOTIME {
        mintime = minout
    }
    if minack<mintime && minack!=Const.NOTIME {
        mintime = minack
    }

    Gvt.SetLocalMin(mintime,data.IndexLP)
    data.GvtFlag = true		// local min has been set
    
    data.Acked.Init()
}


func setGvt(gvt DT.Time, data *Local.LocalData){

    if Shared.State[data.IndexLP] == Const.LPSTOPPED { return }

    if gvt < data.Gvt {
        fmt.Println(data.IndexLP,", GO-WARP, ERROR: THE NEW GVT VALUE IS LOWER THAN THE PREVIOUS ONE!")
        os.Exit(1)
    }
    data.GvtFlag = false
    data.Gvt = gvt

    fossilCollection(gvt, data)
}


func fossilCollection(t DT.Time, data *Local.LocalData) {
    DT.DeleteBefore(t, data.ProcessedEvents)
    DT.DeleteBefore(t, data.MsgSent)
    Shared.State[data.IndexLP] = Const.LPRUNNING

    data.Acked.Init()
}


func sendAck(msg *DT.Message, data *Local.LocalData) {
    var e *DT.Event

    if data.GvtFlag {
        e = DT.CreateEvent(msg.Ev.Id,Const.ACK,DT.Info{0,0,Const.YOURS})
    } else {
        e = DT.CreateEvent(msg.Ev.Id,Const.ACK,DT.Info{0,0,Const.MINE})
    }
    ack := DT.CreateMessage(msg.Receiver, msg.Sender, *e)
    Communication.Send(ack)
}


func gotAck(msg *DT.Message,data *Local.LocalData) {
    var found bool = false

    el := data.OutgoingMsg.Front()
    Loop: for el != nil {
        m := el.Value.(DT.TimedMessage)
        if m.M.Receiver == msg.Sender && m.M.Ev.Id == msg.Ev.Id && m.M.Ev.Time != Const.ACK { 
            if msg.Ev.Type.Flag == Const.MINE {
                data.OutgoingMsg.Remove(el)
            } else if msg.Ev.Type.Flag == Const.YOURS {
                DT.Insert(m, data.Acked)
                data.OutgoingMsg.Remove(el)
            } else {
                fmt.Println(data.IndexLP,"- GO-WARP, ERROR: UNKNOWN FLAG TYPE")
                os.Exit(1)
            }
            found = true
            break Loop
        }
        el = el.Next()
    }
    if !found { fmt.Println("GO-WARP, ERROR: WHAT ABOUT THIS ACK?") }

}


func checkAllIdle() bool {
    var ret bool = true
    Loop: for i:=0;i<Shared.Lpnum;i++ {
        if Shared.State[i] != Const.LPIDLE {
            ret = false
            break Loop
        }
    }
    return ret
}


func killall(data *Local.LocalData) {
    ev := DT.CreateEvent(0,Const.ABORTMSG,DT.Info{0,0,0})

    for i:=0;i<Shared.Lpnum;i++ {
            m := DT.CreateMessage(data.IndexLP,DT.Pid(i),*ev)
            Communication.Send(m)
    }
}
