Ext.define('canopsis.view.ReportView', {
	extend: 'canopsis.view.Tabs.Content',


	setContent: function() {
		//----------------Setting globals options----------------
		var items = this.view.items;

		if (items.length == 1) {
			this.set_one_item(items);
		} else {
			this.set_many_items(items);
		}

		/*
		for(var i= 0; i < items.length; i++) {

			log.debug(' - Item '+i+':', this.logAuthor)
			var item = items[i]
			if(item.data){
				item = item.data
			}

			log.debug('   + Add: '+item.xtype, this.logAuthor)

			item['mytab'] = this

			item['width'] = '100%'
			item['border'] = false

			item['reportMode'] = true;
			item['exportMode'] = true;

			//log.debug('start timestamp is : ' + export_from, this.logAuthor)
			//log.debug('stop timestamp is : ' + export_to, this.logAuthor)

			item.export_from = export_from
			item.export_to = export_to

			//Set default options
			if (! item.nodeId) { item.nodeId=nodeId}
			if (! item.refreshInterval) { item.refreshInterval=refreshInterval}
			if (! item.rowHeight) { item.height=rowHeight }else{ item.height=item.rowHeight }
			if (item.title){ item.border = true }


			this.add(item);
			//log.debug(item);
		}
		* */

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
		if (item.data) {
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

		//Set default options
		//if (! item.nodeId) { item.nodeId=nodeId}
		//if (! item.refreshInterval) { item.refreshInterval=refreshInterval}

		this.add(item);
	},

	set_many_items: function(items) {
		//setting jqgridable
		this.jqgridable = Ext.create('canopsis.view.Tabs.JqGridableViewer', {view_widgets: items});
		this.add(this.jqgridable);

		this.jqgridable._load(items);

		var widget_list = this.jqgridable.get_ext_widget_list();

		//here processing on data

		//globalNodeId
		//refresh interval

		//////////////////////////

		for (var i in widget_list) {
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
