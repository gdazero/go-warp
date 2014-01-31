/*
	Go-based implementation of the PHOLD synthetic benchmark
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package main

import( 
    "../src/DT"
    "../src/Sim"
    "../src/Shared"
    "../src/Random"
    "../src/Local"
    "fmt"
    "flag"
    "os"
    "bufio"
    "strings"
    "strconv"
    "time"
    "runtime"
    "sync"
)

const(
    usage="Main.out #LPs [if 0 -> autoconf] #ENTITIES"
    conf="./PHOLD/phold.conf"
    cpufile="/proc/cpuinfo"
    cpustr="processor"
)

var(
    lpnum int
    entitynum int
    density float
    n_events int
    endtime DT.Time
    nFPops int
    randGen *Random.RNG

    initEv []DT.Event

    idcount int32

    startT int64
    endT int64
    simterm bool
    n_term int
    print sync.Mutex

   n_cores int
)


func main() {
    n_lp,n_ent := readParams()

    initPhold(n_lp, n_ent)

    fmt.Println("GO-WARP: the simulator will use",runtime.GOMAXPROCS(-1),"COREs")
    fmt.Println("GO-WARP: the simulation will use",n_lp,"LPs")

    startT = time.Nanoseconds()
    for i:=1;i<n_lp;i++ {
        go launchLP(DT.Pid(i), n_ent/n_lp)
    }
    launchLP(0,n_ent/n_lp)

    for (n_term != lpnum) {
        time.Sleep(1e3)
    }
    printStats(endT)
}


func readParams() (nlp int, nent int) {
    flag.Parse()
    narg := flag.NArg()
    if narg != 2 {
        fmt.Printf("%s\n",usage)
        os.Exit(1)
    }
    nlp,_ = strconv.Atoi(flag.Arg(0))
    nent,_ = strconv.Atoi(flag.Arg(1))

    if nlp == 0 {
        nlp = getNCpu()
        fmt.Println("GO-WARP: value 0 as LP number, it means that the simulator will use all the available CPU cores")
    }

    fmt.Println("PARAMS:",nlp,nent,)
    return nlp,nent
}


func getNCpu() int {
    var line string
    var rdErr os.Error
    var count int = 0

    file,err := os.Open(cpufile, os.O_RDONLY,0)
    if err != nil {
        fmt.Println("GO-WARP, error opening:",err)
        os.Exit(1)
    }
    rd := bufio.NewReader(file)

    for i:=0;rdErr == nil;i++ {
        line, rdErr = rd.ReadString('\n')
        if len(line) > 0 {
            count+=strings.Count(line,cpustr)
        }
    }

    return count
}


func initPhold(nlp int, nent int) {
    var rdErr os.Error
    var line string
    var str []string = make([]string, 2)

    file,err := os.Open(conf, os.O_RDONLY,0)
    if err != nil {
        fmt.Println("GO-WARP, error opening the PHOLD configuration file:",err)
        os.Exit(1)
    }
    rd := bufio.NewReader(file) 
    
    for i:=0;rdErr == nil;i++ { 
        line, rdErr = rd.ReadString('\n') 
        if len(line) > 0 { 
            str = strings.Split(line,"#",2)
            num := strings.Replace(str[0], "\t", "", -1)
            num = strings.Replace(num, " ", "", -1)
            if i == 0 {
                density,_ = strconv.Atof(num)
            } else if i == 1 {
                t,_ := strconv.Atoi(num)
                endtime = DT.Time(t)
            } else if i == 2 {
                nFPops,_ = strconv.Atoi(num)
            }
        }
    }
    fmt.Println("GO-WARP: read from file:",density,endtime,nFPops)

    lpnum = nlp
    entitynum = nent
    n_events = int(float(nent)*density)
    randGen = Random.RandInit(int64(lpnum+entitynum))
    idcount = 0

    initEv = make([]DT.Event, n_events)

    Sim.Setup(lpnum, endtime, ProcessEvent)

    for i:=0;i<n_events;i++ {
        e := generateEvent(nil)
        initEv[i] = *e
    }

    simterm = false
    n_term = 0
}


func launchLP(index DT.Pid, n_entity int) {
    var data *Local.LocalData
    data = Sim.Initialize(index)

    getEvents(index,data)

    Sim.Simulate(data)

    terminate(data)
}


/* each event in the system is generated in this function */
func generateEvent(oldev *DT.Event) *DT.Event {
    var mitt int
    var dest int
    var id int32
    var t DT.Time

    if oldev==nil {
        mitt = int(randGen.RandIntUniform(0,int32(entitynum)))
        t = 0 // basetime
    } else {
        mitt = oldev.Type.To
        t = oldev.Time // basetime
    }

    dest = int(randGen.RandIntUniform(0,int32(entitynum-1)))
    for ;mitt==dest; {
        dest = int(randGen.RandIntUniform(0,int32(entitynum-1)))
    }
    id = idcount
    idcount++
    t += DT.Time(randGen.RandIntExponential())

    e := DT.CreateEvent(id, t, DT.Info{mitt,dest,0})
    return e
}


/* each LP gets his events from those that have been generated at start up */
func getEvents(index DT.Pid, data *Local.LocalData) {
    for i:=0;i<n_events;i++ {
        if DT.Pid(initEv[i].Type.From/(entitynum/lpnum)) == index {
            data.FutureEvents.Insert(&initEv[i])
        }
    }
}


func ProcessEvent(ev *DT.Event, l *Local.LocalData) {
    newev := generateEvent(ev)
    lp := e2lp(newev.Type.To,entitynum,lpnum)
    Sim.NoticeEvent(newev, DT.Pid(lp), l)
    compute()
}


func e2lp(e, en, lpn int)  int {
    var lp int
    m := en % lpn
    d := en / lpn
    if m == 0 {
        lp = e / d
    } else {
        lp = e / (d+1)
        if lp >= m {
            ee := e - m*(d+1)
            lp = m + ee/d
        }
    }
    return lp
}


func compute() float64 {
	var z,x float64
	z=2
	x=0.5
	
    for i:=0;i<nFPops/5;i++ {
        x = 0.5 * x * (3 - z * x * x)
    }
    return x
}


func terminate(data *Local.LocalData) {
    if !simterm {
        simterm = true
        endT = time.Nanoseconds()-startT
    }
    
    print.Lock()

    fmt.Println("|----------------------------------------------|")
    fmt.Println("LOGICAL PROCESS",data.IndexLP)
    fmt.Println("Number of processed events =",data.N_PROCESSED)

    n_term++
    print.Unlock()
}


func printStats(time int64) {
    print.Lock()

    fmt.Println("SIMULATION IS COMPLETED: TIME REACHED VALUE",endtime)
    fmt.Println("Wall Clock Time spent (ms):",float(time)/float(1e6))

    fmt.Println("Number of GVT evaluations:",Shared.N_gvt)

    sum := 0
    for i:=0;i<lpnum;i++ {
        sum += Shared.N_rollback[i]
    }
    fmt.Println("Total number of rollbacks:",sum)

    print.Unlock()
}
