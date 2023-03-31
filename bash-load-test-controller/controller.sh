#!/usr/bin/env bash

reconcile () {
    namespace=$1
    name=$2

    echo "Reconciling LoadTest (namespace = $namespace / name = $name)"

    loadtest=$(kubectl get loadtest --ignore-not-found -ojson -n $namespace $name)
	[[ "$loadtest" == "" ]] && echo "[LoadTest $namespace/$name]: Not found, ignoring." && exit 0

    status=$(echo $loadtest | jq -r '.status')
	[[ "$status" != "null" ]] && echo "[LoadTest $namespace/$name]: Already run, ignoring." && exit 0
    
    method=$(echo $loadtest | jq -r '.spec.method')
	address=$(echo $loadtest | jq -r '.spec.address')
	duration=$(echo $loadtest | jq -r '.spec.duration')

    echo "Running Load Test (Method = $method, Address = $address, Duration = $duration)"
    results=$(echo "$method $address" | vegeta attack -duration=$duration | vegeta report --type json)
    requests=$(echo $results | jq -r '.requests')

    kubectl patch loadtest -n $namespace $name --type='json' --patch='[{"op": "replace", "path": "/status", "value":{"finished": true, "requests": '$requests'}}]'
}

export -f reconcile

echo "Starting..."
kubectl get loadtests --all-namespaces --watch --output go-template='{{printf "%s %s\n" .metadata.namespace .metadata.name}}' | xargs -L1 bash -c 'reconcile "$@"' _