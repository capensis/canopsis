#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2016 "Capensis" [http://www.capensis.com]
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


class Entity:
    def __init__(self, id_, name, impacts, depends, enable_history, measurements, enabled, infos, type_):
        self.id_ = id_
        self.name = name
        self.impacts = impacts
        self.depends = depends
        self.enable_history = enable_history
        self.measurements = measurements
        self.enabled = enabled
        self.infos = infos
        self.type = type_

    def to_dict(self):
        return {
            'name': self.name,
            'impact': self.impacts,
            'depends': self.depends,
            'enable_history': self.enable_history,
            'measurements': self.measurements,
            'enabled': self.enabled,
            'infos': self.infos,
            'type': self.type,
            '_id': self.id_,
            'entity_id': self.id_
        }
