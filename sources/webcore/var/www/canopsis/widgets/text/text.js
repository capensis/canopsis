//need:app/lib/view/cwidget.js
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
Ext.define('widgets.text.text' , {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.text',

	perfdataMetricList: undefined,
	logAuthor: '[textWidget]',
	specialCharRegex: /\//g,

	aggregate_method: 'LAST',

	baseUrl: '/rest/events',

	useLastRefresh: false,

	initComponent: function() {
		//get special values by parsing
		var raw_vars = this.extractVariables(this.text);

		if(raw_vars.length !== 0) {
			this.perfdataMetricList = {};

			var vars = this.cleanVars(raw_vars);

			log.debug('Rebuild vars form user template', this.logAuthor);

			Ext.Object.each(vars, function(key, value){
				log.debug(' + ' + key, this.logAuthor);
				var var_name = value[0];

				if((var_name === "perfdata" || var_name === "perf_data") && value.length >= 2) {
					var_name = "perf_data";
					var metric = value[1];
					var attribut = value.length>=3? value[2] : 'value';

					log.debug('attribute' + attribut);

					var tpl_name = var_name + Math.ceil(Math.random() * 1000);

					this.text = this.text.replace(new RegExp(key), '{' + tpl_name + '}');

					this.perfdataMetricList[tpl_name] = {
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

		// contains the html
		this.HTML = '';

		// Initialization globale of the template
		this.callParent(arguments);
	},

	/**
	* This method can be called recursively twice if an event is selected.
	* 
	* The first time, input data is event.
	* The second time, input data is event perfstore values.
	*/
	onRefresh: function(data, from, to) {

		if(data && this.perfdataMetricList) { // if an event is selected, add metrics property to perf_data

			var perf_data = {};

			log.debug('The event ' + data + ' is selected ', this.logAuthor);
			log.dump(data);

			log.debug('Parse perf_data_array', this.logAuthor);

			var metrics = [];
			
			Ext.Object.each(this.perfdataMetricList, function(key, value) {
				void(key);
				var metric = value.metric;
				metrics.push(metric);
			});

			// prepare parameters for ajax request
			var filter = {'$and': [{'co': data['component']}, {'re': data['resource']}, {'me': {'$in': metrics}}]};
			var metrics_params = {'filter': Ext.JSON.encode(filter), 'limit': 0, 'show_internals': true};

			Ext.Ajax.request({
				url: '/perfstore',
				scope: this,
				params: metrics_params,
				method: 'GET',
				success: function(response) {
					var _data = Ext.JSON.decode(response.responseText);
					_data = _data.data;
					log.dump(_data);

					log.debug('Get perf_data id from event: ', this.logAuthor);

					var event_ids = [];

					for (var i=0; i<_data.length; i++) {
						var __data = _data[i];
						event_ids.push({id: __data['_id']});
					}

					var _from = parseInt(from / 1000);
					var _to = parseInt(to / 1000);

					// ajax parameters
					var perfdata_params = {
						'nodes': Ext.JSON.encode(event_ids),
						'aggregate_method': this.aggregate_method,
						'aggregate_max_points': 1,
						'timezone': new Date().getTimezoneOffset() * 60,
					};

					Ext.Ajax.request({
						url: '/perfstore/values/' + _from + '/' + _to,
						scope: this,
						params: perfdata_params,
						method: 'POST',
						success: function(response) {
							var _data = Ext.JSON.decode(response.responseText);
							_data = _data.data;
							log.dump(_data);

							log.debug('Get perf_data from ids : ', this.logAuthor);	

							for(var i=0; i<_data.length; i++) {
								var __data = _data[i];
								perf_data[__data.metric] = __data;
								perf_data[__data.metric].value = __data.values[0][1];
								perf_data[__data.metric].unit = __data.bunit;
							}

							log.dump(perf_data);

							Ext.Object.each(this.perfdataMetricList, function(key, value) {
								var metric = value.metric;
								var attribut = value.attribut;
								var tpl_name = key;

								log.debug(' + ' + metric + '(' + tpl_name + ')' + ': ' + attribut, this.logAuthor);

								var perf = perf_data[metric];

								if(perf && perf[attribut] !== undefined) {
									var unit = perf["unit"];
									value = perf[attribut];

									if(Ext.isNumeric(value) && unit) {
										log.dump(this);

										if(this.humanReadable) {
											value = rdr_humanreadable_value(value, unit);
										}
										else if(unit) {
											value = value + ' ' + unit;
										}
									}

									log.debug('   + ' + value, this.logAuthor);

									data[tpl_name] = value;
								}
							});

							this.fillData(data, from, to);
						},
						failure: function(result, request) {
							void(result);
							log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
						}
					});
				},
				failure: function(result, request) {
					void(result);
					log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		} else {
			data = {};
			this.fillData(data, from, to);
		}

	},

	fillData: function(data, from, to) {

		data.timestamp = rdr_tstodate(data.timestamp);

		try {
			if(from) {
				data.from = rdr_tstodate(parseInt(from / 1000));
			}

			if(to) {
				data.to = rdr_tstodate(parseInt(to / 1000));
			}

			this.HTML = this.myTemplate.apply(data);
		} catch (err) {
			log.error(err);
			this.HTML = _('The model widget template is not supported, check if your variables use the correct template.');
		}

		this.setHtml(this.HTML);
	},

	extractVariables: function(text) {
		log.debug("extractVariables:", this.logAuthor);

		//search specific value
		var loop = true;
		var _string = text;
		var var_array = [];

		while(loop) {
			//search for val
			var begin = _string.search(/{(.+:)+.+}/);

			if(begin !== -1) {
				//search end of val
				var end = begin;

				while(_string.charAt(end) !== '}' && end <= _string.length) {
					end = end + 1;
				}

				var_array.push(_string.slice(begin, end + 1));
				_string = _string.slice(end, _string.length);
			}
			else {
				loop = false;
			}
		}

		log.dump(var_array);
		return var_array;
	},

	// return :  ['{var1:var2}',['var1','var2']]
	cleanVars: function(array) {
		var output = {};

		for(var i = 0; i < array.length; i++) {
			output[array[i]] = array[i].slice(1, -1).split(':');
		}

		return output;
	},

	getNodeInfoParams: function(from, to) {
		void(from);
		void(to);
		var result = this.callParent(arguments);
		result['noInternal'] = false;
		result['limit'] = 0;
		return result;
	}
});
