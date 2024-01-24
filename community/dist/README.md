# Makefile params

When running make with the provided Makefile, override these two vars :


Variable     | Default value | CI value
-------------|---------------|----------------------------
VERSION      | HEAD          | Tagged version ($CI_COMMIT_TAG)
BUILD_TARGET | el8, el9      | el8
