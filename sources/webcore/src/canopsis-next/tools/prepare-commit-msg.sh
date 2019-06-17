#!/bin/bash
#
# Automatically adds branch name and branch description to every commit message.
# Modified from the gist here https://gist.github.com/bartoszmajsak/1396344
#

# This way you can customize which branches should be skipped when
# prepending commit message.

# Export the vars in .env into your shell:

if [ -f .env.local ]; then
  echo "FILE WAS FOUND"
  export $(cat .env.local | sed 's/#.*//g' | xargs)
fi

echo "CANOPSIS:$CANOPSIS_PREPARE_COMMIT_MSG_HOOK"

if [ "$CANOPSIS_PREPARE_COMMIT_MSG_HOOK" == "true" ]; then
  echo "INTO"

  if [ -z "$BRANCHES_TO_SKIP" ]; then
    BRANCHES_TO_SKIP=(master develop)
  fi

  if [ -z "$ROOT_PATH" ]; then
    $ROOT_PATH=./
  fi

  # Get branch name and description
  BRANCH_NAME=$(git --git-dir "$ROOT_PATH/$GIT_DIR" --work-tree $ROOT_PATH branch | grep '*' | sed 's/* //')

  IFS='/'; BRANCH_NAME_PARTS=($BRANCH_NAME); unset IFS;

  COMMIT_PREFIX="${BRANCH_NAME_PARTS[0]}(${BRANCH_NAME_PARTS[1]}): "

  # Branch name should be excluded from the prepend
  BRANCH_EXCLUDED=$(printf "%s\n" "${BRANCHES_TO_SKIP[@]}" | grep -c "^$BRANCH_NAME$")

  # A developer has already prepended the commit in the format BRANCH_NAME]
  BRANCH_IN_COMMIT=$(grep -c $COMMIT_PREFIX "$ROOT_PATH/$1")

  if [ -n "$BRANCH_NAME" ] && ! [[ $BRANCH_EXCLUDED -eq 1 ]] && ! [[ $BRANCH_IN_COMMIT -ge 1 ]]; then
    sed -i.bak -e "1s~^~$COMMIT_PREFIX~" "$ROOT_PATH/$1"
  fi
fi