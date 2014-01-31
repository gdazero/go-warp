#
# This file is part of GO-WARP.  GO-WARP is free software, you can
# redistribute it and/or modify it under the terms of the Revised BSD License.
#
# For more information please see the LICENSE file.
#
# Copyright 2014, Gabriele D'Angelo, Moreno Marzolla, Pietro Ansaloni
# Computer Science Department, University of Bologna, Italy
#

#!/bin/bash
if [ $# -lt 2 ]
then
    echo "test-scalability: WRONG NUMBER OF PARAMETERS"
    echo "usage: ./test-scalability.sh #entities #repetitions"
    exit
fi

ENTITIES=$1
REPETITIONS=$2

rm -f logs/*

# number of cores used
for GOMAX in 1 2 3 4
do
	# number of LPs
	for LP in 1 2 3 4
	do
		# it is really a nonsense to simulate many LPs on the same CPU core
		if [ $LP -le $GOMAX ]
		then
		  echo "LP: $LP, GOMAXPROCS = $GOMAX"
		  export GOMAXPROCS=$GOMAX; ./test-main.sh $LP $ENTITIES $REPETITIONS
		fi
	done
	
	sleep 3
done
