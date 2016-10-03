# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2014 Jonathan Labéjof <jonathan.labejof@gmail.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.
# --------------------------------------------------------------------

"""Code from http://code.activestate.com/recipes/277940-decorator-for-\
bindingconstants-at-compile-time/

Decorator for automatic code optimization.
If a global is known at compile time, replace it with a constant.
Fold tuples of constants into a single constant.
Fold constant attribute lookups into a single constant.

Modifications:

- Add constants values from opmap constants (STORE_GLOBAL, etc.) in order to
    avoid to update globals.
- Modify verbose argument which is None or use a function with one argument
    which can be bound to a print function or a logging function.
- Set attributes from originary function such as __dict__, __module__, etc."""

from opcode import opmap, HAVE_ARGUMENT, EXTENDED_ARG

from types import FunctionType, ModuleType

from functools import reduce

from six import exec_, PY3
from six.moves import builtins



__all__ = [
    'SAFE_BUILTINS', 'safe_eval', 'safe_exec', 'bind_all', 'make_constants',
    'singleton_per_scope', 'getcodeobj'
]


BUILTIN_IO_PROPS = [
    'open', '__name__', '__debug__', '__doc__', '__import__', '__package__',
    'compile', 'copyright', 'credits', 'eval', 'execfile', 'exit', 'file',
    'globals', 'help', 'input', 'intern', 'license', 'locals', 'open', 'print',
    'quit', 'raw_input', 'reload'
]  #: set of builtin objects to remove from a safe builtin.


SINGLETONS_PER_SCOPES = {}

def singleton_per_scope(_cls, _scope=None, _renew=False, *args, **kwargs):
    """Instanciate a singleton per scope."""

    result = None

    singletons = SINGLETONS_PER_SCOPES.setdefault(_scope, {})

    if _renew or _cls not in singletons:
        singletons[_cls] = _cls(*args, **kwargs)

    result = singletons[_cls]

    return result


def _safebuiltins():
    """Construct a safe builtin environment without I/O functions.

    :rtype: dict"""

    result = {}

    objectnames = [
        objectname for objectname in dir(builtins)
        if objectname not in BUILTIN_IO_PROPS
    ]

    for objectname in objectnames:
        result[objectname] = getattr(builtins, objectname)

    return result

SAFE_BUILTINS = {'__builtins__': _safebuiltins()}  #: safe builtins.


def _safe_processing(nsafefn, source, _globals=None, _locals=None):
    """Do a safe processing of input fn in using SAFE_BUILTINS.

    :param fn: function to call with input parameters.
    :param source: source object to process with fn.
    :param dict _globals: global objects by name.
    :param dict _locals: local objects by name.
    :return: fn processing result"""

    if _globals is None:
        _globals = SAFE_BUILTINS

    else:
        _globals.update(SAFE_BUILTINS)

    return nsafefn(source, _globals, _locals)


def safe_eval(source, _globals=None, _locals=None):
    """Process a safe evaluation."""

    return _safe_processing(eval, source, _globals, _locals)


def safe_exec(source, _globals=None, _locals=None):
    """Do a safe python execution."""

    return _safe_processing(exec_, source, _globals, _locals)


STORE_GLOBAL = opmap['STORE_GLOBAL']
LOAD_GLOBAL = opmap['LOAD_GLOBAL']
LOAD_CONST = opmap['LOAD_CONST']
LOAD_ATTR = opmap['LOAD_ATTR']
BUILD_TUPLE = opmap['BUILD_TUPLE']
JUMP_FORWARD = opmap['JUMP_FORWARD']

WRAPPER_ASSIGNMENTS = ('__doc__', '__annotations__', '__dict__', '__module__')


