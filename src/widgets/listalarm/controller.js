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



        // var FormFactory = container.lookupFactory('factory:form');
        // var ModelFormController = container.lookupFactory('form:modelform');
        // var formOptions = {
        //     subclass: ModelFormController
        // };
        // var form = FormFactory('pbehaviorform', {
        //     init:function () {
        //         this._super();
        //         Ember.set(this, 'partials.buttons', ['formbutton-cancel', 'formbutton-ack', 'formbutton-ackandproblem'])
        //     }
        // }, formOptions);
        // application.register('form:pbehaviorform', form);
        // Ember.TEMPLATES['pbehaviorform'] = Ember.TEMPLATES['modelform'];

        // u=9;

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
              mx
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
              'state': 'v.state',
              'status': 'v.status',
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
              'output': 'v.state.m',
              'pbehaviors': 'v.pbehaviors',
              'extra_details': 'v.extra_details'
            },

            mandatoryFields: [
              {
                getValue: 'v.connector',
                name: 'connector',
                humanName: 'connector'
              },
              {
                getValue: 'v.connector_name',
                name: 'connector_name',
                humanName: 'connector_name'
              },
              {
                getValue: 'v.component',
                name: 'component',
                humanName: 'component'
              },
              {
                getValue: 'v.resource',
                name: 'resource',
                humanName: 'resource'
              },
              {
                getValue: 'v.state',
                name: 'state',
                humanName: 'state'
              }
            ],

            extraDeatialsEntities: [
              {
                name: 'snooze',
                value: 'v.snooze'
              },
              {
                name: 'ticket',
                value: 'v.ticket'
              },
              {
                name: 'ack',
                value: 'v.ack'
              },
              {
                name: 'pbehaviors',
                value: 'v.pbehaviors'
              }
            ],

            manualUpdateAlarms: 0,

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

                // this.set('filters_list', []);
                // this.get('filters_list').push(Ember.Object.create({filter: this.get('model.mixins').findBy('name', 'customfilterlist').default_filter, title: 'test', isActive: false}))

                // try {
                //   set(this, 'defaultFilter', this.get('model.mixins').findBy('name', 'customfilterlist').default_filter);
                // } catch (e) {
                //    set(this, 'defaultFilter', '');
                // }
                this.showParams();
                // this.setFields();
                this.loadTemplates(this.get('model.popup'));


                // var timestamps = this.getLivePeriod();
                try {
                  var fil = this.get('user_filters').findBy('isActive', true);
                  var filter = fil ? fil.filter : undefined;
                } catch (err) {
                  // console.error('error while selecting a filter', err);
                  var filter = undefined;
                }
                var filterState = this.get('model.alarms_state_filter.state') || 'resolved';


                var timestamps = this.defultTimestamps(filterState);

                this.set('alarmSearchOptions', {
                  // tstart: timestamps.tstart,
                  // tstop: timestamps.tstop,
                  // tstart: 0,
                  // tstop: 0,
                  opened: filterState == 'opened',
                  resolved: filterState == 'resolved',
                  // consolidations: [],
                  filter: filter,
                  // filter: {},                  
                  // search: '',
                  // sort_key: this.get('model.default_sort_column.property'),
                  // sort_dir: this.get('model.default_sort_column.direction'),
                  // skip: 0,
                  // limit: this.get('model.itemsPerPage') || 50
                });
            },

            defultTimestamps: function (state) {
              var tstart = 0, tstop = 0;
              if (state == 'opened') {
                tstart = 0;
                tstop = new Date().getTime();
              } else {
                var d = new Date();
                tstart = d.setMonth(d.getMonth() - 1)
                tstop = new Date().getTime();
              }
              return {
                tstart: tstart,
                tstop: tstop
              }
              // return {
              //   tstart: 0,
              //   tstop: 0
              // }
            },

            filtersObserver: function() {
              // console.error(this.get('selected_filter.filter'));
              this.set('alarmSearchOptions.filter', this.get('selected_filter.filter') || {});  
            }.observes('selected_filter'),

            
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
                    console.error('there are no suitable alarms');
                      return;
                  }
              }
              this.processEvent(event_type, 'handle', [crecords]);
              this.setPendingOperation(crecords);
              console.groupEnd();
            },

            alarmsTimestamp: 0,
            // for updating list of alarms
            timelineListener: function() {
              if (this.get('controllers.application.interval.timestamp.$gte')) {
                this.set('alarmSearchOptions.tstart', this.get('controllers.application.interval.timestamp.$gte') || 0);
                this.set('alarmSearchOptions.tstop', this.get('controllers.application.interval.timestamp.$lte') || 0);
              } else {
                var def = this.defultTimestamps(this.get('model.alarms_state_filter.state') || 'resolved');
                this.set('alarmSearchOptions.tstart', def.tstart);
                this.set('alarmSearchOptions.tstop', def.tstop);
              }
            }.observes('controllers.application.interval.timestamp'),

            alarmss: function() {
              var controller = this;
              this.set('loaded', false);              
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
                        'alarmSearchOptions.skip', 'alarmSearchOptions.limit', 'alarmSearchOptions.tstart',
                        'alarmSearchOptions.tstop', 'manualUpdateAlarms'),

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
                  // alarm['v']['pbehaviors'] = [
                  //   {
                  //     "tstop": 1483311600,
                  //     "enabled": false,
                  //     "name": "downtime",
                  //     "tstart": 1483225200,
                  //     "rrule": "FREQ=WEEKLY"
                  //   }
                  // ];

                  alarm['v']['extra_details'] = {};
                  controller.get('extraDeatialsEntities').forEach(function(item) {
                    alarm['v']['extra_details'][item.name] = Ember.Object.create(alarm).get(item.value);
                  })



                  

                  var newAlarm = Ember.Object.create();


                  controller.get('mandatoryFields').forEach(function(field) {
                      var val = get(Ember.Object.create(alarm), field.getValue);
                      
                      newAlarm[field.name] = val;
                      newAlarm[field.humanName] = val;
 
                  });

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

                  newAlarm['isExpanded'] = false;
                  
                  
                  newAlarm['id'] = alarm._id;
                  newAlarm['entity_id'] = alarm.d;

                  // newAlarm['state']['canceled'] = alarm.v.canceled;
                  newAlarm.set('state.canceled', alarm.v.canceled);

                  
                  // newAlarm['entity_id'] = '/resource/feeder/feeder/feeder_component/feeder_resource';
                  
                  // /resource/feeder/feeder/feeder_component/feeder_resource
                  // newAlarm['cancelled'] = alarm.v.cancelled;

                  newAlarm['changed'] = new Date().getTime();

                  // data for testing
                  // controller.set(newAlarm, 'linklist', {
                  //   'event_links': [
                  //     {
                  //       'url': 'http://tasks.info/?co=Demo',
                  //       'label': 'test'
                  //     }
                  //   ]
                  // });

                  // newAlarm['state'] = 1;
                  // newAlarm['v.state.val'] = 1;


                  newAlarm['linklist'] = {
                    'event_links': [
                      {
                        'url': 'http://tasks.info/?co=Demo',
                        'label': 'test'
                      }
                    ]
                  };

                  if (alarm.d.search('/resource/') == 0) {
                    newAlarm['source_type'] = 'resource';
                  };
                  if (alarm.d.search('/component/') == 0) {
                    newAlarm['source_type'] = 'component';
                  };
                  // newAlarm['pbehaviors'] = [
                  //   {
                  //     "tstop": 1483311600,
                  //     "enabled": false,
                  //     "name": "downtime",
                  //     "tstart": 1483225200,
                  //     "rrule": "FREQ=WEEKLY"
                  //   }
                  // ];

                  return newAlarm;
                });
              this.set('defTotal', Ember.totalAlarms);
              this.set('loaded', true);  
              return alarmsArr;

            }.property('alarmss.@each', 'fields.[]'),

            currPage: function () {
              this.set('alarmSearchOptions.limit', this.get('itemsPerPage'));
              this.set('alarmSearchOptions.skip', this.get('itemsPerPage') * (this.get('currentPage') - 1));
              
              // console.error('current page', this.get('currentPage'));
              // console.error('itemsPerPage', this.get('itemsPerPage'));              
            }.observes('currentPage', 'itemsPerPage'),

            paginationLastItemIndexx: function () {
                var itemsPerPage = get(this, 'itemsPerPage');
                var start = itemsPerPage * (this.currentPage - 1);
                return Math.min(start + itemsPerPage, get(this, 'itemsTotal'));
            }.property('widgetData', 'itemsTotal', 'itemsPerPage', 'currentPage'),

            // -------------------------------------------------------

            // filtersObserver: function() {
            //   try {
            //     var userFilters = this.get('user_filters');
            //     var filtersList = this.get('filters_list');
            //     if (userFilters || filtersList) {
            //       var filter = userFilters.findBy('isActive', true) || filtersList.findBy('isActive', true);
            //       if (filter) {
            //         var f = filter.filter || filter.get('filter');         
            //         // console.error(f.replace('state', 'v.state.val'));
            //         this.set('alarmSearchOptions.filter', f);
            //       } else {
            //         // console.error('there is no filter');
            //         this.set('alarmSearchOptions.filter', undefined);                
            //       }
            //     } else {
            //       this.set('alarmSearchOptions.filter', undefined);                                
            //     }
            //   } catch (err) {
            //       this.set('alarmSearchOptions.filter', undefined);                                
            //       // console.error('error while selecting a filter', err);
            //   }

            // }.observes('user_filters.@each.isActive', 'filters_list.@each.isActive'),

            loadTemplates: function (templates) {
                try {
                  Ember.columnTemplates = templates.map(function (obj) {
                    return {
                      columnName: obj.column,
                      columnTemplate: Ember.View.extend({
                        template: Ember.HTMLBars.compile(obj.template)
                      })                    
                    }
                  })
                } catch (err) {
                  // console.error('error while loading column templates');
                }
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
              this.set('manualUpdateAlarms', new Date().getTime());
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
                    // console.error('alerts::', alerts);
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
                    // console.error(result.data[0]);
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
              var column = get(this, 'fields').findBy('humanName', get(this, 'controller.default_sort_column.property'));
              if (!column) {
                column = get(this, 'fields.firstObject');
                try {
                  column['isSortable'] = true;
                  column['isASC'] = get(this, 'controller.default_sort_column.property');                
                } catch (err) {
                  // console.error('alarm!!!');
                }
                console.error('the column "' + get(this, 'controller.default_sort_column.property') + '" was not found.');
                return column;
              }
              return column;
            }.property('controller.default_sort_column.property', 'fields.[]'),


            updateAlarm: function (alarmId) {
              var controller = this;

              // this.get('alarms.firstObject').set('component', 'test');
              // t=4;
              var aa = this.get('alarms').findBy('id', alarmId);
              if (aa) {
                var self = this;
                var filterState = this.get('model.alarms_state_filter.state') || 'resolved';
                var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');

                  var f = {
                    'd': aa.get('entity_id')
                  }
                  adapter.findQuery('alarm', 'get-current-alarm', {'entity_id': aa.get('entity_id')}).then(function (a) {
                  // adapter.findQuery('alarm', 'get-current-alarm', { 'opened': filterState == 'opened', 'resolved': filterState == 'resolved', 'filter': "{'d': 'as'}" }).then(function (a) {
                    
                    if (a.success) {
                    // console.error('teest', a);
                      var fields = self.get('fields');
                      var alarm = a.data[0];
                  // var alarmsArr = get(this, 'alarmss').map(function(alarm) {
                      alarm.v =alarm.value;
                      // aa._id = 'aa';
                      aa.entity_id= alarm.data_id;
                      aa.d= alarm.data_id;

                      aa.set('extra_details', Ember.Object.create());
                      controller.get('extraDeatialsEntities').forEach(function(item) {
                        // aa['extra_details'][item.name] = Ember.Object.create(alarm).get(item.value);
                        aa.set('extra_details.' + item.name, Ember.Object.create(alarm).get(item.value));
                      })
                    
                      // alarm['v']['pbehaviors'] = [
                      //   {
                      //     "tstop": 1483311600,
                      //     "enabled": false,
                      //     "name": "downtime",
                      //     "tstart": 1483225200,
                      //     "rrule": "FREQ=WEEKLY"
                      //   }
                      // ];
                      var newAlarm = Ember.Object.create();

                      fields.forEach(function(field) {
                          if (field.humanName != 'extra_details') {
                            var val = get(Ember.Object.create(alarm), field.getValue);
                            // controller.set(newAlarm, field.name, val);
                            // controller.set(newAlarm, field.humanName, val);
                            
                            // newAlarm[field.name] = val;
                            // newAlarm[field.humanName] = val;

                            // aa.set(field.name, val);
                            aa.set(field.humanName, val);
                          }

                      });
                      // controller.set(newAlarm, 'isSelected', false);
                      
                      aa.set('isSelected', false);
                      // controller.set(newAlarm, 'id', alarm.get('_id'));

                      aa.set('isExpanded', false);

                      aa.set('canceled', get(Ember.Object.create(alarm), 'v.canceled'));

                      // aa.set('status.val', 4);
                      
                      
                      // newAlarm['id'] = alarm._id;
                      // newAlarm['entity_id'] = alarm.d;
                      // newAlarm['cancelled'] = alarm.v.cancelled;

                      // data for testing
                      // controller.set(newAlarm, 'linklist', {
                      //   'event_links': [
                      //     {
                      //       'url': 'http://tasks.info/?co=Demo',
                      //       'label': 'test'
                      //     }
                      //   ]
                      // });

                      // newAlarm['state'] = 2;
                      // newAlarm['v.state.val'] = 2;

                      // newAlarm['changed'] = new Date().getTime;


                      // newAlarm['linklist'] = {
                      //   'event_links': [
                      //     {
                      //       'url': 'http://tasks.info/?co=Demo',
                      //       'label': 'test'
                      //     }
                      //   ]
                      // };

                      // if (alarm.d.search('/resource/') == 0) {
                      //   aa.set['source_type'] = 'resource';
                      // };
                      // if (alarm.d.search('/component/') == 0) {
                      //   aa['source_type'] = 'component';
                      // };

                      // var t = self.get('alarms').objectAt(0);
                      // Ember.set(t, 'state', 1);

                      // var aa = this.get('alarms').findBy('entity_id', alarmId);
                      
                      
                      Ember.set(aa, 'changed', new Date().getTime());
                      // self.set('alarms.firstObject', 1);
                      t=2;


                    } else {
                      console.error('unsuccessful request');
                    }
              })
                  
              } else {
                console.error('alarm not found');
              }

    
                // newAlarm['pbehaviors'] = [
                //   {
                //     "tstop": 1483311600,
                //     "enabled": false,
                //     "name": "downtime",
                //     "tstart": 1483225200,
                //     "rrule": "FREQ=WEEKLY"
                //   }
                // ];

                // return newAlarm;
              // });

              // var a = this.get('alarms.firstObject');


              // var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
            
              //   adapter.findQuery('alarm', 'get-current-alarm', {'entity_id': a.get('entity_id')}).then(function (alarms) {
              //     console.error('teest', alarms);
              // })

            },



            actions: {
              massAction: function (action) {
                this.sendEventt(action.mixin_name);
              },

              sendCustomAction: function (action, alarm) {
                // console.error('controller', action, alarm);
                this.sendEventt(action.mixin_name, alarm);
              },

              updateSortField: function (field) {
                this.set('alarmSearchOptions.sort_key', field.name);
                this.set('alarmSearchOptions.sort_dir', field.isASC ? 'ASC' : 'DESC');
              },
              
              search: function (text) {
                var controller = this;
                // console.error('search ', text);
                this.isValidExpression(text).then(function(result) {
                  controller.set('isValidSearchText', result);
                  if (result) {
                    controller.set('alarmSearchOptions.search', text);
                    controller.set('manualUpdateAlarms', new Date().getTime());
                  }
                })
              },

            }

        }, listOptions);

        application.register('widget:listalarm', widget);
    }
});