#!/bin/bash

# Get the current branch name
BRANCH_NAME=$(git branch --show-current)

# Exit if we can't get the branch name
if [ -z "$BRANCH_NAME" ]; then
    exit 0
fi

# Skip if we're on main, master, develop, or release branches
case "$BRANCH_NAME" in
    main|master|develop|release/*|hotfix/*)
        exit 0
        ;;
esac

# Extract commit type from branch name
COMMIT_TYPE=""
case "$BRANCH_NAME" in
    feat/*|feature/*)
        COMMIT_TYPE="feat"
        ;;
    fix/*|bugfix/*)
        COMMIT_TYPE="fix"
        ;;
    docs/*|doc/*)
        COMMIT_TYPE="docs"
        ;;
    style/*)
        COMMIT_TYPE="style"
        ;;
    refactor/*)
        COMMIT_TYPE="refactor"
        ;;
    test/*|tests/*)
        COMMIT_TYPE="test"
        ;;
    chore/*)
        COMMIT_TYPE="chore"
        ;;
    perf/*|performance/*)
        COMMIT_TYPE="perf"
        ;;
    ci/*)
        COMMIT_TYPE="ci"
        ;;
    build/*)
        COMMIT_TYPE="build"
        ;;
esac

# If we found a commit type, prepend it to the commit message
if [ -n "$COMMIT_TYPE" ] && [ -f "$1" ]; then
    # Read the current commit message
    COMMIT_MSG=$(cat "$1")
    
    # Check if the message already starts with the commit type
    if [[ ! "$COMMIT_MSG" =~ ^$COMMIT_TYPE: ]]; then
        # Prepend the commit type to the message
        echo "$COMMIT_TYPE: $COMMIT_MSG" > "$1"
    fi
fi