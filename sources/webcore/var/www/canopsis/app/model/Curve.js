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
Ext.define('canopsis.model.Curve', {
	extend: 'Ext.data.Model',
	idProperty: '_id',
	fields: [
		{name: '_id'},
		{name: 'metric'},
		{name: 'crecord_name'},

		{name: 'line_color',	defaultValue: undefined },
		{name: 'area_color',	defaultValue: undefined },
		{name: 'area_opacity',	defaultValue: undefined, type: 'int'},
		{name: 'invert',	defaultValue: false, type: 'boolean' },
		{name: 'zIndex',	defaultValue: 0, type: 'int'},
		{name: 'dashStyle',	defaultValue: 'Solid' },
		{name: 'label',	defaultValue: undefined },

		{name: 'aaa_access_owner', defaultValue: ['r', 'w']},
		{name: 'aaa_access_group', defaultValue: ['r', 'w'] },
		{name: 'aaa_access_other', defaultValue: ['r'] },
		{name: 'aaa_admin_group', defaultValue: 'group.CPS_curve_admin'},
		{name: 'aaa_group', defaultValue: 'group.CPS_curve_admin' },
		{name: 'aaa_owner'}
	]
});
