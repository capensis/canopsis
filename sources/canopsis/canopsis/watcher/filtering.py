#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
Filter watchers.
"""

from __future__ import unicode_literals

import copy


class WatcherFilter(object):
    """
    Filter out active_pb_some and active_pb_all.
    """

    def __init__(self):
        self._all = None
        self._some = None
        self._watcher = None
        self._include_types = set()
        self._exclude_types = set()

    def all(self):
        """
        :returns: True, False or None
        """
        return self._all

    def some(self):
        """
        :returns: True, False or None
        """
        return self._some

    def watcher(self):
        """
        :returns: True, False or None if the watcher has an active pbehavior.
        """
        return self._watcher

    def include_types(self):
        """
        :returns: set of pbehavior types to filter on.
        :rtype: set[str]
        """
        return self._include_types

    def exclude_types(self):
        """
        oposit of include_types
        """
        return self._exclude_types

    @staticmethod
    def to_bool(value):
        """
        :rtype: bool
        """
        return value in ['1', 1, "true", "True", True]

    def _filter_dict(self, dictdoc):
        cdoc = copy.deepcopy(dictdoc)
        for k, v in dictdoc.items():
            if k == 'active_pb_all':
                self._all = self.to_bool(v)
                del cdoc[k]

            elif k == 'active_pb_some':
                self._some = self.to_bool(v)
                del cdoc[k]

            elif k in ['active_pb_type', 'active_pb_include_type']:
                self._include_types.add(str(v).strip().lower())
                del cdoc[k]

            elif k in ['active_pb_exclude_type']:
                self._exclude_types.add(str(v).strip().lower())
                del cdoc[k]

            elif k == 'active_pb_watcher':
                self._watcher = self.to_bool(v)
                del cdoc[k]

            else:
                nv = self._filter(v)
                if nv is not None or v is None:
                    cdoc[k] = nv
                else:
                    del cdoc[k]

        return cdoc

    def _filter_list(self, listdoc):
        cdoc = copy.deepcopy(listdoc)
        j = 0
        for i, item in enumerate(listdoc):
            v = self._filter(item)
            if v is not None:
                cdoc[j] = v
                j += 1
            else:
                del cdoc[j]

        return cdoc

    def _filter(self, doc):
        if isinstance(doc, dict):
            ndoc = self._filter_dict(doc)
            if len(ndoc) == 0 and len(doc) != 0:
                return None
            return ndoc

        elif isinstance(doc, list):
            ndoc = self._filter_list(doc)
            if len(ndoc) == 0 and len(doc) != 0:
                return None
            return ndoc

        return doc

    def _filter_pb_types(self, pb_types):
        """
        :param pb_types set:
        """
        if not pb_types or (not self.exclude_types() and not self.include_types()):
            return True

        for pb_type in pb_types:
            normalized_pb = pb_type.strip().lower()

            if normalized_pb in self.exclude_types():
                return False

            if normalized_pb in self.include_types():
                return True

        if not self.include_types():
            return True

        return False

    def filter(self, doc):
        """
        :rtype: dict
        """
        res = self._filter(doc)
        if res is None:
            return {}
        return res

    def match(self, allstatus, somestatus, watcherstatus, pb_types=None):
        """
        Call WatcherFilter.filter(filter_) first before calling this function.

        :param bool allstatus: watcher has all entities with an active pbehavior
        :param bool somestatus: watcher has some or all entities with an active pbehavior
        :param bool watcherstatus: watcher has a pbehavior or not. If the filter contains a filter on that, it will be
        :param list[str] pb_types: list of pbehavior types to filter on. If None or empty, any types will be ok.
        """
        if not isinstance(allstatus, bool):
            raise TypeError('wrong allstatus value: not a bool')
        if not isinstance(somestatus, bool):
            raise TypeError('wrong somestatus value: not a bool')
        if not isinstance(watcherstatus, bool):
            raise TypeError('wrong watcherstatus value: not a bool')

        if allstatus is True and somestatus is False:
            raise ValueError('allstatus cannot be true if somestatus is false')

        if pb_types is None:
            pb_types = list()

        if not isinstance(pb_types, (list, set)):
            raise TypeError('wrong pb_types value: not a list')

        logic_some_all = False
        logic_watcher = False
        logic_pb_types = self._filter_pb_types(pb_types)

        # check watcher status
        if self.watcher() is None:
            logic_watcher = True

        elif self.watcher() is True and watcherstatus is True:
            logic_watcher = True

        elif self.watcher() is False and watcherstatus is False:
            logic_watcher = True

        # check entities status
        if self.all() is None and self.some() is None:
            logic_some_all = True

        elif self.all() is not None and self.some() is not None:
            # cannot be simplified as only this test.
            if self.all() == allstatus and self.some() == somestatus:
                logic_some_all = True

        elif self.all() is allstatus and self.all() is not None:
            logic_some_all = True

        elif self.some() is somestatus and self.some() is not None:
            logic_some_all = True

        # check if a pbehavior on a watcher can be included into pbehavior on entities
        if self.watcher() is None and watcherstatus is True:
            if self.some() is not False:
                logic_watcher = True
                logic_some_all = True
            if self.some() is False:
                logic_watcher = False
                logic_some_all = False

        return logic_watcher and logic_some_all and logic_pb_types
