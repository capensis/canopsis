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
Ext.define('canopsis.model.Event', {
	extend: 'Ext.data.Model',

	fields: [
		{name: '_id'},
		{name: 'connector'},
		{name: 'connector_name'},
		{name: 'event_type'},
		{name: 'source_type'},
		{name: 'component'},
		{name: 'resource'},
		{name: 'timestamp'},
		{name: 'state'},
		{name: 'state_type'},
		{name: 'output'},
		{name: 'long_output'},
		{name: 'perf_data'},
		{name: 'perf_data_array'},
		{name: 'tags'},
		{name: 'id'},

		{name: 'event_id'},
		{name: 'derogation_name'},
		{name: 'derogation_description'},
		{name: 'ack'},

		{name: 'ticket'},
		
		{name: 'ref_rk'}
	]
});
