/*
	GO-WARP: a Time Warp simulator written in Go
	http://pads.cs.unibo.it
  
	This file is part of GO-WARP.  GO-WARP is free software, you can
	redistribute it and/or modify it under the terms of the Revised BSD License.

	For more information please see the LICENSE file.

	Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
	Computer Science Department, University of Bologna, Italy
*/

package Random

import "fmt"
import "math"

/*
 * Linear Congruential Generator LGC 16807
 */

type RNG struct{ Seed, Prev int64 }


const (
    module int64 = 1<<31 -1	// RNG module
    coeff int64 = 16807		// RNG coefficient
    LAMBDA = 5.0		// exponential distribution parameter
    SD = 1.0			// normal distribution parameter
)


func RandInit(seed int64) *RNG{
    var rngptr *RNG = new(RNG)

    if seed == 0 {
        fmt.Println("GO-WARP, ERROR: NON NULL SEED EXPECTED!")
    }

    *rngptr = RNG{seed, seed}
    return rngptr
}


/* This function randomically generates a float number in the range 0..1 */
func (rng *RNG) RandFloat() float64 {
    var n int64
    var fl float64

    n = (coeff*(rng.Prev)) % module
    rng.Prev = n
    
    fl = float64(n)/float64(module)
    
    return fl
}


func (rng *RNG) RandIntUniform(min int32, max int32) int32 {
    var ret int32 = 0
    
    if min<=max {
        ret = int32(rng.RandFloat()*float64(max-min+1))+min
    }
    return ret
}


func (rng *RNG) RandIntExponential() int32 {
    return int32( -LAMBDA * math.Log(rng.RandFloat()) + 1 )
}
