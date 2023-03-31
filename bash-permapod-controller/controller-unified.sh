#!/usr/bin/env bash

echo "Controller starting"

# Kill the whole process group as we want to kill
# both background jobs when terminating this script.
# Note the nested `trap` resets the script's SIGTERM
# response to the default `kill` behaviour; we don't
# want to cause a `trap` inception. The `kill -- -$$`
# sends a SIGTERM to the whole process group, thus
# killing all descendants, see:
# https://stackoverflow.com/a/2173421/16406860
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM

reconcile_ns() {
  local ns=$1
  echo "Reconciling namespace: $ns"
  kubectl get pod -n $ns permapod >/dev/null 2>&1
  return_value=$?

  if [ $return_value -eq 1 ]; then
    echo "Permapod does not exist, creating pod"
    kubectl run permapod -n $ns --image ubuntu -- sleep 10000
  else
    echo "Permapod EXISTS, doing nothing"
  fi
}

# In bash all commands in a pipeline execute in a subshell, so no need to
# wrap this command in parenthesis. Appending the ampersand will turn it
# into a background process.
kubectl get namespaces -w -o=jsonpath='{.metadata.name}{"\n"}' | while IFS= read -r ns; do
  reconcile_ns "$ns"
done &

kubectl get pods -A -w -o=jsonpath='{.metadata.namespace}{"\n"}' | while IFS= read -r ns; do
  reconcile_ns "$ns"
done &

# Wait forever, or until we receive a Term signal.
# If we don't wait here this script finishes and
# there is no way to trap the exit code and act on it.
wait
