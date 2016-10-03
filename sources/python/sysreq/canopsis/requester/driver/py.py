# -*- coding: utf-8 -*-

# --------------------------------------------------------------------
# The MIT License (MIT)
#
# Copyright (c) 2016 Jonathan Labéjof <jonathan.labejof@gmail.com>
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

"""Python driver module."""

__all__ = ['PyDriver', 'processcrud', 'create', 'read', 'update', 'delete']

from .base import Driver

from operator import (
    lt, le, eq, ne, ge, gt, not_, truth, is_, is_not, abs, add, floordiv,
    truediv, invert, mod, mul, neg, or_, pow, rshift, lshift, sub,
    xor, concat, countOf, indexOf, repeat, sequenceIncludes, iadd, iand,
    getitem, setitem, delitem, getslice, setslice, delslice, iconcat,
    ifloordiv, ilshift, imod, imul, ior, ipow, irepeat, irshift, isub,
    itruediv, ixor, contains
)

from re import match

from ..request.expr import FuncName
from ..request.crud.create import Create
from ..request.crud.read import Read, Join
from ..request.crud.update import Update
from ..request.crud.delete import Delete

from random import random

from soundex import getInstance

from md5 import md5

from time import time
from datetime import datetime

soundex = getInstance().soundex

DATETIMEFORMAT = '%Y-%m-%d %H:%M:%S'


class PyDriver(Driver):
    """In charge of accessing to data from a list of dictionaries or objects."""

    version = '0.1'
    name = 'py'  # driver name

    def __init__(self, values=None, *args, **kwargs):
        """
        :param list values: list of data. Data are dictionaries. Default is [].
        """

        super(PyDriver, self).__init__(*args, **kwargs)

        self.values = [] if values is None else values

    def _process(self, request, crud, **kwargs):

        if kwargs:
            raise ValueError(
                'Driver {0} does not support additional arguments {1}'.format(
                    self, kwargs
                )
            )

        if request.query not in request.ctx:
            self._processquery(request.query, request.ctx)

        processcrud(request=request, items=self.values, crud=crud)

        return request

    def _processquery(self, query, ctx):

        if query in ctx:
            return ctx[query]

        else:
            pass

            ctx[query] = query


def create(items, create):
    """Apply input Create element to items.

    :param list items: items to process with input Create.
    :param Create create: data to add to input items.
    :return: created item.
    :rtype: list"""

    items.append(create.values)

    return items


def read(items, read):
    """Return application of input Read to items.

    :param list items: items to read.
    :param Read read: read resource to apply on items.
    :return: read list.
    :rtype: list
    """

    result = items

    if read.select():
        result = []
        for item in list(items):
            fitem = {}
            for sel in read.select():
                if sel in item:
                    fitem[sel] = item[sel]
            result.append(fitem)

    if read.offset():
        result = result[read.offset():]

    if read.limit():
        result = result[:read.limit()]

    if read.orderby():
        for orderby in read.orderby():
            result.sort(key=lambda item: item.get(orderby))

    if read.groupby():
        groupbyresult = {}
        _groupbyresult = []
        for groupby in read.groupby():
            if _groupbyresult:
                for item in _groupbyresult:
                    pass
            _groupbyresult = {groupby: []}

            for res in result:
                if groupby in res:
                    groupbyresult[groupby] = res.pop(groupby)

            #FIX: do the same for sub groupby...

    if read.join() not in ('FULL', None):
        raise NotImplementedError(
            'Driver {0} does not support join {1}'.format(
                self, read.join()
            )
        )

    items[:] = result

    return result


def update(items, update):
    """Apply update to items.

    :param list items: items to update.
    :param Update update: update rule.
    :return: updated items.
    :rtype: list"""

    result = []

    for item in items:
        if update.name in item:
            item[update.name] = update.values
            result.append(item)

    return result


def delete(items, delete):
    """Apply deletion rule to items.

    :param list items: items to modify.
    :param Delete delete: deletion rule.
    :rtype: list
    :return: modified/deleted items."""

    result = []

    if delete.names:
        for name in delete.names:
            for item in items:
                if name in item:
                    del item[name]
                    result.append(item)

    else:
        result = items
        items[:] = []

    return result


