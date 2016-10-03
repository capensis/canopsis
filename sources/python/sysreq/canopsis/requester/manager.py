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

"""Manager module."""

__all__ = ['RequestManager']

from .request import Request

from .request.crud.base import CRUD
from .request.crud.create import Create
from .request.crud.read import Read
from .request.crud.update import Update
from .request.crud.delete import Delete
from .driver.base import Driver


class RequestManager(Driver):
    """Driver in charge of instanciate requests.

    All instanciated requests will be associated to this manager."""

    def __init__(self, driver, requests=None, *args, **kwargs):

        super(RequestManager, self).__init__(*args, **kwargs)

        self.driver = driver
        self._requests = [] if requests is None else requests
        self._rindex = 0

    def request(request=None):

        result = None

        if request is not None:
            result = Request(
                manager=self,
                ctx=request.ctx, query=request.query, cruds=request.cruds
            )

        else:
            result = Request(manager=self)

        self._requests.append(result)

        return result

    def process(self, request, **kwargs):

        result = self.driver.process(request, **kwargs)

        return result

    @property
    def processed(self):

        return self._requests[:self._rindex]

    @property
    def waiting(self):

        return self._requests[self._rindex:]

    def clean(self, processed=True, waiting=True):

        if processed and waiting:
            self._requests = []

        elif processed:
            self._requests = self.waiting

        elif waiting:
            self._requests = self.processed

        return self

    def cruds2request(self, *cruds, **kwargs):
        """Transform input cruds to an input request.

        :param tuple cruds: list of crud operations to exe.
        :param dict kwargs: request parameters such as ctx, query.
        :rtype: Request
        :return: created request"""

        result = Request(driver=self, cruds=cruds, **kwargs)

        return result

    def create(self, name, **values):
        """Create an element with input kwargs values."""

        create = Create(name=name, params=[values])

        result = self.crud2request(create)

        return result

    def read(self, *names, **kwargs):
        """Read data.

        :param tuple names: data to retrieve. Default all.
        :param dict kwargs: Read parameters (limit, skip, etc.).
        """

        read = Read(name=names, **kwargs)

        result = self.crud2request(read)

        return self.ctx[read.name]

    def update(self, name, **values):

        update = Update(name=name, params=[values])

        result = self.crud2request(update)

        return self.ctx[update.name]

    def delete(self, *names):

        delete = Delete(name=names)

        result = self.crud2request(delete)

        return self.ctx[delete.name]

    def exe(self, name, *params):

        crud = CRUD(name, params=params)

        result = self.crud2request(crud)

        return self.ctx[crud.name]

    def __getitem__(self, key):

        return self.read(key)

    def __setitem__(self, key, item):

        return self.update(name=key, **item)

    def __delitem__(self, key):

        return self.delete(key)
