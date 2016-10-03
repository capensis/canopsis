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

"""Driver module."""

__all__ = ['Driver']


try:
    from threading import thread

except ImportError:
    from dummy_threading import Thread


class Driver(object):
    """In charge of accessing to data from a request."""

    __DEFAULTEXPLAIN__ = False  #: default explain value.
    __DEFAULTASYNC__ = False  #: default async value.

    name = None  # driver name

    def __init__(self, name=None, *args, **kwargs):

        super(Driver, self).__init__(*args, **kwargs)

        if name is not None:
            self.name = name

    def process(
            self, request, crud,
            explain=__DEFAULTEXPLAIN__, async=__DEFAULTASYNC__, callback=None,
            **kwargs
    ):
        """Process input request and crud element.

        :param Request request: request to process.
        :param bool explain: give additional information about the request
            execution.
        :param bool async: if True (default False), execute input crud in a
            separated thread.
        :param callable callback: callable function which takes in parameter the
            function result. Commonly used with async equals True.
        :param dict kwargs: additional parameters specific to the driver.
        :return: request.
        :rtype: Request
        """

        def process(
                request=request, crud=crud, explain=explain, callback=callback,
                **kwargs
        ):

            result = self._process(
                request=request, crud=crud, explain=explain, **kwargs
            )

            if callback is not None:
                callback(result)

            return result

        if async:
            Thread(target=process).start()

        else:
            return process()

    def _process(self, request, crud, explain=False, **kwargs):
        """Generic method to override in order to crud input data related to
        request and kwargs.

        :param Request request: request to process.
        :param bool explain: give additional information about the request
            execution.
        :param dict kwargs: additional parameters specific to the driver.
        :return: request.
        :rtype: Request
        """

        raise NotImplementedError()

    def __repr__(self):

        return 'Driver({0})'.format(self.name)
