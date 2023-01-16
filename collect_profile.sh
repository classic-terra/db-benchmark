#!/bin/bash

counter=0
rm -rf profile
mkdir profile

while true; do
    # collecting inuse memory profile for each 30s
    # inuse_space: memory allocated but not yet released
    # inuse_objects: number of objects allocated but not yet released
    echo "collecting inuse memory profile"
    go tool pprof --inuse_space -png http://localhost:6060/debug/pprof/heap > profile/inuse_space_$counter.png
    go tool pprof --inuse_objects -png http://localhost:6060/debug/pprof/heap > profile/inuse_objects_$counter.png
    counter = $((counter + 1))
    
    # check counter divisible by two, checking for every 1 minute
    if [ $((counter % 2)) -eq 0 ]; then
        # if counter is divisible by two, then collect cpu profile
        # alloc_space: total amount of memory allocated
        # alloc_objects: total number of objects allocated
        echo "collecting total allocations"
        go tool pprof --alloc_space -png http://localhost:6060/debug/pprof/heap > profile/alloc_space_$counter.png
        go tool pprof --alloc_objects -png http://localhost:6060/debug/pprof/heap > profile/alloc_objects_$counter.png
    fi

    sleep 30
done