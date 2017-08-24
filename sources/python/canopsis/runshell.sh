#!/bin/bash

# Build PYTHONPATH
canopsis_python_path="$PYTHONPATH"

projects=$(find -type d -maxdepth 1 -mindepth 1 -printf "%f\n" 2>/dev/null)

for project in $projects
do
     project=$(echo $project | cut -d'/' -f2)
     tmppath="$(pwd)/$project"
     canopsis_python_path="$tmppath:$canopsis_python_path"
done

# Build python command
python_command="ipython"

which $python_command >/dev/null 2>&1

if [ $? -ne 0 ]
then
     python_command="python"
fi

# Execute shell
PYTHONPATH="$canopsis_python_path" $python_command
