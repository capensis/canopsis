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
    after: ['NotificationUtils', 'TimeWindowUtils', 'DataUtils', 'WidgetFactory', 'UserconfigurationMixin', 'RinfopopMixin', 'SchemasLoader', 'CustomfilterlistMixin', 'CustomSendeventMixin'],
    initialize: function (container, application) {
        var timeWindowUtils = container.lookupFactory('utility:timewindow'),
            dataUtils = container.lookupFactory('utility:data'),
            formsUtils = container.lookupFactory('utility:forms');

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
        var widget = WidgetFactory('listalarm', {
            needs: ['login', 'application'],

            viewMixins: [
                mx
            ],

            /**
             * @property searchText
             */
            searchText: '',

            /**
             * @property isValidSearchText
             */
            isValidSearchText: true,

            /**
             * @property humanReadableColumnNames
             * @description Maps an alarm's properties and human readeble label
             */
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
                'pbehaviors': 'pbehaviors',
                'extra_details': 'v.extra_details',
                'initial_output': 'v.initial_output'
            },

            alarmDBColumn: ["status",
                "resolved",
                "resource",
                "tags",
                "ack",
                "extra",
                "component",
                "creation_date",
                "connector",
                "canceled",
                "state",
                "connector_name",
                "steps",
                "initial_output",
                "last_update_date",
                "snooze",
                "ticket",
                "hard_limit",
                "display_name"],

            entityDBColumn: ["impact",
                "name",
                "measurements",
                "depends",
                "infos",
                "type",
                "disable_history",
                "enabled",
                "enable_history"],


            /**
             * @property mandatoryFields
             * @description Properties that always have to present in order to correctly send events
             */
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
                },
                {
                    getValue: 'v.extra_details',
                    name: 'extra_details',
                    humanName: 'extra details'
                },
                {
                    getValue: 'd',
                    name: 'd',
                    humanName: 'd'
                }
            ],

            /**
             * @property extraDeatialsEntities
             * @description List of properties that are displayed in extra_details column
             */
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
                    name: 'done',
                    value: 'v.done'
                },
                {
                    name: 'ack',
                    value: 'v.ack'
                },
                {
                    name: 'pbehaviors',
                    value: 'pbehaviors'
                }
            ],

            /**
             * @property manualUpdateAlarms
             */
            manualUpdateAlarms: 0,


            /**
             * @property alarmsTimestamp
             */
            alarmsTimestamp: 0,

            /**
             * Create the widget and set widget params into Ember vars
             * @method init
             */
            init: function () {
                this._super.apply(this, arguments);

                set(this, 'loaded', false);
                set(this, 'rights', { list_filters: { checksum: true } });
                set(this, 'store', DS.Store.extend({
                    container: get(this, 'container')
                }));

                this.showParams();
                this.loadTemplates(this.get('model.popup'));

                try {
                    var fil = this.get('user_filters').findBy('isActive', true);
                    var filter = fil ? fil.filter : undefined;
                } catch (err) {
                    var filter = undefined;
                }

                var filterState = this.get('model.alarms_state_filter.state') || 'resolved';
                var hide_resources = this.get('model.hide_resources')
                var timestamps = this.defultTimestamps(filterState);

                this.set('alarmSearchOptions', {
                    opened: filterState == 'opened',
                    resolved: filterState == 'resolved',
                    hide_resources: hide_resources,
                    lookups: JSON.stringify(["pbehaviors"]) ,
                    filter: filter,
                    sort_key: this.get('model.default_sort_column.property'),
                    sort_dir: this.get('model.default_sort_column.direction'),
                });
            },

            /**
             * @method defultTimestamps
             * @argument state
             */
            defultTimestamps: function (state) {
                var tstart = 0, tstop = 0;
                if (state == 'opened') {
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
            },

            /**
             * @method filtersObserver
             */
            filtersObserver: function () {
                this.set('alarmSearchOptions.filter', this.get('selected_filter.filter') || {});
            }.observes('selected_filter'),

            /**
             * @property totalPagess
             */
            totalPagess: function () {
                if (get(this, 'itemsTotal') === 0) {
                    this.set('totalPages', 0);
                    // return 0;
                } else {
                    var itemsPerPage = get(this, 'itemsPerPage');
                    this.set('totalPages', Math.ceil(get(this, 'itemsTotal') / itemsPerPage));
                }
            }.observes('itemsPerPagePropositionSelected', 'itemsTotal', 'itemsPerPage'),



            /**
             * @method sendEventCustom
             * @desciption Is rewritten from SendEventMixin
             */
            sendEventCustom: function (event_type, crecord) {
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

                    // custom event mixin does not support 'bulk_pbehavior' assignment
                    // controller is responsible for that
                    if (event_type === 'bulk_pbehavior') {
                        this.createMultiplePBehaviors(selected);
                        return;
                    }

                    crecords = this.filterUsableCrecords(event_type, selected);
                    console.log('events:', event_type, crecords);
                    if (!crecords.length) {
                        console.error('there are no suitable alarms');
                        return;
                    }
                }
                this.processEvent(event_type, 'handle', [crecords]);
                this.setPendingOperation(crecords);
                console.groupEnd();
            },

            /**
             * @method timelineListener
             */
            timelineListener: function () {
                if (this.get('controllers.application.interval.timestamp.$gte')) {
                    this.set('alarmSearchOptions.tstart', this.get('controllers.application.interval.timestamp.$gte') || 0);
                    this.set('alarmSearchOptions.tstop', this.get('controllers.application.interval.timestamp.$lte') || 0);
                } else {
                    var def = this.defultTimestamps(this.get('model.alarms_state_filter.state') || 'resolved');
                    this.set('alarmSearchOptions.tstart', def.tstart);
                    this.set('alarmSearchOptions.tstop', def.tstop);
                }
            }.observes('controllers.application.interval.timestamp'),


			replaceColumnsName: function (search, fields, columnMap) {
				var conditions = search.split(/ OR | AND /)
				var conditionOps = ["<=", "<", "=", "!=", ">=", ">", "CONTAINS", "LIKE"]

				for (itCond = 0; itCond < conditions.length; itCond++) {
					for(itCondOp = 0; itCondOp < conditionOps.length; itCondOp++) {
						if (conditions[itCond].indexOf(conditionOps[itCondOp]) >= 0) {

							humanName = conditions[itCond].split(conditionOps[itCondOp])[0].trim()

							var hasNot = false
							if (humanName.indexOf("NOT") >= 0) {
								hasNot = true
								humanName = humanName.replace("NOT", "")
								humanName = humanName.trim()
							}

							var found = false
							var technicalName = ""
							var itField = 0

							while (!found && itField < fields.length){
								if(fields[itField]["humanName"] == humanName){
									var technicalName = fields[itField].name
									found=true
								}
								itField ++
							}

							if(!found){
								technicalName = columnMap[humanName.toLowerCase()] || ""
							}

							if(technicalName === "" && humanName.startsWith("infos.")){
								technicalName = "entity." + humanName
							}

							if(technicalName.startsWith("infos.")){
								technicalName = "entity." + technicalName
							}

							if (technicalName !== "") {
								updatedCondition = conditions[itCond].replace(humanName, technicalName)
								search = search.replace(conditions[itCond], updatedCondition)
							}
						}
					}
				}
				return search
			},

            /**
             * @property originalAlarms
             */
            originalAlarms: function () {
                var controller = this;
                this.set('loaded', false);
                var options = this.get('alarmSearchOptions');

                // fix for applying a filter after refreshing the widget
                this.set('alarmSearchOptions.filter', this.get('selected_filter.filter') || {});

                //don't touch this or the backend will explode
                if (!options.filter){
                    options.filter = "{}";
				}
				if(options.search !== undefined){
					options.search = options.search.trim()
				}

				options['natural_search'] = true;
				controller.set('isNaturalSearch', true);
				if(options.search !== undefined && options.search.startsWith("- ")){
					options.search = options.search.substring(2)
					options.search = this.get("replaceColumnsName")(options.search, this.get("fields"), this.get('humanReadableColumnNames'))
					options['natural_search'] = false;
					controller.set('isNaturalSearch', false);
				}

                var columns = []
				get(this, 'model.widget_columns').forEach(function(element) {
					columns.push(element.value)
				})

                var prefixed_columns = [];
                for (idx = 0; idx < columns.length; idx++) {
                    depth_one = columns[idx].split(".", 1)[0];
                    if (this.get("entityDBColumn").indexOf(depth_one) >= 0) {
                        prefixed_columns.push("entity." + columns[idx]);
                    }
                    if (this.get("alarmDBColumn").indexOf(depth_one) >= 0) {
                        prefixed_columns.push("v." + columns[idx]);
                    }
                }
                options["active_columns"] = prefixed_columns;

                var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
                return DS.PromiseArray.create({
                    promise: adapter.findQuery('alerts', options).then(function (alarms) {
                        if (alarms.success) {
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


            /**
             * @property fields
             * @description Stores choosen by user fileds
             */
            fields: function() {
                return this.parseFields(get(this, 'model.widget_columns'));
            }.property('model.widget_columns'),

            /**
             * @property widgetDataMetas
             */
            widgetDataMetas: function () {
                return { total: this.get('defTotal') || 0 };
            }.property('defTotal'),

            /**
             * @property alarms
             * @description Stores all alams that have to be displayed
             */
            alarms: function () {
                var controller = this;
                var fields = get(this, 'fields');
                var alarmsArr = get(this, 'originalAlarms').map(function (alarm) {
                    alarm['v']['extra_details'] = {};
                    controller.get('extraDeatialsEntities').forEach(function (item) {
                        alarm['v']['extra_details'][item.name] = Ember.Object.create(alarm).get(item.value);
                    })
                    var newAlarm = Ember.Object.create();

                    controller.get('mandatoryFields').forEach(function (field) {
                        var val = get(Ember.Object.create(alarm), field.getValue);
                        newAlarm[field.name] = val;
                        newAlarm[field.humanName] = val;

                    });

                    fields.forEach(function (field) {
                        var val = get(Ember.Object.create(alarm), field.getValue);
                        newAlarm[field.name] = val;
                        newAlarm[field.humanName] = val;
                    });

                    newAlarm['isSelected'] = false;
                    newAlarm['isExpanded'] = false;
                    newAlarm['id'] = alarm._id;
                    newAlarm['entity_id'] = alarm.d;
                    if (newAlarm.get('state') !== undefined){
                        newAlarm.set('state.canceled', alarm.v.canceled);
                    }
                    newAlarm['changed'] = new Date().getTime();
                    newAlarm.links = alarm.links;

                    if (alarm.v.resource) {
                        newAlarm['source_type'] = 'resource';
                    } else {
                        newAlarm['source_type'] = 'component';
                    }
                    return newAlarm;
                });
                this.set('defTotal', Ember.totalAlarms);
                this.set('loaded', true);
                return alarmsArr;

            }.property('originalAlarms.@each', 'fields.[]'),

            /**
             * @method currPage
             */
            currPage: function () {
                this.set('alarmSearchOptions.limit', this.get('itemsPerPage'));
                this.set('alarmSearchOptions.skip', this.get('itemsPerPage') * (this.get('currentPage') - 1));
            }.observes('currentPage', 'itemsPerPage'),

            /**
             * @method paginationLastItemIndexx
             */
            paginationLastItemIndexx: function () {
                var itemsPerPage = get(this, 'itemsPerPage');
                var start = itemsPerPage * (this.currentPage - 1);
                return Math.min(start + itemsPerPage, get(this, 'itemsTotal'));
            }.property('widgetData', 'itemsTotal', 'itemsPerPage', 'currentPage'),

            /**
             * @method loadTemplates
             * @description Loads templates for popups columns
             */
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
                }
            },

            /**
             * @method showParams
             * @desctiption Show all user_configuration form's params in a console
             */
            showParams: function () {
                var controller = this;
                var params = ["popup", "title"];
                params.forEach(function (param) {
                    console.log(param + ': ' + controller.get('model.' + param));
                });

            },

            /**
             * Set the reload to true in order to redraw events
             * extend the native refreshContent method
             * @method refreshContent
             */
            refreshContent: function () {
                this.set('manualUpdateAlarms', new Date().getTime());

                // this.set('manualUpdateAlarms', new Date().getTime());
                // Not implemented because backend too long, feature not useful for this widget
            },

            /**
             * Get the Alarms from the backend using the adapter
             * @method fetchAlarms
             */
            fetchAlarms: function (params) {
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
                    var alerts = get(result, 'data');
                    controller.setAlarmsForShow(alerts[0]['alarms']);
                }, function (reason) {
                    console.error('ERROR in the adapter: ', reason);
                });
            },



            /**
             * @method parseFields
             */
            parseFields: function (columns) {
                var controller = this;
                var fields = [];
                var sortColumn = this.get('model.default_sort_column.property');
                var order = this.get('model.default_sort_column.direction');

				if (columns === undefined){
					columns = []
				}

                fields = columns.map(function (column) {
                    var obj = {};
                    if (column.value.startsWith('infos')){
                      obj['name'] = column.value;
                      obj['humanName'] = column.label || column.value;
                      obj['isSortable'] = column.value == sortColumn;
                      obj['isASC'] = order == 'ASC';
                      obj['getValue'] = column.value;
                    } else {
                      obj['name'] = controller.get('humanReadableColumnNames')[column.value] || 'v.' + column.value;
                      obj['humanName'] = column.label || column.value;
                      obj['isSortable'] = column.value == sortColumn;
                      obj['isASC'] = order == 'ASC';
                      obj['getValue'] = controller.get('humanReadableColumnNames')[column.value] || 'v.' + column.value;
                    }
                    return obj;
                });

                return fields;
            },

            /**
             * @property sortColumn
             * @description Current sortable column
             */
            sortColumn: function () {
                var column = get(this, 'fields').findBy('humanName', get(this, 'controller.default_sort_column.property'));
                if (!column) {
                    column = get(this, 'fields.firstObject');
                    try {
                        column['isSortable'] = true;
                        column['isASC'] = get(this, 'controller.default_sort_column.property');
                    } catch (err) {
                        console.log("No defautl sort column property")
                    }
                    console.warn('the column "' + get(this, 'controller.default_sort_column.property') + '" was not found.');
                    return column;
                }
                return column;
            }.property('controller.default_sort_column.property', 'fields.[]'),


            /**
             * @method updateAlarm
             * @description Update an alarm after performaning an action
             */
            updateRecord: function (alarmId) {
                var controller = this;
                var alarm_record = this.get('alarms').findBy('id', alarmId);
                if (alarm_record) {
                    var self = this;
                    var filterState = this.get('model.alarms_state_filter.state') || 'resolved';
                    var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
                    var f = {
                        'd': alarm_record.get('entity_id')
                    }
                    adapter.findQuery('alarm',
                        {
                            lookups: JSON.stringify(["pbehaviors"]),
                            'filter': ('{"$or":[{"_id":"' + alarm_record.get('id') + '"}]}')
                        }).then(function (found_alarm) {
                            if (found_alarm.success) {
                                var fields = self.get('fields');
                                var alarm = found_alarm.data[0].alarms[0];
                                alarm_record.entity_id = alarm._id;
                                alarm_record.d = alarm._id;
                                alarm_record.set('extra_details', Ember.Object.create());
                                controller.get('extraDeatialsEntities').forEach(function (item) {
                                    alarm_record.set('extra_details.' + item.name, Ember.Object.create(alarm).get(item.value));
                                })

                                var newAlarm = Ember.Object.create();

                                fields.forEach(function (field) {
                                    if (field.humanName != 'extra_details') {
                                        var val = get(Ember.Object.create(alarm), field.getValue);
                                        alarm_record.set(field.humanName, val);
                                    }

                                });
                                alarm_record.set('isSelected', false);
                                alarm_record.set('isExpanded', false);
                                alarm_record.set('canceled', get(Ember.Object.create(alarm), 'v.canceled'));

                                Ember.set(alarm_record, 'changed', new Date().getTime());

                            } else {
                                console.error('unsuccessful request');
                            }
                        })
                } else {
                    console.error('alarm not found');
                }

            },

            createMultiplePBehaviors: function(selection) {
                if(selection.length === 0) {
                    return;
                }
                var obj = Ember.Object.create({'crecord_type': 'pbehaviorform'});
                var confirmform = formsUtils.showNew('modelform', obj, {
                    title: 'Put a pbehavior on these elements ?'
                });
                this.stopRefresh();
                confirmform.submit.then(function(form) {
                    if(!Array.isArray(selection)) {
                        selection = [selection];
                    }

                    var listId = [];

                    for (var i = 0, l = selection.length; i < l; i++)
                        listId.push(selection[i]['d']);

                    var payload = {
                        'tstart':form.get('formContext.start'),
                        'tstop': form.get('formContext.type_') === 'Pause' ? 2147483647 : form.get('formContext.end'),
                        'rrule':form.get('formContext.rrule'),
                        'name':form.get('formContext.name'),
                        'type_': form.get('formContext.type_'),
                        'reason': form.get('formContext.reason'),
                        'author': window.username,
                        'filter': {}
                    };

                    if(!payload.rrule) {
                        delete(payload.rrule);
                    }

                    payload.filter = {
                        '_id': {
                            '$in': listId
                        }
                    }

                    //$.post(url)
                    return $.ajax({
                        type: 'POST',
                        url: '/api/v2/pbehavior',
                        data: JSON.stringify(payload),
                        contentType: 'application/json',
                        dataType: 'json',
                        success: function () {
                            console.log('pbehavior is sent');
                        },
                        statusCode: {
                            500: function () {
                                console.error("Failure to send pbehavior");
                            }
                        }
                    });

                });
            },

            actions: {
                massAction: function (action) {
                    this.sendEventCustom(action.mixin_name);
                },

                sendCustomAction: function (action, alarm) {
                    this.sendEventCustom(action.mixin_name, alarm);
                },

                updateSortField: function (field) {
                    this.set('alarmSearchOptions.sort_key', field.name);
                    this.set('alarmSearchOptions.sort_dir', field.isASC ? 'ASC' : 'DESC');
                },

                search: function (text) {
                    var controller = this;
                    controller.set('isValidSearchText', true);
                    controller.set('alarmSearchOptions.search', text);
                    controller.set('manualUpdateAlarms', new Date().getTime());
                },
            }

        }, listOptions);

        application.register('widget:listalarm', widget);
    }
});
