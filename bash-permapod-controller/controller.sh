#!/bin/bash

echo "Controller starting"

reconcile () {
    ns=$1
    echo "Reconciling Namespace: $ns"

    kubectl get pod -n $ns permapod 2>&1 > /dev/null

    return_value=$?

    if [ $return_value -eq 1 ]; then
        echo "Permapod does not exist, creating pod"
        kubectl run permapod -n $ns --image ubuntu -- sleep 1000000
    else
        echo "Permapod EXISTS, doing nothing"
    fi
}

# The function needs to be exported so that is can be called with xargs.
export -f reconcile

kubectl get namespaces -w -o=jsonpath='{.metadata.name}{"\n"}' | xargs -L1 bash -c 'reconcile "$@"' _
