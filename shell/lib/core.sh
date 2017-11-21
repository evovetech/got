#!/bin/bash

function commit_sha() {
  echo $(git rev-parse --short $1)
}

function indent() {
  (( INDENT_LEVEL++ ))
}

function outdent() {
  (( INDENT_LEVEL-- ))
}

function current_indention() {
  echo $(( INDENT_LEVEL * INDENT_SIZE ))
}

function println_indent() {
  echo "$( printf "%$1s%s" '' "${@:2}" )"
}

function println() {
  println_indent $(current_indention) "$*"
}

function subprint() {
  local indent=$(( (INDENT_LEVEL + 1) * INDENT_SIZE ))
  println_indent $indent "~ " "$*"
}

function log() {
  println "$@"
}

function log_cmd_error() {
  local _st=$?
  while read err; do
    if [[ -n "$err" ]]; then
      subprint "ERROR: '$err'" 1>&2
      _st=1
    fi
  done
  return $_st
}

function log_cmd() {
  local _st=$?
  while read msg
  do
    if [[ -n "$msg" ]]; then
      subprint "$msg"
    fi
  done
  return $_st
}

function call_cmd() {
  log \$ "$@"
  if [[ -n "$AAMOB_VERBOSE" ]]; then
    log_cmd <<< "$( "$@" 2> >(log_cmd_error) )"
  else
    log_cmd <<< "$( "$@" 1>/dev/null 2> >(log_cmd_error) )"
  fi
}

function call_function() {
  log ""
  log ""
  log \$ "$@"
  log "---------"

  indent
  $1 "${@:2}"
  local _st=$?
  outdent

  local msg='SUCCESS'
  if [[ $_st -ne 0 ]]; then
    msg='FAILURE'
  fi

  log "---------"
  log "  ~ $msg -- '"$@"'"
  return $_st
}

function print_args() {
  log "$1 {"
  indent
  local a
  for a in "${@:2}"; do
    log "'$a'"
  done
  outdent
  log "}"
}
