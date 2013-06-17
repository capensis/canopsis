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
		var raw_vars = this.extractVariables(this.text);

		if (raw_vars.length != 0) {
			this.perfdataMetricList = {};

			var vars = this.cleanVars(raw_vars);

			log.debug('Rebuild vars form user template', this.logAuthor);
			Ext.Object.each(vars, function(key, value){
				log.debug(' + ' + key, this.logAuthor);
				var var_name = value[0];
				if ((var_name == "perfdata" || var_name == "perf_data") && value.length == 3){

					var_name = "perf_data";
					var metric = value[1];
					var attribut = value[2];

					var tpl_name = var_name  + Math.ceil(Math.random() * 1000)
	
					this.text = this.text.replace(new RegExp(key), '{' + tpl_name + '}');

					this.perfdataMetricList[metric] = {
						metric: metric,
						attribut: attribut,
						tpl_name: tpl_name
					};
				}
			}, this);

			log.dump(this.perfdataMetricList);
		}


		//Initialisation of ext JS template
		this.myTemplate = new Ext.XTemplate('<div>' + this.text + '</div>');

		//Compilation of template ( to accelerate the render )
		this.myTemplate.compile();
		this.HTML = ''; // contains the html
		this.callParent(arguments); // Initialization globale of the template
	},

	onRefresh: function(data) {
		perf_data = {};

		if (data){
			if (data.perf_data_array){
				log.debug('Parse perf_data_array', this.logAuthor);
				for (var i = 0; i < data.perf_data_array.length; i++)
					perf_data[data.perf_data_array[i]["metric"]] = data.perf_data_array[i];

				log.dump(perf_data);
			}

			if (this.perfdataMetricList){
				log.debug('Parse template perf_data', this.logAuthor);

				Ext.Object.each(this.perfdataMetricList, function(key, value){
					var metric = key;
					var attribut = value.attribut;
					var tpl_name = value.tpl_name;

					log.debug(' + ' + metric + '(' + tpl_name + ')' + ': ' + attribut, this.logAuthor)

					var perf = perf_data[metric];

					if (perf && perf[attribut] != undefined){
						var unit = perf["unit"];
						var value = perf[attribut];
					
						if (Ext.isNumeric(value) && unit)
							value = rdr_humanreadable_value(value, unit)

						log.debug('   + ' + value, this.logAuthor);

						data[tpl_name] = value;
					}
				});
			}

			try {
				data.timestamp = rdr_tstodate(data.timestamp);
				
				if (this.lastRefresh)
					data.lastRefresh = rdr_tstodate(parseInt(this.lastRefresh / 1000));

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
		log.debug("extractVariables:", this.logAuthor);
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
				while (_string.charAt(end) != '}' && end <= _string.length)
					end = end + 1;

				var_array.push(_string.slice(begin, end + 1));
				_string = _string.slice(end, _string.length);
			}else {
				loop = false;
			}
		}

		log.dump(var_array);
		return var_array;
	},

	// return :  ['{var1:var2}',['var1','var2']]
	cleanVars: function(array) {
		var output = {};
		for (var i = 0; i < array.length; i++)
			output[array[i]] = array[i].slice(1, -1).split(':')

		return output;
	}


});
