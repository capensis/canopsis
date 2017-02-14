/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'ListAlarmWidget',
    after: ['TimeWindowUtils', 'DataUtils', 'WidgetFactory', 'UserconfigurationMixin', 'RinfopopMixin', 'SchemasLoader', 'CustomfilterlistMixin'],
    initialize: function(container, application) {
		    var timeWindowUtils = container.lookupFactory('utility:timewindow'),
            dataUtils = container.lookupFactory('utility:data'),
			      WidgetFactory = container.lookupFactory('factory:widget'),
			      UserConfigurationMixin = container.lookupFactory('mixin:userconfiguration');

            mx = container.lookupFactory('mixin:customfilterlist');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        // load the viewMixin
        var listOptions = {
            mixins: [
                UserConfigurationMixin,
                mx
            ]
        };

        /**
         * This widget allows to display alarms, with action possible on them.
         *
         * @memberOf canopsis.frontend.brick-listalarm
         * @mixes UserConfigurationMixin
         * @class WidgetListAlarm
         * @widget listalarm
         */
        var widget = WidgetFactory('listalarm',{
            // needs: ['login', 'application'],

            viewMixins: [
                ],

            searchText: '',
            isValidSearchText: true,
            humanReadableColumnNames: {
              'uid': '_id',
              'connector': 'v.connector',
              'connector_name': 'v.connector_name',
              'component': 'v.component',
              'resource': 'v.resource',
              'entity_id': 'd',
              'state': 'v.state.val',
              'status': 'v.status.val',
              'snooze': 'v.snooze',
              'ack': 'v.ack',
              'cancel': 'v.cancel',
              'ticket.val': 'v.ticket.val',
              'output': 'output',
              'opened': 't',
              'resolved': 'v.resolved',
              'domain': 'v.extra.domain',
              'perimeter': 'v.extra.perimeter',
              'last_state_change': 'v.state.t',
              'output': 'v.state.m'
            },


            hook: function() {
              this.set('model.user_filters', this.get('user_filters'));
              // this.saveUserConfiguration();
              console.error('waaat');
            }.observes('user_filters'),

            // rights: Ember.computed.alias('controllers'),


            /**
             * Create the widget and set widget params into Ember vars
             * @method init
             */
            init: function() {
                this._super.apply(this, arguments);
				        this.fetchAlarms();
                // this.valideExpression();
                set(this, 'rights', {list_filters: {checksum: true}});
				        set(this, 'store', DS.Store.extend({
                    container: get(this, 'container')
                }));
                this.showParams();
                this.setFields();
                this.loadRadioButtonView();                
                this.loadTemplates(this.get('model.popup'));
            },

            loadRadioButtonView: function () {
                view = Ember.View.extend({
                    tagName : "input",
                    type : "radio",
                    attributeBindings : [ "name", "type", "value", "checked:checked:" ],
                    click : function() {
                        console.error(this);
                        this.set("selection", this.$().val())
                    },
                    checked : function() {
                        return this.get("value") == this.get("selection");   
                    }.property()
                });
                try {
                  if (!Ember.RadioButton) {
                    Ember.RadioButton = view;
                  }
                } catch (err) {

                }
            },

            loadTemplates: function (templates) {
                Ember.columnTemplates = templates.map(function (obj) {
                  return {
                    columnName: obj.column,
                    columnTemplate: Ember.View.extend({
                      template: Ember.HTMLBars.compile(obj.template)
                    })                    
                  }
                })
            },

            showParams: function () {
                var controller = this;
                var params = ["popup", "title"];
                console.error("brick's parameters");
                params.forEach(function(param) {
                    console.error(param + ': ' + controller.get('model.' + param));
                });
                console.error("default_sort_column: " + controller.get('model.default_sort_column.property') + "-" + controller.get('model.default_sort_column.direction'));
                console.error("columnts: " + controller.get('model.columns').join(' ')),
                console.error("alarms_state_filter: " + controller.get('model.alarms_state_filter.state'))
                
            },

            /**
             * Set the reload to true in order to redraw events
             * extend the native refreshContent method
             * @method refreshContent
             */
            refreshContent: function () {
				          // Not implemented because backend too long, feature not useful for this widget
            },

            /**
             * Get the Alarms from the backend using the adapter
             * @method fetchAlarms
             */
            fetchAlarms: function(params) {
              var controller = this;
              var iParams = params || {};
              var filterState = this.get('model.alarms_state_filter.state') || 'opened';

              var query = {
                tstart: iParams['tstart'] || 0,
								tstop: iParams['tstop'] || 0,
								sort_key: iParams['sort_key'] || this.get('model.default_sort_column.property'),
           			sort_dir: iParams['sort_dir'] || this.get('model.default_sort_column.direction'),
                // filter: iParams['filter'] || this.get('model.filter'),
                search: iParams['search'] || '',
                opened: filterState == 'opened',
                resolved: filterState == 'resolved'
							};

              var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
            	adapter.findQuery('alerts', query).then(function (result) {
                    // onfullfillment
					          var alerts = get(result, 'data');
                    controller.setAlarmsForShow(alerts[0]['alarms']);

                    console.error('alerts::', alerts);
              }, function (reason) {
                    // onrejection
                    // console.error('ERROR in the adapter: ', reason);
              });
            },

            /**
             * Get the Alarms from the backend using the adapter
             * @method valideExpression
             */
            isValidExpression: function (expression) {
              var controller = this;

              var query = {
                expression: expression
              };

              var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alertexpression');
            	adapter.findQuery('alertexpression', query).then(function (result) {
                    // onfullfillment
					          var result = get(result, 'data');
                    console.error('alertexpression result', result);
                    controller.set('isValidSearchText', result[0]);
                    if (result[0]) {
                      var params = {};

                      params['search'] = expression;                      
                      
                      controller.fetchAlarms(params);
                    }
                    
              }, function (reason) {
                    // onrejection
                    console.error('ERROR in the adapter: ', reason);
              });
            },

            setFields: function () {
              this.set('fields', this.parseFields(get(this, 'model.columns')));              
            },

            setAlarmsForShow: function (alarms) {
              var fields = get(this, 'fields');
              var controller = this;
              var alarmsArr = alarms.map(function(alarm) {
                var newAlarm = {};
                fields.forEach(function(field) {
                    var val = get(Ember.Object.create(alarm), field.getValue);
                    newAlarm[field.name] = val;
                    newAlarm[field.humanName] = val;
                })
                return newAlarm;
              });

              this.set('alarms', alarmsArr);
            },

						parseFields: function (columns) {
              var controller = this;
							var fields = [];
							var sortColumn = this.get('model.default_sort_column.property');
							var order = this.get('model.default_sort_column.direction');
							

							fields = columns.map(function(column) {
								var obj = {};

								obj['name'] = controller.get('humanReadableColumnNames')[column] || 'v.' + column;
                obj['humanName'] = column;
								obj['isSortable'] = column == sortColumn;
								obj['isASC'] = order == 'ASC';
                obj['getValue'] = controller.get('humanReadableColumnNames')[column] || 'v.' + column;
								return obj;
							});

							return fields;
						},

            sortColumn: function() {
              var column = get(this, 'fields').findBy('name', get(this, 'controller.default_sort_column.property'));
              if (!column) {
                column = get(this, 'fields.firstObject');
                try {
                  column['isSortable'] = true;
                  column['isASC'] = get(this, 'controller.default_sort_column.property');                
                } catch (err) {
                  console.error('alarm!!!');
                }
                console.error('the column "' + get(this, 'controller.default_sort_column.property') + '" was not found.');
                return column;
              }
              return column;
            }.property('controller.default_sort_column.property', 'fields.[]'),



            actions: {
              updateSortField: function (field) {
                var params = {};

                params['sort_key'] = field.name;
                params['sort_dir'] = field.isASC ? 'ASC' : 'DESC';
                
                this.fetchAlarms(params);
              },
              
              search: function (text) {
                console.error('search', text);
                // console.error(this.isValidExpression(text));
                this.isValidExpression(text);
                  // console.error('request for search')
                // } else {
                  // this.set('isValidSearchText', false)
                // }
              }
            }

        }, listOptions);

        application.register('widget:listalarm', widget);
    }
});
