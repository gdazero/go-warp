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
if [ $# -lt 1 ]
then
    echo "statistics: WRONG NUMBER OF PARAMETERS"
    echo "usage: logfile"
fi

file=$1

tmp="./tmp"
grep "Wall Clock Time" $file >> $tmp
cat $tmp | awk 'BEGIN{avg=0} {avg+=$6} END{printf "Average: %5.3f\n",avg/NR}'
media=$(cat $tmp | awk 'BEGIN{avg=0} {avg+=$6} END{print avg/NR}')

cat $tmp | awk 'BEGIN{max=0} {if ($6>max) max=$6} END{printf "Max: %5.3f\n",max}'

cat $tmp | awk 'BEGIN{min=1000000000} {if ($6<min) min=$6} END{printf "Min: %5.3f\n",min}'

cat $tmp | awk 'BEGIN{var=0} {var+=($6-'$media')^2} END{printf "Standard deviation: %5.3f\n", sqrt(var/NR)}'

grep "GVT" $file | awk 'BEGIN{avg=0} {avg+=$5} END{printf "GVT average: %5.3f\n",avg/NR}'
grep "rollback" $file | awk 'BEGIN{avg=0} {avg+=$5} END{printf "Rollback average: %5.3f\n",avg/NR}'

rm $tmp
