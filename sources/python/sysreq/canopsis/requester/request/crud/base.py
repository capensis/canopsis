
"""Base crud module."""

__all__ = ['CRUD']

from ..expr import BaseElement
from ...driver.base import Driver

from enum import IntEnum, unique


@unique
class CRUD(IntEnum):

    CREATE = 1
    READ = 2
    UPDATE = 3
    DELETE = 4


class CRUDElement(BaseElement):
    """Base crud operation.

    Can be associated to a request."""

    __slots__ = ['request', 'result'] + BaseElement.__slots__

    def __init__(self, request=None, result=None, *args, **kwargs):
        """
        :param canopsis.requester.Request request:
        :param result: result of this crud processing.
        """

        super(CRUDElement, self).__init__(*args, **kwargs)

        self.request = request
        self.result = result

    def __call__(
            self,
            explain=Driver.__DEFAULTEXPLAIN__, async=Driver.__DEFAULTASYNC__,
            callback=None, **kwargs
    ):
        """Execute this CRUD element.

        :param bool async: if False (default), result is this element result,
            otherwise, process the function in an isolated thread.
        :param callback: callback callable which takes in parameter the CRUD
            element result.
        :param dict kwargs: driver kwargs.
        :return: this execution result if async is False."""

        if self.request is None:
            raise RuntimeError(
                'Impossible to execute this without associate it to a request.'
            )

        else:
            return self.request.processcrud(
                crud=self, explain=explain, async=async, callback=callback,
                *kwargs
            )
