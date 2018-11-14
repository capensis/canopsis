#!/usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2018 "Capensis" [http://www.capensis.com]
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

"""
Enumerations for event filter rules.
"""

from __future__ import unicode_literals

from canopsis.common.enumerations import FastEnum


class RuleField(FastEnum):
    """
    The RuleField enumeration defines the names of the fields of a rule.
    """
    id = '_id'
    type = 'type'
    pattern = 'pattern'
    priority = 'priority'
    enabled = 'enabled'
    external_data = 'external_data'
    actions = 'actions'
    on_success = 'on_success'
    on_failure = 'on_failure'


class RuleType(FastEnum):
    """
    The RuleType enumeration defines the available types of rules.
    """
    break_ = "break"
    drop = "drop"
    enrichment = "enrichment"


class RuleOutcome(FastEnum):
    """
    The RuleOutcome enumeration defines the available outcomes for a rule.
    """
    pass_ = "pass"
    break_ = "break"
    drop = "drop"


ENRICHMENT_FIELDS = set((
    RuleField.external_data,
    RuleField.actions,
    RuleField.on_success,
    RuleField.on_failure))
