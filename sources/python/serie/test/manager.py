#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from unittest import TestCase, main

from canopsis.middleware.core import Middleware
from canopsis.perfdata.manager import PerfData
from canopsis.context.manager import Context
from canopsis.serie.manager import Serie

##TODO4-01-2017
#class BaseTestSerieManager(TestCase):
#    def setUp(self):
#        self.ctx_storage = Middleware.get_middleware_by_uri(
#            'storage-composite-test_context://',
#            path=Context.DEFAULT_CONTEXT
#        )
#        self.ctx_manager = Context()
#        self.ctx_manager[Context.CTX_STORAGE] = self.ctx_storage
#
#        self.perfdata_storage = Middleware.get_middleware_by_uri(
#            'storage-timed-test_perfmeta://'
#        )
#        self.perf_manager = PerfData()
#        self.perf_manager[PerfData.PERFDATA_STORAGE] = self.perfdata_storage
#        self.perf_manager[PerfData.CONTEXT_MANAGER] = self.ctx_manager
#
#        self.serie_storage = Middleware.get_middleware_by_uri(
#            'storage-default-test_serie://'
#        )
#        self.serie_manager = Serie()
#        self.serie_manager[Serie.SERIE_STORAGE] = self.serie_storage
#        self.serie_manager[Serie.CONTEXT_MANAGER] = self.ctx_manager
#        self.serie_manager[Serie.PERFDATA_MANAGER] = self.perf_manager
#
#
#class TestSerieManager(BaseTestSerieManager):
#    def test_get_metrics(self):
#        raise NotImplementedError()
#
#    def test_get_metrics_subset(self):
#        raise NotImplementedError()
#
#    def test_get_perfdata(self):
#        raise NotImplementedError()
#
#    def test_subset_perfdata_superposed(self):
#        raise NotImplementedError()
#
#    def test_aggregation_round_time(self):
#        raise NotImplementedError()
#
#    def test_aggregation_no_round_time(self):
#        raise NotImplementedError()
#
#    def test_consolidation_round_time(self):
#        raise NotImplementedError()
#
#    def test_consolidation_no_round_time(self):
#        raise NotImplementedError()
#
#    def test_calculate(self):
#        raise NotImplementedError()
#
#    def test_get_series(self):
#        raise NotImplementedError()
#

if __name__ == '__main__':
    main()
