/*
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
Ext.define('canopsis.model.Derogation', {
	extend: 'Ext.data.Model',
	idProperty: '_id',
	fields: [
		{name: '_id'},

		{name: 'crecord_name'},
		{name: 'description'},

		{name: 'forTs', defaultValue: undefined},
		{name: 'ts_unit', defaultValue: undefined},
		{name: 'ts_window', defaultValue: undefined},

		{name: 'actions', defaultValue: []},
		{name: 'conditions', defaultValue: []},
		{name: 'time_conditions', defaultValue: []},

		{name: 'ids', defaultValue: []},
		{name: 'name', defaultValue: undefined},

		{name: 'aaa_access_group', defaultValue: ['r', 'w'] },
		{name: 'aaa_access_owner', defaultValue: ['r', 'w'] },
		{name: 'aaa_access_other', defaultValue: ['r'] },
		{name: 'aaa_admin_group'},
		{name: 'aaa_group', defaultValue: 'CPS_derogation_admin' },
		{name: 'aaa_owner'},
		{name: 'crecord_type'},

		{name: 'selector_name'},

		{name: 'enable', defaultValue: true},
		{name: 'active', defaultValue: false},

		{name: 'tags', defaultValue: [], type: 'array'}
	]
});
