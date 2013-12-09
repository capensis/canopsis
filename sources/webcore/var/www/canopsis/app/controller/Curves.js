//need:app/lib/controller/cgrid.js,app/view/Curves/Grid.js,app/view/Curves/Form.js
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
Ext.define('canopsis.controller.Curves', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Curves.Grid', 'Curves.Form'],
	stores: ['Curves'],
	models: ['Curve'],

	logAuthor: '[controller][Curves]',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'CurvesForm';
		this.listXtype = 'CurvesGrid';

		this.modelId = 'Curve';

		this.callParent(arguments);

		global.curvesCtrl = this;

		Ext.require('Ext.menu.ColorPicker');
    },

	preSave: function(record) {
		record.data['_id'] = $.encoding.digests.hexSha1Str(record.data['metric']);
		record.data['crecord_name'] = record.data['metric'];
		return record;
	},

	beforeload_EditForm: function(form) {
		var field = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=metric]')[0];

		if(field) {
			field.hide();
		}
	},

	getRenderInfo: function(metric) {
		if(metric) {
			var _id   = $.encoding.digests.hexSha1Str(metric);
			var store = Ext.data.StoreManager.lookup('Curves');
			var info  = store.getById(_id);

			if(info) {
				return info;
			}
		}
	},

	getRenderColors: function(metric, index) {
		if(!index) {
			index = 0;
		}

		var line_color = global.default_colors[index];
		var area_color = line_color;
		var area_opacity = 75;

		var info = this.getRenderInfo(metric);

		if(info) {
			line_color = info.get('line_color');
			area_color = info.get('area_color');
			area_opacity = info.get('area_opacity');
		}

		if(line_color) {
			line_color = '#' + line_color;
		}
		else {
			line_color = undefined;
		}

		if(area_color) {
			area_color = '#' + area_color;
		}
		else {
			area_color = undefined;
		}

		return [line_color, area_color, area_opacity];
	}
});