def processcrud(request, items, crud):
    """Apply the right rule on input items.

    :param list items: items to process.
    :param CRUDElement crud: crud rule to apply.
    :rtype: list
    :return: list"""

    if isinstance(crud, Create):
        processresult = create(items=items, create=crud)

    elif isinstance(crud, Read):
        processresult = read(items=items, read=crud)

    elif isinstance(crud, Update):
        processresult = update(items=items, update=crud)

    elif isinstance(crud, Delete):
        processresult = delete(items=items, delete=crud)

    request.ctx[crud] = processresult

    return items


_OPERTORS_BY_NAME = {
    FuncName.LT.value: lt,
    FuncName.LE.value: le,
    FuncName.EQ.value: eq,
    FuncName.NE.value: ne,
    FuncName.GE.value: ge,
    FuncName.GT.value: gt,
    FuncName.NOT.value: not_,
    FuncName.TRUTH.value: truth,
    FuncName.IS.value: is_,
    FuncName.ISNOT.value: is_not,
    FuncName.ABS.value: abs,
    FuncName.ADD.value: add,
    FuncName.FLOORDIV.value: floordiv,
    FuncName.DIV.value: truediv,
    FuncName.INDEX.value: indexOf,
    FuncName.INVERT.value: invert,
    FuncName.MOD.value: mod,
    FuncName.LIKE.value: match,
    FuncName.MUL.value: mul,
    FuncName.NEG.value: neg,
    FuncName.OR.value: or_,
    FuncName.POW.value: pow,
    FuncName.RSHIFT.value: rshift,
    FuncName.LSHIFT.value: lshift,
    FuncName.SUB.value: sub,
    FuncName.XOR.value: xor,
    FuncName.CONCAT.value: concat,
    FuncName.COUNTOF.value: countOf,
    FuncName.REPEAT.value: repeat,
    FuncName.INCLUDE.value: sequenceIncludes,
    FuncName.IADD.value: iadd,
    FuncName.IAND.value: iand,
    FuncName.IOR.value: ior,
    FuncName.IXOR.value: ixor,
    FuncName.GETITEM.value: getitem,
    FuncName.SETITEM.value: setitem,
    FuncName.DELITEM.value: delitem,
    FuncName.GETSLICE.value: getslice,
    FuncName.SETSLICE.value: setslice,
    FuncName.DELSLICE.value: delslice,
    FuncName.ICONCAT.value: iconcat,
    FuncName.IDIV.value: itruediv,
    FuncName.IFLOORDIV.value: ifloordiv,
    FuncName.ILSHIFT.value: ilshift,
    FuncName.IMOD.value: imod,
    FuncName.IMUL.value: imul,
    FuncName.IPOW.value: ipow,
    FuncName.IREPEAT.value: irepeat,
    FuncName.IRSHIFT.value: irshift,
    FuncName.ISUB.value: isub,
    FuncName.COUNT.value: len,
    FuncName.LENGTH.value: len,
    FuncName.AVG.value: lambda v: sum(v) / (len(v) or 1),
    FuncName.MEAN.value: lambda v: sum(v) / (len(v) or 1),
    FuncName.MAX.value: max,
    FuncName.MIN.value: min,
    FuncName.SUM.value: sum,
    FuncName.EXISTS.value: contains,
    FuncName.ISNULL.value: lambda data: data is None,
    FuncName.BETWEEN.value: lambda data, inf, sup: inf <= data <= sup,
    FuncName.IN.value: contains,
    FuncName.HAVING.value: contains,
    FuncName.UNION.value: lambda seq1, seq2: set(seq1) + set(seq2),
    FuncName.INTERSECT.value: lambda seq1, seq2: set(seq1) & set(seq2),
    FuncName.ALL.value: all,
    FuncName.ANY.value: any,
    FuncName.VERSION.value: PyDriver.version,
    FuncName.CONCAT.value: str.__add__,
    FuncName.ICONCAT.value: str.__add__,
    FuncName.REPLACE.value: str.replace,
    FuncName.SOUNDEX.value: soundex,
    FuncName.SUBSTRING.value: lambda data, start, end=None: data[start:end],
    FuncName.LEFT.value: lambda data, count: str[:-count],
    FuncName.RIGHT.value: lambda data, count: str[count:],
    FuncName.REVERSE.value: reversed,
    FuncName.TRIM.value: str.strip,
    FuncName.LTRIM.value: str.lstrip,
    FuncName.RTRIM.value: str.rstrip,
    FuncName.LPAD.value: str.ljust,
    FuncName.RPAD.value: str.rjust,
    FuncName.UPPER.value: str.upper,
    FuncName.LOWER.value: str.lower,
    FuncName.UCASE.value: str.upper,
    FuncName.LCASE.value: str.lower,
    FuncName.LOCATE.value: lambda val, data, *args: str.find(
        data, val, *args
    ) + 1,
    FuncName.INSTR.value: lambda val, data, *args: str.find(
        data, val, *args
    ) + 1,
    FuncName.RAND.value: random,
    FuncName.ROUND.value: round,
    FuncName.MD5.value: lambda data: md5(data).digest(),
    FuncName.NOW.value: lambda: datetime.now().strftime(DATETIMEFORMAT),
    FuncName.SEC_TO_TIME.value: lambda date: datetime.fromtimestamp(date).strftime(DATETIMEFORMAT),
    FuncName.DATEDIFF.value: lambda date1, date2: datetime.strpformat(date1, DATETIMEFORMAT) - datetime.strpformat(date2, DATETIMEFORMAT),
    FuncName.MONTH.value: lambda date=None: datetime.strpformat(date, DATETIMEFORMAT).month,
    FuncName.YEAR.value: lambda date=None: datetime.strpformat(date, DATETIMEFORMAT).year,
}


