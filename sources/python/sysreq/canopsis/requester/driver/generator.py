
"""Driver generator module."""

__all__ = [
    'FunctionalDriver',
    'func2crodprocessing', 'obj2driver', 'DriverAnnotation',
    'CreateAnnotation', 'ReadAnnotation', 'UpdateAnnotation', 'DeleteAnnotation'
]

from canopsis.annotation import Annotation
from canopsis.schema import data2schema, data2schemacls, Schema
from canopsis.schema.lang.python import FunctionSchema

from six import string_types

from .base import Driver
from .py import read
from ..request.crud.base import CRUD
from ..request.crud.create import Create
from ..request.crud.read import Read
from ..request.crud.update import Update
from ..request.crud.delete import Delete


class FunctionalDriver(Driver):
    """Driver with fine grained implementation of crud functions.

    This driver uses at most five set of functions for five respective crud
    types.

    The processing execute the right function for all request crud objects.

    Functions must takes in parameters a 'crud' object, a 'request' object and
    kwargs for specific driver uses (like explain for example).

    Functions must return a request."""

    def __init__(
            self, creates=None, reads=None, updates=None, deletes=None,
            functions=None, *args, **kwargs
    ):
        """
        :param list creates: creation functions. Take in parameters a request,
            a crud operation and specific kwargs. Return created item
            with additional fields such as id.
        :param list reads: reading functions. Take in parameters a request,
            a crud operation and specific kwargs. Return read items.
        :param list updates: updating functions. Take in parameters a request,
            a crud operation and specific kwargs. Return updated items.
        :param list deletes: deletion functions. Take in parameters a request,
            a crud operation and specific kwargs. Return deleted items.
        :param dict functions: function to process for specific query functions.
        """

        super(FunctionalDriver, self).__init__(*args, **kwargs)

        self.creates = [] if creates is None else creates
        self.reads = [] if reads is None else reads
        self.updates = [] if updates is None else updates
        self.deletes = [] if deletes is None else deletes
        self.functions = {} if functions is None else functions

    def _process(self, request, crud, **kwargs):

        result = request

        crudname = type(crud).__name__.lower()

        funcs = getattr(self, '{0}s'.format(crudname))

        if funcs:
            for func in funcs:
                result = func(crud=crud, request=result, **kwargs)

        else:
            raise NotImplementedError(
                'No implementation found for {0}'.format(crudname)
            )

        return result


def func2crudprocessing(func=None):

    if func is not None and not isinstance(func, Schema):
        func = data2schema(func)

    def _processing(crud, request, _func=func, **kwargs):

        funckwargs = {}
        funcvarargs = []

        if isinstance(crud, (Create, Update)):
            funckwargs.update(crud.values)
            # todo : specific _func args

        for param in _func.params:

            if param.name in request.ctx:

                funckwargs[param.name] = request.ctx[param.name]

        funcresult = list(_func(*funcvarargs, **funckwargs))

        if isinstance(crud, Read):
            read(read=crud, items=funcresult)

        request.ctx[crud] = funcresult

        return request

    return _processing


def obj2driver(
        obj, name=None, creates=None, reads=None, updates=None, deletes=None,
):
    """Convert an object to a driver.

    :param str name: driver name. obj type name by default.
    :param list creates: create function names to retrieve from the obj.
    :param list reads: read function names to retrieve from the obj.
    :param list updates: update function names to retrieve from the obj.
    :param list deletes: delete function names to retrieve from the obj.
    :rtype: FunctionalDriver
    """

    fname = type(obj).__name__ if name is None else name

    # build obj schema if necessary
    if not isinstance(obj, Schema):
        fobj = data2schema(obj)

    else:
        fobj = obj

    fcreates = []
    freads = []
    fupdates = []
    fdeletes = []

    _locals = locals()
    # start to get information from annotations
    crudannotations = set(
        _CRUDAnnotation.get_annotations(obj) +
        _CRUDAnnotation.get_annotations(fobj)
    )
    for crudannotation in crudannotations:
        for target in crudannotation.targets:
            targetname = targetname
            fobjtarget = getattr(fobj, targetname)
            _locals['f{0}s'.format(crudannotation.name)].append(fobjtarget)

    # then parse function parameters
    for crudname in (crudname.lower() for crudname in  CRUD.__members__):

        cruds = _locals['{0}s'.format(crudname)]

        for crud in cruds:

            crudfunc = getattr(fobj, crud)

            fcrudfunc = func2crudprocessing(crudfunc)

            _locals['f{0}s'.format(crudname)].append(fcrudfunc)

    return FunctionalDriver(
        name=fname,
        creates=fcreates,
        reads=freads,
        updates=fupdates,
        deletes=fdeletes
    )


class DriverAnnotation(Annotation):
    """Generate a deriver from class content."""

    def __init__(
            self,
            name=None, creates=None, reads=None, updates=None, deletes=None,
            *args, **kwargs
    ):
        """
        :param str name: driver name to generate. Default target type name.
        :param list creates: create function names. Related functions must
            return created objects with additional fields such as the id.
        :param list reads: read function names. Related functions must return
            list of read objects.
        :param list updates: update function names. Related functions must
            return updated data.
        :param list deletes: delete function names. Related functions must
            return deleted objects."""

        super(DriverAnnotation, self).__init__(*args, **kwargs)

        self.name = name
        self.creates = creates
        self.reads = reads
        self.updates = updates
        self.deletes = deletes

    def getdriver(self, obj):
        """Get a driver corresponding to input target instance related
        to this attributes.

        :param obj: instance to transform to a functional driver.
        :rtype: FunctionalDriver"""

        kwargs = {'name': self.name}

        for crud in (crud.lower() for crud in CRUD.__members__):

            fcrud = '{0}s'.format(crud)

            funcnames = getattr(self, fcrud)

            for funcname in funcnames:
                func = getattr(obj, funcname)

                kwargs.setdefault(fcrud, []).append(func)

        return obj2driver(**kwargs)


class _CRUDAnnotation(Annotation):

    def __init__(self, crud, *args, **kwargs):
        """
        :param str crud: crud name.
        """

        super(_CRUDAnnotation, self).__init__(self)

        self.crud = crud.name if isinstance(crud, CRUDElement) else crud


class CreateAnnotation(_CRUDAnnotation):

    def __init__(self, crud=CRUD.READ, *args, **kwargs):

        super(CreateAnnotation, self).__init__(crud=crud, *args, **kwargs)


class ReadAnnotation(_CRUDAnnotation):

    def __init__(self, crud=CRUD.READ, *args, **kwargs):

        super(ReadAnnotation, self).__init__(crud=crud, *args, **kwargs)


class UpdateAnnotation(_CRUDAnnotation):

    def __init__(self, crud=CRUD.UPDATE, *args, **kwargs):

        super(UpdateAnnotation, self).__init__(crud=crud, *args, **kwargs)


class DeleteAnnotation(_CRUDAnnotation):

    def __init__(self, crud=CRUD.DELETE, *args, **kwargs):

        super(DeleteAnnotation, self).__init__(crud=crud, *args, **kwargs)
