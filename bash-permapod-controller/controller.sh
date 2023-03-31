#!/bin/bash

k=$(which kubectl)
current_user=$(whoami)

echo "cleanup all kubectl(s) running in the background for user ${current_user}"
for process in $(ps -eo user,pid,pgid,tgid,args| awk 'NR == 1 || ($3 != -1 && $2 != $3)' |grep kubectl|awk -v user=${current_user} '{if ($1 == user) {print $2}}'); do kill -9 ${process} > /dev/null 2>&1 ; done

reconcile () {
  local ns=${1}
  local k=$(which kubectl)

  echo "reconciling namespace: ${ns}"
  ${k} get po alpine -n ${ns} > /dev/null 2>&1
  local exit_code=${?}
  if [[ ${exit_code} == 1 ]]; then
    echo "pod not found in ns ${ns}... CREATING"
    ${k} -n ${ns} run alpine --image=alpine:latest -- sleep 10000000
  else
    echo "pod found in ns ${ns}... DO NOTHING"
  fi
}

export -f reconcile

echo "Starting controller in the background..."
echo -e "\nView controller logs with: \ntail -F /tmp/controller.log\n"

${k} get ns -w -ojsonpath='{.metadata.name}{"\n"}'|xargs -L1 bash -c 'reconcile "${@}" >> /tmp/controller.log & ' _ > /dev/null 2>&1 &
${k} get po -A -w -ojsonpath='{.metadata.namespace}{"\n"}'|xargs -L1 bash -c 'reconcile "${@}" >> /tmp/controller.log & ' _ > /dev/null 2>&1 &