def innerjoin(litems, ritems):

    return [item for item in litems if item in ritems]


def leftjoin(litems, ritems):

    return litems


def leftexjoin(litems, ritems):

    return [item for item in litems if item not in ritems]


def rightjoin(litems, ritems):

    return ritems


def rightexjoin(litems, ritems):

    return [item for item in ritems if item not in litems]


def fulljoin(litems, ritems):
    """Apply full join on litems and rtimes.

    :param list litems:
    :param list ritmes:
    :return: new list of items.
    :rtype: list"""

    return litems + [item for item in ritems if item not in litems]


def fullexjoin(litems, ritems):

    return leftexjoin(litems, ritems) + rightexjoin(litems, ritems)


def crossjoin(litems, ritems):

    return [(litem, ritem) for litem in litems for ritem in ritems]


def selfjoin(litems, ritems):

    return crossjoin(litems, litems)


def naturaljoin(litems, ritems):

    result = []

    for litem in litems:

        for ritem in ritems:

            issame = False

            for key in ritem:

                if key in litem:
                    if litem[key] == ritem[key]:
                        issame = True

                    else:
                        break

            else:
                if issame:
                    item = deepcopy(litem)
                    item.update(deepcopy(ritem))
                    result.append(item)

    return result

def unionjoin(litems, ritems):

    return litems + ritems


_JOINBYNAME = {
    Join.INNER.name: innerjoin,
    Join.LEFT.name: leftjoin,
    Join.LEFTEX.name: leftexjoin,
    Join.RIGHT.name: rightjoin,
    Join.RIGHTEX.name: rightexjoin,
    Join.FULL.name: fulljoin,
    Join.FULLEX.name: fullexjoin,
    Join.CROSS.name: crossjoin,
    Join.SELF.name: selfjoin,
    Join.NATURAL.name: naturaljoin,
    Join.UNION.name: unionjoin
}


def applyjoin(litems, ritems, join=Join.FULL):

    if isinstance(join, Join):
        join = join.name

    func = _JOINBYNAME[join]

    return func(litems, ritems)
