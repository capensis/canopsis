/*
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
*/
Ext.define('canopsis.model.Rule', {
	extend: 'Ext.data.Model',
	fields: [
		{name: '_id'},
		{name: 'id', mapping: '_id'},
		{name: 'crecord_type'},
		{name: 'crecord_name'},
		{name: 'priority', defaultValue: 99, type: 'int'},
		{name: 'mfilter'},
		{name: 'action', defaultValue: 'pass'},
		{name: 'enable', defaultValue: true},

		{name: 'aaa_access_group', defaultValue: undefined},
		{name: 'aaa_access_other', defaultValue: undefined},
		{name: 'aaa_access_owner', defaultValue: ['r', 'w']},
		{name: 'aaa_admin_group', defaultValue: ['r']},
		{name: 'aaa_group', defaultValue: 'group.CPS_rule_admin'},
		{name: 'aaa_owner', defaultValue: undefined}
	]
});
