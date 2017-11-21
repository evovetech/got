#!/bin/bash

function call_git() {
  call_cmd git "$@"
}

function git_rm() {
  call_cmd rm -f "$1"
  call_git add -- "$1"
}

function git_resolve() {
  call_git checkout "--$1" -- "$2"
  call_git add -- "$2"
}

function reset() {
  call_git merge --abort
  call_git checkout "$target"
  call_git reset --hard "$target_sha"
  call_git branch -D "$1"
  return 0
}

function resolve_ours() {
  local status="$1"
  local f="$2"
  case "$status" in
    D* | UA* )
      git_rm "$f"
      ;;
    * )
      git_resolve ours "$f"
      ;;
  esac
}

function resolve_theirs() {
    local status="$1"
    local f="$2"
    case "$status" in
      ?D* | AU* )
        git_rm "$f"
        ;;
      * )
        git_resolve theirs "$f"
        ;;
    esac
}

function resolve_file() {
  local status="$1"
  local cmd="$2"
  local f="$3"

  if [[ "$cmd" == "ours" ]]; then
    resolve_ours "$status" "$f"
  else
    resolve_theirs "$status" "$f"
  fi
}

function merge_resolve_conflicts() {
  local cmd="$1"
  local _abort=0

  case "$cmd" in
    "ours" )
      ;;
    * )
      cmd="theirs"
      ;;
  esac

  log "---"
  log "merge_resolve_conflicts --$cmd"
  indent

  while read f; do
    local _unmerged="unmerged"
    log "---"
    indent

    local status=""
    log "f=" "$f"
    if [[ -n "$f" ]]; then
      status="$(git status -s -- "$f")"
    fi
    log "in='$status'"

    local _st=0
    case "$status" in
      "" )
        log "unknown:" "$f"
        _st=1
        ;;
      DD* )
        git_rm "$f"
        ;;
      * )
        resolve_file "$status" "$cmd" "$f"
        _st=$?
        ;;
    esac

    if [[ $_st -ne 0 ]]; then
      _abort=1
      log "-- Error: not found '$status' --"
    fi

    if [[ -n "$f" ]]; then
      status="$(git status -s -- "$f")"
    fi
    log "out='$status'"

    outdent
  done

  outdent
  log "---"

  return $_abort
}

function merge_delete_untracked() {
  local remove='?* '
  while read status; do
    case "$status" in
      \?\?* )
        local f=${status#$remove}
        log "untracked: $f"
        call_cmd rm -f "$f"
        ;;
    esac
  done <<< "$(git status -s -u)"
}

function merge_keep() {
  local which="$1"
  local cmd="theirs"

  case "$which" in
    "parent" )
      merge_commit="$parent_sha"
      merge_msg="$merge_parent_msg"
      cleanup_branch="$merge_parent"
      call_git checkout -b "$merge_parent" "$target_sha"
      ;;
    "target" )
      merge_commit="$target_sha"
      merge_msg="$merge_target_msg"
      cleanup_branch="$merge_target"
      call_git checkout -b "$merge_target" "$parent_sha"
      ;;
    * )
      which="theirs"
      merge_commit="$merge_parent"
      merge_msg="merging conflicts, resolving with ${parent} changes"
      cleanup_branch="$merge_commit"
      ;;
  esac

  ## perform initial merge
  call_git merge --no-commit -X "$cmd" "$merge_commit"

  # resolve unmerged files
  git diff --name-only --diff-filter=UXB | merge_resolve_conflicts "$cmd"
  local _abort=$?
  merge_delete_untracked

  if [[ ${_abort} -eq 0 ]]; then
    call_git commit -am "$merge_msg"
    _abort=$?
  fi

  if [[ ${_abort} -ne 0 ]]; then
    reset "$cleanup_branch"
  else
    if [[ "$which" == "theirs" ]]; then
      # cleanup the temp merge branches
      call_git branch -D "$cleanup_branch"
    fi
  fi
  return ${_abort}
}

function try_normal_merge() {
  call_git checkout "$target"
  call_git merge --commit -X ignore-all-space "$parent"
  local _st=$?
  if [[ ${_st} -ne 0 ]]; then
    reset ""
  fi
  return ${_st}
}

function ff_merge() {
  call_git checkout "$target"
  call_git merge --ff-only "$merge_target"
  local _st=$?
  if [ ${_st} -ne 0 ]; then
    reset "$merge_parent"
  fi
  call_git branch -D "$merge_target"
  return ${_st}
}
