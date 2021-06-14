//need:app/view/Tabs/Content.js
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
Ext.define('canopsis.view.ReportView', {
	extend: 'canopsis.view.Tabs.Content',

	setContent: function() {
		// Setting globals options
		var items = this.view.items;

		if(items.length === 1) {
			this.set_one_item(items);
		}
		else {
			this.set_many_items(items);
		}

		//recalculate document height
		actual_height = $('#' + this.jqgridable.id).height();
		document.body.style.height = actual_height + 100;

		$('#' + this.jqgridable.id).height(actual_height + 100);
		this.doLayout();

		var task = new Ext.util.DelayedTask(function() {
			window.status = 'ready';
		});

		task.delay(3000);
	},

	set_one_item: function(items) {
		log.debug(' + Use full mode ...', this.logAuthor);
		this.layout = 'fit';

		//check if standart item view or jqgridable item
		var item = items[0];

		if(item.data) {
			item = item.data;
		}

		log.debug('   + Add: ' + item.xtype, this.logAuthor);

		item['width'] = '100%';
		item['title'] = '';
		item['mytab'] = this;

		item['reportMode'] = true;
		item['exportMode'] = true;

		item.export_from = export_from;
		item.export_to = export_to;

		this.add(item);
	},

	set_many_items: function(items) {
		//setting jqgridable
		this.jqgridable = Ext.create('canopsis.view.Tabs.JqGridableViewer', {
			view_widgets: items
		});

		this.add(this.jqgridable);

		this.jqgridable._load(items);

		var widget_list = this.jqgridable.get_ext_widget_list();

		for(var i = 0; i < widget_list.length; i++) {
			var item = items[i].data;

			item['reportMode'] = true;
			item['exportMode'] = true;

			item['height'] = widget_list[i].height;

			item.export_from = export_from;
			item.export_to = export_to;

			widget_list[i].add(item);
		}
	}
});
