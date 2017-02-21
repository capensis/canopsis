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
    after: ['NotificationUtils' ,'TimeWindowUtils', 'DataUtils', 'WidgetFactory', 'UserconfigurationMixin', 'RinfopopMixin', 'SchemasLoader', 'CustomfilterlistMixin', 'CustomSendeventMixin'],
    initialize: function(container, application) {
		    var timeWindowUtils = container.lookupFactory('utility:timewindow'),
            dataUtils = container.lookupFactory('utility:data'),
			      WidgetFactory = container.lookupFactory('factory:widget'),
			      UserConfigurationMixin = container.lookupFactory('mixin:userconfiguration');
			      SendeventMixin = container.lookupFactory('mixin:customsendevent');            
            notificationUtils = container.lookupFactory('utility:notification');

            mx = container.lookupFactory('mixin:customfilterlist');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        // load the viewMixin
        var listOptions = {
            mixins: [
                UserConfigurationMixin,
                mx,
                SendeventMixin
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
            needs: ['login', 'application'],

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

            /**
             * Create the widget and set widget params into Ember vars
             * @method init
             */
            init: function() {
                this._super.apply(this, arguments);
				        // this.fetchAlarms();
                // this.valideExpression();
                set(this, 'loaded', false);
                set(this, 'rights', {list_filters: {checksum: true}});
				        set(this, 'store', DS.Store.extend({
                    container: get(this, 'container')
                }));
                this.showParams();
                // this.setFields();
                this.loadRadioButtonView();                
                this.loadTemplates(this.get('model.popup'));


                // var timestamps = this.getLivePeriod();
                var fil = this.get('user_filters').findBy('isActive', true);
                var filter = fil ? fil.filter : undefined;
                var filterState = this.get('model.alarms_state_filter.state') || 'opened';

                this.set('alarmSearchOptions', {
                  tstart: 0,
                  tstop: 0,
                  opened: filterState == 'opened',
                  resolved: filterState == 'resolved',
                  // consolidations: [],
                  filter: filter,
                  search: '',
                  sort_key: this.get('model.default_sort_column.property'),
                  sort_dir: this.get('model.default_sort_column.direction'),
                  skip: 0,
                  limit: this.get('model.itemsPerPage') || 50
                });
            },

            // rewrite totalPages
            totalPagess: function() {
                if (get(this, 'itemsTotal') === 0) {
                  this.set('totalPages', 0);
                    // return 0;
                } else {
                    var itemsPerPage = get(this, 'itemsPerPage');
                    this.set('totalPages', Math.ceil(get(this, 'itemsTotal') / itemsPerPage));
                }
            }.observes('itemsPerPagePropositionSelected', 'itemsTotal', 'itemsPerPage'),


            sendEventt: function(event_type, crecord) {
              console.group('sendEvent:', arguments);
              this.stopRefresh();
              var crecords = [];
              if (!isNone(crecord)) {
                  console.log('event:', event_type, crecord);
                  crecords.pushObject(crecord);
              }
              else {
                  if (this.get('loaded')) {
                    var content = get(this, 'alarms');
                  } else {
                    var content = Ember.A();
                  }
                  var selected = content.filterBy('isSelected', true);
                  crecords = this.filterUsableCrecords(event_type, selected);
                  console.log('events:', event_type, crecords);
                  if(!crecords.length) {
                      notificationUtils.warning(
                          'No matching event found for event:',
                          event_type
                      );
                      return;
                  }
              }
              this.processEvent(event_type, 'handle', [crecords]);
              this.setPendingOperation(crecords);
              console.groupEnd();
            },

            // getLivePeriod: function () {
            //     var tw = timeWindowUtils.getFromTo(
            //         get(this, 'time_window'),
            //         get(this, 'time_window_offset')
            //     );
            //     var from = tw[0],
            //         to = tw[1];

            //     /* live reporting support */
            //     var liveFrom = get(this, 'from'),
            //         liveTo = get(this, 'to');
            //     if (!isNone(liveFrom)) {
            //         from = liveFrom;
            //     }
            //     if (!isNone(liveTo)) {
            //         to = liveTo;
            //     }
            //     return {'from': from, 'to': to};
            // },

            findItems: function() {
              console.error('FINDITMENS');
            }.observes('widgetData'),

            alarmss: function() {
              var controller = this;
              this.set('loaded', false);              

              // var timestamps = this.getLivePeriod();
              // this.set('alarmSearchOptions.tstart', 0);
              // this.set('alarmSearchOptions.tstop', 0);              

              var options = this.get('alarmSearchOptions');
              console.error('reload original alarms with params', options);              
              var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
              
              return DS.PromiseArray.create({
                promise: adapter.findQuery('alerts', options).then(function (alarms) {
                  if (alarms.success) {
                    console.error('loaded alarms: ', get(alarms, 'data.firstObject.alarms'));
                    Ember.totalAlarms = get(alarms, 'data.firstObject.total');
                    return get(alarms, 'data.firstObject.alarms');
                  } else {
                    throw new Error(get(alarms, 'data.msg'));
                  }
                }, function (reason) {
                  console.error('ERROR in the adapter: ', reason);
                  return [];
                })
                .catch(function (err) {
                  console.error('unexpected error ', err);
                  return [];                
                })
              })

            }.property('alarmSearchOptions.search', 'alarmSearchOptions.resolved',
                        'alarmSearchOptions.sort_key', 'alarmSearchOptions.sort_dir', 'alarmSearchOptions.filter',
                        'alarmSearchOptions.skip', 'alarmSearchOptions.limit'),

            fields: function() {
              return this.parseFields(get(this, 'model.columns'));
            }.property('model.columns'),

            widgetDataMetas: function () {
              return {total: this.get('defTotal') || 0};
            }.property('defTotal'),

            alarms: function() {
              var controller = this;
              var fields = get(this, 'fields');
              var alarmsArr = get(this, 'alarmss').map(function(alarm) {
                  var newAlarm = Ember.Object.create();
                  fields.forEach(function(field) {
                      var val = get(Ember.Object.create(alarm), field.getValue);
                      // controller.set(newAlarm, field.name, val);
                      // controller.set(newAlarm, field.humanName, val);
                      
                      newAlarm[field.name] = val;
                      newAlarm[field.humanName] = val;
 
                  });
                  // controller.set(newAlarm, 'isSelected', false);
                  
                  newAlarm['isSelected'] = false;
                  // controller.set(newAlarm, 'id', alarm.get('_id'));
                  
                  newAlarm['id'] = alarm._id;

                  // data for testing
                  // controller.set(newAlarm, 'linklist', {
                  //   'event_links': [
                  //     {
                  //       'url': 'http://tasks.info/?co=Demo',
                  //       'label': 'test'
                  //     }
                  //   ]
                  // });


                  newAlarm['linklist'] = {
                    'event_links': [
                      {
                        'url': 'http://tasks.info/?co=Demo',
                        'label': 'test'
                      }
                    ]
                  }
                  return newAlarm;
                });
              this.set('defTotal', Ember.totalAlarms);
              this.set('loaded', true);  
              return alarmsArr;

            }.property('alarmss.@each', 'fields.[]'),

            currPage: function () {
              this.set('alarmSearchOptions.limit', this.get('itemsPerPage'));
              this.set('alarmSearchOptions.skip', this.get('itemsPerPage') * (this.get('currentPage') - 1));
              
              console.error('current page', this.get('currentPage'));
              console.error('itemsPerPage', this.get('itemsPerPage'));              
            }.observes('currentPage', 'itemsPerPage'),

            paginationLastItemIndexx: function () {
                var itemsPerPage = get(this, 'itemsPerPage');
                var start = itemsPerPage * (this.currentPage - 1);
                return Math.min(start + itemsPerPage, get(this, 'itemsTotal'));
            }.property('widgetData', 'itemsTotal', 'itemsPerPage', 'currentPage'),

            // -------------------------------------------------------

            filtersObserver: function() {
              var filter = this.get('user_filters').findBy('isActive', true);
              if (filter) {
                console.error(filter.filter || filter.get('filter'));
                this.set('alarmSearchOptions.filter', filter.filter || filter.get('filter'));
              } else {
                console.error('there is no filter');
                this.set('alarmSearchOptions.filter', undefined);                
              }
            }.observes('user_filters.@each.isActive'),

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

              return DS.PromiseObject.create({
                promise: adapter.findQuery('alertexpression', query).then(function (result) {
                  if (result.success) {
                    console.error(result.data[0]);
                    return result.data[0];
                  } else {
                    throw new Error(result.data.msg);
                  }
                }, function (reason) {
                  console.error('ERROR in the adapter: ', reason);
                  return false;
                })
                .catch(function (err) {
                  console.error('unexpected error ', err);
                  return false;                
                })
              })

              // var controller = this;

              // var query = {
              //   expression: expression
              // };

              // var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alertexpression');
            	// adapter.findQuery('alertexpression', query).then(function (result) {
              //       // onfullfillment
					    //       var result = get(result, 'data');
              //       console.error('alertexpression result', result);
              //       controller.set('isValidSearchText', result[0]);
              //       if (result[0]) {
              //         var params = {};

              //         params['search'] = expression;                      
                      
              //         controller.fetchAlarms(params);
              //       }
                    
              // }, function (reason) {
              //       // onrejection
              //       console.error('ERROR in the adapter: ', reason);
              // });
            },

            // setFields: function () {
            //   this.set('fields', this.parseFields(get(this, 'model.columns')));              
            // },

            // setAlarmsForShow: function (alarms) {
            //   var fields = get(this, 'fields');
            //   var controller = this;
            //   var alarmsArr = alarms.map(function(alarm) {
            //     var newAlarm = {};
            //     fields.forEach(function(field) {
            //         var val = get(Ember.Object.create(alarm), field.getValue);
            //         newAlarm[field.name] = val;
            //         newAlarm[field.humanName] = val;
            //     })
            //     return newAlarm;
            //   });

            //   this.set('alarms', alarmsArr);
            // },

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
              var column = get(this, 'fields').findBy('humanName', get(this, 'controller.default_sort_column.property'));
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
              massAction: function (action) {
                this.sendEventt(action.mixin_name);
              },

              sendCustomAction: function (action, alarm) {
                console.error('controller', action, alarm);
                this.sendEventt(action.mixin_name, alarm);
              },

              updateSortField: function (field) {
                this.set('alarmSearchOptions.sort_key', field.name);
                this.set('alarmSearchOptions.sort_dir', field.isASC ? 'ASC' : 'DESC');
              },
              
              search: function (text) {
                var controller = this;
                console.error('search ', text);
                this.isValidExpression(text).then(function(result) {
                  controller.set('isValidSearchText', result);
                  if (result) {
                    controller.set('alarmSearchOptions.search', text);
                  }
                })
              }
            }

        }, listOptions);

        application.register('widget:listalarm', widget);
    }
});
