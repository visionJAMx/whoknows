#!/usr/bin/env bash

SCRIPT_DIR=$(dirname "$0")
PYTHON_SCRIPT_PATH="${1:-${SCRIPT_DIR}/backend/app.py}"

while true; do
    if python3 "$PYTHON_SCRIPT_PATH"; then
        echo "Script stopped normally. Restarting..." >&2
    else
        exit_code=$?
        echo "Script crashed with exit code ${exit_code}. Restarting..." >&2
    fi

    sleep 1
done
