#!/bin/bash

function fail() {
  echo ERROR: ${*} >&2

	cat<<-EOU
  
  Uses:
  $0 [-p <pald port>] get|set <service>
  $0 [-p <pald port>] del     <port>

  Default port is ${DEF_PALD_PORT}
	EOU

  exit 1
}

function validport() {
  [[ "${1}" =~ ^[0-9]+$ ]] ||
  {
    fail ${2} is not valid
  }

  [ "${1}" -le 65535 ] ||
  {
    fail ${2} number is not valid
  }
}

DEF_PALD_PORT=49200
PALD_PORT=${DEF_PALD_PORT}

while getopts ":p:" opt
do
  case ${opt} in
    p )
      validport ${OPTARG} "pald port"
      PALD_PORT=${OPTARG}
      ;;
    \? )
      fail invalid parameter: -${OPTARG}
      ;;
  esac
done

shift "$((OPTIND - 1))"

case ${1} in
  get | set )
    curl --fail --data service=${2} http://localhost:${PALD_PORT}/${1}
    ;;
  del )
    validport ${2} "resource port"
    curl --fail --data port=${2} http://localhost:${PALD_PORT}/${1}
    ;;
  * )
    fail invalid request name: ${1}
    ;;
esac

