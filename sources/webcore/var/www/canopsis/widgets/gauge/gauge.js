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
Ext.define('widgets.gauge.gauge' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.gauge',

	logAuthor: '[gauge]',

	colorStart: '#6FADCF',
	colorStop: '#8FC0DA',
	colorWarn: '#FFA500',
	gaugeColor: '#E1E6FA',
	titleFontColor: '#3E576F',
	gaugeWidthScale: 1,
	showMinMax: true,
	shadowOpacity: 0.7,

	labelSize: 25,
	maxValue: 100,
	minValue: 0,

	// Internals
	gauge: undefined,

	label: '',
	gaugeTitle: '',
	gaugeLabel: '',
	lastValue: 0,
	displayUnit: true,
	time_window: 0,

	initComponent: function() {
		this.gaugeTitle = this.title;
		this.title = '';

		this.nodesByID = parseNodes(this.nodes);

		this.haveCounter = false;
		//search counter
		for(var i = 0; i < this.nodes.length; i++) {
			var node = this.nodes[i];
			if(node['type'] && node['type'] === 'COUNTER') {
				this.haveCounter = true;
			}
		}

		this.callParent(arguments);
	},

	createGauge: function(value) {
		var me = this;

		var renderer = function(val) {
			if(me.humanReadable) {
				return rdr_humanreadable_value(val, me.symbol);
			}
			else {
				if(me.symbol) {
					return val + ' ' + me.symbol;
				}
				else {
					return val;
				}
			}
		};

		if(!value) {
			value = 0;
		}

		if(this.autoTitle && this.nodesByID) {
			var node = this.nodesByID[Ext.Object.getKeys(this.nodesByID)[0]];

			var component = node.component;
			var source_type = node.source_type;

			if (source_type === 'resource') {
				var resource = node.resource;
				this.gaugeTitle = resource + ' ' + _('on') + ' ' + component;
			}
			else {
				this.gaugeTitle = component;
			}
		}

		var opts = {
			id: this.wcontainerId,
			value: value,
			gaugeWidthScale: this.gaugeWidthScale,
			titleFontColor: this.titleFontColor,
			showMinMax: this.showMinMax,
			levelColorsGradient: true,
			min: this.minValue,
			max: this.maxValue,
			shadowOpacity: this.shadowOpacity,
			title: this.gaugeTitle,
			label: this.gaugeLabel,
			gaugeColor: this.gaugeColor,
			textRenderer: renderer,
			symbol: this.bunit
		};

		if(this.exportMode) {
			opts['showInnerShadow'] = 0;
			opts['shadowVerticalOffset'] = 0;
			opts['shadowOpacity'] = 0;
			opts['shadowSize'] = 0;
			opts['startAnimationType'] = 1;
			opts['refreshAnimationTime'] = 1;
		}

		if(this.levelThresholds) {
			opts.levelColorsGradient = false;
			opts.levelColors = [this.colorStart];
			opts.levelThresholds = [];

			if(this.warnValue) {
				opts.levelColors.push(this.colorWarn);
				opts.levelThresholds.push(this.warnValue);
			}

			if(this.critValue) {
				opts.levelColors.push(this.colorStop);
				opts.levelThresholds.push(this.critValue);
			}
		}
		else {
			opts.levelColors = [this.colorStart];
			opts.levelThresholds = [];
		}

		log.debug('Gauge options:', this.logAuthor);
		log.dump(opts);

		opts['cwidget'] = this;
		this.gauge = new JustGage(opts);
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);

		delete this.gauge;
		this.createGauge();
		this.gauge.refresh(this.lastValue);
	},

	getNodeInfo: function(from,to) {
		this.processNodes();

		if(!this.haveCounter || !this.time_window) {
			from = to;
		}

		if(this.nodesByID) {
			Ext.Ajax.request({
				url: '/perfstore/values' + '/' + parseInt(from/1000) + '/' + parseInt(to/1000),
				scope: this,
				params: this.post_params,
				method: 'POST',
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);

					if(Ext.Object.getSize(this.nodesByID) > 1) {
						data = data.data;
					}
					else {
						data = data.data[0];
					}

					this._onRefresh(data);
				},

				failure: function(result, request) {
					void(result);

					log.error('Impossible to get Node informations, Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		}
	},

	processNodes: function() {
		var post_params = [];

		Ext.Object.each(this.nodesByID, function(id, node) {
			post_params.push({
				id: id,
				metrics: node.metrics
			});
		}, this);

		this.post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points
		};
	},

	onRefresh: function(data) {
		log.debug('onRefresh', this.logAuthor);

		if(data) {
			if(this.getEl().isMasked && !this.isDisabled()) {
				this.getEl().unmask();
			}

			if(data.min) {
				this.minValue = data.min;
			}

			if(data.max) {
				this.maxValue = data.max;
			}

			if(data.thld_warn) {
				this.warnValue = data.thld_warn;
			}

			if(data.thld_crit) {
				this.critValue = data.thld_crit;
			}

			if(data.bunit && this.displayUnit) {
				this.bunit = data.bunit;
			}

			var fields = undefined;

			//get first node
			fields = this.nodesByID[Ext.Object.getKeys(this.nodesByID)[0]];


			if(fields) {
				//update metric name
				if(fields.label) {
					this.gaugeLabel = fields.label;
				}
				else {
					this.gaugeLabel = data.metric;
				}

				//update metric value
				if(fields.max) {
					this.maxValue = fields.max;
				}

				if(fields.min) {
					this.minValue = fields.min;
				}

				if(fields.thld_warn) {
					this.warnValue = fields.thld_warn;
				}

				if(fields.thld_crit) {
					this.critValue = fields.thld_crit;
				}
			}

			try {
				if(data.values) {
					this.lastValue = data.values[data.values.length - 1][1];

					if(!this.gauge) {
						this.createGauge(this.lastValue);
					}
					else {
						this.gauge.refresh(this.lastValue);
					}
				}
			}
			catch (err) {
				log.error('Error while set value:' + err, this.logAuthor);
			}
		}
		else {
			this.getEl().mask(_('No data received from webserver'));
			log.debug('No data', this.logAuthor);
		}
	}
});
