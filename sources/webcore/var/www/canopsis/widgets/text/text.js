/*
#--------------------------------
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
# ---------------------------------
*/
Ext.define('widgets.text.text' , {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.text',

	//templateVars : undefined,
	perfdataMetricList: undefined,
	logAuthor: '[textWidget]',
	specialCharRegex: /\//g,

	initComponent: function() {
		//get special values by parsing
		var raw_var = this.extractVariables(this.text);

		if (raw_var.length != 0) {
			//this.templateVars = []
			this.perfdataMetricList = [];

			var extracted_vars = this.cleanVars(raw_var);

			//replacing var in this.text by new var name
			for (var i=0; i < extracted_vars.length; i++) {
				var one_var = extracted_vars[i];
				var metric_exploded_name = one_var[1];

				//get rid of perfdata
				if (metric_exploded_name[0] == 'perfdata')
					metric_exploded_name = metric_exploded_name.slice(1, metric_exploded_name.length);

				var new_var_name = '{' + metric_exploded_name.join('') + '}';
				new_var_name = new_var_name.replace(this.specialCharRegex, '');
				this.text = this.text.replace(new RegExp(one_var[0]), new_var_name);

				//keep track of metrics : {load: value, load: unit }
				try {
					this.perfdataMetricList.push([metric_exploded_name[0], metric_exploded_name[1]]);
				}catch (err) {
					log.debug('No attribut specified for var ' + metric_exploded_name[0], this.logAuthor);
				}
			}
		}


		//Initialisation of ext JS template
		this.myTemplate = new Ext.XTemplate('<div>' + this.text + '</div>');

		//Compilation of template ( to accelerate the render )
		this.myTemplate.compile();
		this.HTML = ''; // contains the html
		this.callParent(arguments); // Initialization globale of the template
	},

	onRefresh: function(data) {
		if (data) {
			if (data.perf_data_array && data.perf_data_array.length) {
				if (this.perfdataMetricList && this.perfdataMetricList.length != 0) {
					//loop on var in required perfdata
					for (var i=0; i < this.perfdataMetricList.length; i++) {
						var metric = this.perfdataMetricList[i][0];
						var attribut = this.perfdataMetricList[i][1];

						log.debug('Metric searched is ' + metric, this.logAuthor);
						log.debug('Attribut is ' + attribut, this.logAuthor);

						//search the right metric
						if (metric != undefined && metric != null) {
							for (var j=0; j < data.perf_data_array.length; j++) {
								if (data.perf_data_array[j].metric == metric) {
									log.debug('  + ' + attribut + '  found', this.logAuthor);
									var attributName = this.perfdataMetricList[i].join('');
									attributName = attributName.replace(this.specialCharRegex, '');
									try {
										var value = data.perf_data_array[j][attribut];
										if (value != null && value != undefined)
											if (Ext.isNumeric(value))
												data[attributName] = rdr_humanreadable_value(value);
											else
												data[attributName] = value;
									}catch (err) {
										log.debug('metric : ' + metric + ' have no attribut ' + attribut);
									}
									break;
								}
							}
						}
					}
				}
			}

			try {
				data.timestamp = rdr_tstodate(data.timestamp);
				this.HTML = this.myTemplate.apply(data);
			}catch (err) {
				this.HTML = _('The model widget template is not supported, check if your variables use the correct template.');
			}
		}
		else
		{
			//otherwise we put the text contained in the field
			this.HTML = this.text;
		}
		this.setHtml(this.HTML);
	},

	getNodeInfo: function() {
		//we override the function : if there is'nt any nodeId specified we call the onRefresh function
		if (! this.nodeId)
		{
			this.onRefresh(false);
		}else {
			Ext.Ajax.request({
				url: this.baseUrl,
				scope: this,
				method: 'GET',
				params: {_id: this.nodeId},
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);
					if (this.nodeId.length > 1)
						data = data.data;
					else
						data = data.data[0];
					this._onRefresh(data);
				},
				failure: function(result, request) {
					log.error('Impossible to get Node informations, Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});


		}
		//we call the parent which is applied when there is a nodeId specified.
		//this.callParent(arguments);
	},

	extractVariables: function(text) {
		//search specific value
		var loop = true;
		var _string = text;
		var var_array = [];
		while (loop) {
			//search for val
			var begin = _string.search(/{(.+:)+.+}/);
			if (begin != -1) {
				//search end of val
				var end = begin;
				while (_string[end] != '}' && end <= _string.length)
					end = end + 1;

				var_array.push(_string.slice(begin, end + 1));
				_string = _string.slice(end, _string.length);
			}else {
				loop = false;
			}
		}
		return var_array;
	},

	// return :  ['{var1:var2}',['var1','var2']]
	cleanVars: function(array) {
		var output = [];
		for (var i=0; i < array.length; i++)
			output.push([
					array[i],
					array[i].slice(1, -1).split(':')
				]);
		return output;
	}


});
