
"""Module which specifices a composite of drivers."""

__all__ = ['DriverComposite']

from .base import Driver

from ..request.core import Request
from ..request.expr import Expression, Function, FuncName
from ..request.crud.create import Create
from ..request.crud.read import Read
from ..request.crud.update import Update
from ..request.crud.delete import Delete

from canopsis.schema import data2schema, Schema

try:
    from threading import Thread

except ImportError:
    from dummy_threading import Thread


class DriverComposite(Driver):

    __LAST_DRIVER__ = '__last_driver__'

    def __init__(self, drivers, *args, **kwargs):
        """
        :param list drivers: drivers to use.
        """

        super(DriverComposite, self).__init__(*args, **kwargs)

        self.drivers = {}
        for driver in drivers:
            if not isinstance(driver, Schema):
                driver = data2schema(driver, name=driver.name, _force=True)

            self.drivers[driver.name] = driver

    def _process(self, request, crud, **kwargs):

        result = self._processquery(query=query, ctx=request.ctx, **kwargs)

        if isinstance(crud, (Create, Update)):
            if isinstance(crud.name, Expression):
                ctx = self._processquery(query=crud.name, ctx=ctx).ctx
                names = [crud.name.name]

            else:
                names = [crud.name]

        elif isinstance(crud, (Read, Delete)):
            names = []

            items = crud.select() if isinstance(crud, Read) else crud.names

            for item in items:
                if isinstance(item, Expression):
                    ctx = self._processquery(query=item, ctx=ctx)
                    names.append(item.name)

                else:
                    names.append(item)

        if names:
            driver = self.drivers.get(names[0])

            if driver is None:

                processingdrivers = [
                    driver for driver in self.drivers.values()
                    if names[0] in driver.getschemas()
                ]

            else:
                processingdrivers = [driver]

            request = Request(query=query, ctx=ctx, cruds=crud)

            for processingdriver in processingdrivers:

                pctx = processingdriver.process(request=request)

                ctx.fill(pctx)

        request.ctx = ctx

    def _processquery(self, query, ctx, _lastquery=None, **kwargs):
        """Parse deeply the query from the left to the right and aggregates
        queries which refers to the same system without any intermediary system.

        :param Expression query: query to process.
        :param Context ctx: context execution.
        :param Expression _lastquery: private parameter which save the last
            query to process.
        :return: ctx
        :rtype: ctx
        """

        if _lastquery is None:
            _lastquery = query

        names = query.name.split('.')

        driver = self.drivers.get(names[0])
        lastdriver = ctx.get(DriverComposite.__LAST_DRIVER__, driver)

        if driver is None:
            driver = lastdriver

        else:
            if driver != lastdriver:
                ctx[DriverComposite.__LAST_DRIVER__] = driver

        if isinstance(query, Function):

            isor = query.name == FuncName.OR.value

            for param in query.params:
                if isinstance(param, Expression):

                    pctx = deepcopy(ctx) if isor else ctx

                    self._processquery(
                        query=param, ctx=pctx,
                        _lastquery=query if _lastquery is None else _lastquery,
                        **kwargs
                    )

                    if isor:
                        ctx.fill(pctx)

        if query == _lastquery:

            if lastdriver is None:
                lastdriver = ctx.get(DriverComposite.__LAST_DRIVER__, driver)

            if lastdriver is None:

                processingdrivers = [
                    driver for driver in self.drivers.values()
                    if names[0] in driver.getschemas()
                ]

            else:
                processingdrivers = [lastdriver]

            request = Request(query=query, ctx=ctx)

            for processingdriver in processingdrivers:

                # prepare a read operation for the inner driver
                crudname = query.ctxname  # prepare crud name
                if crudname.startswith(processingdriver.name):
                    crudname = crudname[len(processingdriver.name) + 1:]

                crud = Read(
                    select=[crudname] if crudname else None,
                    alias=query.ctxname  # set alias equals query ctxname
                )

                request.cruds = [crud]

                request = processingdriver.process(request=request, **kwargs)

                ctx = request.ctx

        return ctx

    def __repr__(self):

        return 'CompositeDriver({0}, {1})'.format(self.name, self.drivers)
