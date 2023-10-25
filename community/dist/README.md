# Makefile params

When running make with the provided Makefile, override these two vars :


Variable     | Default value | CI value
-------------|---------------|----------------------------
VERSION      | develop       | Tagged version ($CI_COMMIT_TAG)
BUILD_TARGET | el7, el8      | el8
