#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from canopsis.old.record import Record


class Statemap(Record):
    def __init__(self, statemap=None, *args, **kwargs):
        self.type = "statemap"

        if isinstance(statemap, list):
            self.statemap = statemap

        super(Statemap, self).__init__(_type=self.type, *args, **kwargs)

    def dump(self):
        self.data['statemap'] = self.statemap

        return Record.dump(self)

    def load(self, dump):
        Record.load(self, dump)

        self.statemap = self.data['statemap']

    def get_mapped_state(self, state):
        if state < len(self.statemap):
            return self.statemap[state]
