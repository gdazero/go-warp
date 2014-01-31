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
if [ $# -lt 3 ]
then
  echo "test-main: WRONG NUMBER OF PARAMETERS"
  echo "usage: ./testmain.sh #LPs #entities #repetitions"
  exit 1
fi

LPS=$1
ENTITIES=$2
REPETITIONS=$3

out="logs/out.$LPS-$ENTITIES"

# as many runs as the number of repetitions
k=0
while [ $k -lt $REPETITIONS ]
do
  k=$((k+1))
  echo "--- REPETITION NUMBER: $k ----------------------------------" >> $out
  # simulator execution
  ./builds/Main.out $LPS $ENTITIES >> $out 2>&1
done

# counting the runs that have correctly reached the Wall Clock Time
echo -n "Simulation OK: "
grep "Wall Clock Time" $out | wc -l

# some statistics on the runs
./statistics.sh $out

# log cleaning
make cleanlog
echo