def _make_constants(func, builtin_only=False, stoplist=None, verbose=None):
    """Generate new function where code is an input function code with all
    LOAD_GLOBAL statements changed to LOAD_CONST statements.

    :param function func: code function to transform.
    :param bool builtin_only: only transform builtin objects.
    :param list stoplist: attribute names to not transform.
    :param function verbose: logger function which takes in parameter a message

    .. warning::
        Be sure global attributes to transform are not resolved dynamically."""

    result = func

    if stoplist is None:
        stoplist = []

    try:
        fcode = func.__code__
    except AttributeError:
        return func        # Jython doesn't have a __code__ attribute.
    newcode = list(fcode.co_code) if PY3 else [ord(co) for co in fcode.co_code]
    newconsts = list(fcode.co_consts)
    names = fcode.co_names
    codelen = len(newcode)

    env = vars(builtins).copy()
    if builtin_only:
        stoplist = dict.fromkeys(stoplist)
        stoplist.update(func.__globals__)
    else:
        env.update(func.__globals__)

    # First pass converts global lookups into constants
    changed = False
    i = 0
    while i < codelen:
        opcode = newcode[i]
        if opcode in (EXTENDED_ARG, STORE_GLOBAL):
            return func    # for simplicity, only optimize common cases
        if opcode == LOAD_GLOBAL:
            oparg = newcode[i + 1] + (newcode[i + 2] << 8)
            name = fcode.co_names[oparg]
            if name in env and name not in stoplist:
                value = env[name]
                for pos, val in enumerate(newconsts):
                    if val is value:
                        break
                else:
                    pos = len(newconsts)
                    newconsts.append(value)
                newcode[i] = LOAD_CONST
                newcode[i + 1] = pos & 0xFF
                newcode[i + 2] = pos >> 8
                changed = True
                if verbose is not None:
                    verbose("{0} --> {1}".format(name, value))
        i += 1
        if opcode >= HAVE_ARGUMENT:
            i += 2

    # Second pass folds tuples of constants and constant attribute lookups
    i = 0
    while i < codelen:

        newtuple = []
        while newcode[i] == LOAD_CONST:
            oparg = newcode[i + 1] + (newcode[i + 2] << 8)
            newtuple.append(newconsts[oparg])
            i += 3

        opcode = newcode[i]
        if not newtuple:
            i += 1
            if opcode >= HAVE_ARGUMENT:
                i += 2
            continue

        if opcode == LOAD_ATTR:
            obj = newtuple[-1]
            oparg = newcode[i + 1] + (newcode[i + 2] << 8)
            name = names[oparg]
            try:
                value = getattr(obj, name)
            except AttributeError:
                continue
            deletions = 1

        elif opcode == BUILD_TUPLE:
            oparg = newcode[i + 1] + (newcode[i + 2] << 8)
            if oparg != len(newtuple):
                continue
            deletions = len(newtuple)
            value = tuple(newtuple)

        else:
            continue

        reljump = deletions * 3
        newcode[i - reljump] = JUMP_FORWARD
        newcode[i - reljump + 1] = (reljump - 3) & 0xFF
        newcode[i - reljump + 2] = (reljump - 3) >> 8

        nclen = len(newconsts)
        newconsts.append(value)
        newcode[i] = LOAD_CONST
        newcode[i + 1] = nclen & 0xFF
        newcode[i + 2] = nclen >> 8
        i += 3
        changed = True
        if verbose is not None:
            verbose("new folded constant:{0}".format(value))

    if changed:

        codeobj = getcodeobj(newconsts, newcode, fcode, fcode)

        result = type(func)(
            codeobj,
            func.__globals__, func.__name__, func.__defaults__, func.__closure__
        )

        # set func attributes to result
        for prop in WRAPPER_ASSIGNMENTS:
            try:
                attr = getattr(func, prop)
            except AttributeError:
                pass
            else:
                setattr(result, prop, attr)

    return result


def getcodeobj(consts, intcode, newcodeobj, oldcodeobj):
    """Get code object from decompiled code.

    :param list consts: constants to add in the result.
    :param list intcode: list of byte code to use.
    :param newcodeobj: new code object with empty body.
    :param oldcodeobj: old code object.
    :return: new code object to produce."""

    # get code string
    if PY3:
        codestr = bytes(intcode)

    else:
        codestr = reduce(lambda x, y: x + y, (chr(b) for b in intcode))

    # get vargs
    vargs = [
        newcodeobj.co_argcount, newcodeobj.co_nlocals, newcodeobj.co_stacksize,
        newcodeobj.co_flags, codestr, tuple(consts), newcodeobj.co_names,
        newcodeobj.co_varnames, newcodeobj.co_filename, newcodeobj.co_name,
        newcodeobj.co_firstlineno, newcodeobj.co_lnotab,
        oldcodeobj.co_freevars, newcodeobj.co_cellvars
    ]
    if PY3:
        vargs.insert(1, newcodeobj.co_kwonlyargcount)

    # instanciate a new newcodeobj object
    result = type(newcodeobj)(*vargs)

    return result


_make_constants = _make_constants(_make_constants)  # optimize thyself!


def bind_all(morc, builtin_only=False, stoplist=None, verbose=None):
    """Recursively apply constant binding to functions in a module or class.

    Use as the last line of the module (after everything is defined, but
    before test code). In modules that need modifiable globals, set
    builtin_only to True.

    :param morc: module or class to transform.
    :param bool builtin_only: only transform builtin objects.
    :param list stoplist: attribute names to not transform.
    :param function verbose: logger function which takes in parameter a message
    """

    if stoplist is None:
        stoplist = []

    def _bind_all(morc, builtin_only=False, stoplist=None, verbose=False):
        """Internal bind all decorator function.
        """
        if stoplist is None:
            stoplist = []

        if isinstance(morc, (ModuleType, type)):
            for k, val in list(vars(morc).items()):
                if isinstance(val, FunctionType):
                    newv = _make_constants(
                        val, builtin_only, stoplist, verbose
                    )
                    setattr(morc, k, newv)
                elif isinstance(val, type):
                    _bind_all(val, builtin_only, stoplist, verbose)

    if isinstance(morc, dict):  # allow: bind_all(globals())
        for k, val in list(morc.items()):
            if isinstance(val, FunctionType):
                newv = _make_constants(val, builtin_only, stoplist, verbose)
                morc[k] = newv
            elif isinstance(val, type):
                _bind_all(val, builtin_only, stoplist, verbose)
    else:
        _bind_all(morc, builtin_only, stoplist, verbose)


@_make_constants
def make_constants(builtin_only=False, stoplist=None, verbose=None):
    """Return a decorator for optimizing global references.

    Replaces global references with their currently defined values.
    If not defined, the dynamic (runtime) global lookup is left undisturbed.
    If builtin_only is True, then only builtins are optimized.
    Variable names in the stoplist are also left undisturbed.
    Also, folds constant attr lookups and tuples of constants.
    If verbose is True, prints each substitution as is occurs.

    :param bool builtin_only: only transform builtin objects.
    :param list stoplist: attribute names to not transform.
    :param function verbose: logger function which takes in parameter a message
    """

    if stoplist is None:
        stoplist = []

    if isinstance(builtin_only, type(make_constants)):
        raise ValueError("The bind_constants decorator must have arguments.")

    return lambda func: _make_constants(func, builtin_only, stoplist, verbose)
