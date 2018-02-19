# https://gist.github.com/FredLoney/5454553

from __future__ import unicode_literals

import inspect
import logging

HEADER_FMT = "Call stack at %s, line %d in function %s, frames %d to %d of %d:"
"""The log header message formatter."""

STACK_FMT = "%s, line %d in function %s."
"""The log stack message formatter."""

def log_stack(logger=None, limit=None, start=0):
    """
    Prints the call stack at the point of the caller to the given log.

    Example:

    >>> import logging
    >>> logger = logging.getLogger(__name__)
    >>> logger.setLevel(logging.DEBUG)
    >>>
    >>> import logging_helper
    >>>
    >>> def outer():
    ...     middle()
    >>>
    >>> def middle():
    ...     inner()
    >>>
    >>> def inner():
    ...     logging_helper.log_stack(logger, 2)
    >>>
    >>> outer()
    130424-12:17:35,722 __main__ DEBUG:
         Call stack at /snippet/test_logging_helper.py, line 13 in function inner, frames 2 to 3 of 11:
    130424-12:17:35,722 __main__ DEBUG:
    	 /snippet/test_logging_helper.py, line 10 in function middle.
    130424-12:17:35,722 __main__ DEBUG:
    	 /snippet/test_logging_helper.py, line 7 in function outer.

    @param logger: the logger to use (default use the root logger)
    @param limit: the number of frames to print (default print all remaining frames)
    @param start: the offset of the first frame preceding the caller to print (default 0)
    """

    # Use the default logger, if necessary.
    if not logger:
        from canopsis.logger import Logger
        logger = Logger.get('callstack', '/tmp/callstack.cps.log', level=logging.DEBUG)

    # The call stack.
    stack = inspect.stack()

    # The penultimate frame is the caller to this function.
    here = stack[1]

    # The index of the first frame to print.
    begin = start + 2

    # The index of the last frame to print.
    if limit:
        end = min(begin + limit, len(stack))
    else:
        end = len(stack)

    # Print the stack to the logger.
    file, line, func = here[1:4]
    logger.debug(HEADER_FMT % (file, line, func, start + 2, end - 1, len(stack) - 1))
    # Print the next frames up to the limit.
    for frame in stack[begin:end]:
        file, line, func = frame[1:4]
        logger.debug(STACK_FMT % (file, line, func))