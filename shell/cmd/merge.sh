#!/bin/bash

while [[ $# -gt 0 ]]; do
case "$1" in
  -h | --help)
      echo "merge help"  # Call your function
      # no shifting needed here, we're done.
      exit 0
      ;;
  -s | --strategy)
    strategy="$2"
    shift 2
    ;;
  -*)
      echo "Error: Unknown option: $1" >&2
      shift 1
      ;;
  *)
    break
      ;;
esac
done

parent="$1"
parent_sha="$(commit_sha $parent)"
target="$(git symbolic-ref --short HEAD)"
target_sha="$(commit_sha $target)"

merge_parent="merge_parent_${parent_sha}"
merge_target="merge_target_${target_sha}"

merge_parent_msg="merge ${parent} into ${target} -- CONFLICTS -- resolving with ${parent} changes"
merge_target_msg="merge ${target} into ${parent} -- CONFLICTS -- resolving with ${target} changes"

## import functions
src_lib "merge/core"

echo "parent=$parent, strategy=$strategy"

## First try to merge ignoring whitespace
call_function try_normal_merge
if [ $? -eq 0 ]; then
  log "-- normal merge SUCCESS --"
  log ""
  exit 0
fi

# first create a merge_parent branch on the target branch's commit
# and merge the parent branch. resolve conflicts with target branch's changes
call_function merge_keep parent
if [ $? -ne 0 ]; then
  exit 1
fi

# second create a merge_target branch on the parent branch's commit
# and merge the target branch. resolve conflicts with parent branch's changes
call_function merge_keep target
if [ $? -ne 0 ]; then
  exit 1
fi

# checkout the original target branch and fast-forward merge with the
# merge_parent branch
call_function ff_merge
if [ $? -ne 0 ]; then
  exit 1
fi

# Lastly, do a final merge of the merge_target branch
# into the target branch (now at merge_parent commit) in order to single out
# the conflicts into one merge commit.
# Here, resolve with the parent branch's changes
call_function merge_keep theirs
_exit=$?
if [[ ${_exit} -ne 0 ]]; then
  reset "$merge_parent"
fi
log ""

exit ${_exit}
