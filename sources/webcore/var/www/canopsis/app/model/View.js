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
Ext.define('canopsis.model.View', {
	extend: 'Ext.data.Model',

	fields: [
		{name: '_id'},
		{name: 'id', mapping: '_id'},
		{name: 'leaf'},
		{name: 'crecord_name'},
		{name: 'nodeId'},
		{name: 'items'},
		{name: 'refreshInterval'},
		{name: 'template', defaultValue: false},
		{name: 'reporting', defaultValue: false},

		{name: 'enable'},
		{name: 'view_options', defaultValue: {orientation: 'portrait', pageSize: 'A4'}},

		{name: 'aaa_access_group'},
		{name: 'aaa_access_other'},
		{name: 'aaa_access_owner'},
		{name: 'aaa_admin_group'},
		{name: 'aaa_group'},
		{name: 'aaa_owner'},
		{name: 'crecord_write_time'},
		{name: 'crecord_creation_time'}
	]
});
