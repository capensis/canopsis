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
    name:'CustomSendeventMixin',
    after: ['MixinFactory', 'FormsUtils', 'DatesUtils', 'NotificationUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');
        var formsUtils = container.lookupFactory('utility:forms');
        var datesUtils = container.lookupFactory('utility:dates');
        var notificationUtils = container.lookupFactory('utility:notification');
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;
        /**
         * @mixin sendevent
         */
        var mixin = Mixin('customsendevent', {
            needs: ['application'],
            /**
             * @method init
             */
            init: function (){
                this._super();
                set(this, 'login', get(this, 'controllers.login.record'));
            },
            partials: {
                itemactionbuttons: [
                    'actionbutton-ack',
                    'actionbutton-cancel',
                    'actionbutton-incident',
                    'actionbutton-changestate',
                    'actionbutton-ticketnumber'
                ],
                selectionToolbarButtons: ['actionbutton-cancelselection', 'actionbutton-ackselection']
            },
            /**
             * @method setPendingOperation
             * @argument {array} crecords
             *
             * @description
             * Change crecord value to make it pending in user interface, this allows displaying a loadding glyphicon.
             * Pending status is defined for records that undergoes change and witch are not yet refreshed from server.
             */
            setPendingOperation: function(crecords){
                var safe_mode = get(this, 'controllers.application.frontendConfig.safe_mode');
                if (safe_mode) {
                    for (var i = 0, l = crecords.length; i < l; i++) {
                        console.log('Pending operations on crecord',crecords[i]);
                        // !!!!
                        crecords[i]['pendingOperation'] = true;
                    }
                }
            },
            /**
             * @method getDataFromRecord
             * @argument event_type
             * @argument crecord
             * @argument formRecord
             *
             * @description
             * Generates an object that contains all neceary data to be an event understood by Canopsis
             */
            getDataFromRecord: function(event_type, crecord, formRecord) {
                console.group('getDataFromRecord');
                console.log('get data from record:', event_type, formRecord);
                var record = {
                    // authkey: get(this, 'login.authkey'),
                    author: get(this, 'login._id'),
                    id: get(crecord, 'id'),
                    connector: get(crecord, 'connector'),
                    connector_name: get(crecord, 'connector_name'),
                    event_type: event_type,
                    source_type: get(crecord, 'source_type'),
                    component: get(crecord, 'component'),
                    state: get(crecord, 'state.val'),
                    state_type: get(crecord, 'state_type'),
                    crecord_type: event_type,
                    timestamp:  datesUtils.getNow()
                };
                var extra_fields = ['domain', 'perimeter', 'resource', 'output'];
                for (var i = 0, l = extra_fields.length; i < l; i++) {
                    var field = extra_fields[i];
                    if (!isNone(get(crecord, field))) {
                        set(record, field, get(crecord, field));
                    }
                }
                if(record.source_type === 'resource') {
                    record.resource = get(crecord, 'resource');
                }
                this.processEvent(event_type, 'extract', [record, crecord, formRecord]);
                console.groupEnd();
                return record;
            },
            /**
             * @method submitEvents
             * @argument crecords
             * @argument record
             * @argument event_type
             *
             * @description
             * Sends event to the api that will then send event to the amqp service.
             * Events are built from UI elements
             */
            submitEvents: function(crecords, record, event_type) {
                var me = this;
                var store = get(this, 'widgetDataStore') || Ember['widgetDataStore'];
                var safe_mode = this.get('controllers.application.frontendConfig.safe_mode');
                console.log('safe_mode', safe_mode);
                return new Ember.RSVP.Promise(function(resolve, reject) {
                    void(reject);
                    var post_events = [];
                    for(var i = 0; i < crecords.length; i++) {
                        console.log('Event:', get(record, 'author'), get(record, 'output'));
                        var post_event = me.getDataFromRecord(event_type, crecords[i], record);
                        post_event.author = get(record, 'author');
                        post_event.output = get(record, 'output');
                        if(get(record, 'ticket')) {
                            post_event.ticket = get(record, 'ticket');
                        }
                        if(get(record, 'ack_resources')) {
                            post_event.ack_resources = get(record, 'ack_resources');
                        }
                        post_events.push(post_event);
                        //processes ui event as they are processed by backoffice in order to get direct render
                        if (!safe_mode) {
                            me.processEvent(event_type, 'transform', [crecords[i], post_event]);
                        }
                        // in order to prevent an error when a user ack and fast ack on the same alarm
                        store.destroy(record);

                        me.refreshContent();
                    }
                    //TODO use an adapter for this
                    $.post('/event', {
                        event: JSON.stringify(post_events)
                    }).then(function(response) {
                        if (response.success) {
                            // me.updateAlarm(data.data[0].id.replace(/\..*/, ''));
                            for(var i = 0; i < response.data.length; i++) {
                                store.destroy(response.data[i]);
                            }
                            me.refreshContent();

                        } else {
                            console.error('error while send event', response.data.msg);
                        }
                    });
                });
            },
            /**
             * @method getEventForm
             * @argument event_type
             * @argument record
             * @argument crecords
             * @argument formbuttons
             *
             * @description
             * Generates and displays a form for givent record type
             */
            getEventForm: function(event_type, record, crecords, formType) {
                formType = formType || 'modelform';
                var wizard = formsUtils.showNew(formType, record, {
                    title: __('Add event type: ') + event_type,
                    override_labels: {
                        output: 'comment'
                    },
                    onePageDisplay: true
                });
                var me = this;
                var rollback = function() {
                    me.startRefresh();
                    record.rollback();
                    record.unloadRecord();
                };
                wizard.submit.then(function(form) {
                    console.log('saveRecord:', record, form);
                    record = get(form, 'formContext');
                    notificationUtils.info(__('event sent: ') + event_type);
                    me.submitEvents(crecords, record, event_type);
                    //send a delayed declare ticket when the option reportIncident is matched
                    if ($.inArray('reportIncident', arguments[1]) !== -1) {
                        me.submitEvents(crecords, record, 'declareticket');
                    } else if ($.inArray('doneTicket', arguments[1]) !== -1) {
                        set(record, 'output', '')
                        me.submitEvents(crecords, record, 'done');
                    }
                    rollback();
                }, rollback);
            },
            /**
             * @method getRoutingKey
             * @argument record
             * @return rk
             *
             * @description
             * Generates the routing key for given record.
             * Assume that record is an event.
             */
            getRoutingKey: function(record) {
                var rk = [
                    record.connector,
                    record.connector_name,
                    record.event_type,
                    record.source_type,
                    record.component
                ].join('.');
                if (record.source_type === 'resource') {
                    rk = [record.id, record.resource].join('.');
                }
                return rk;
            },
            /**
             * @method getDisplayRecord
             * @argument event_type
             * @argument crecord
             * @return record
             */
            getDisplayRecord: function(event_type, crecord) {
                var store = get(this, 'widgetDataStore') || Ember['widgetDataStore'];
                var recordData = this.getDataFromRecord(event_type, crecord);
                var record = store.createRecord(event_type, recordData);
                return record;
            },
            /**
             * @method filterUsableCrecords
             * @argument event_type
             * @argument crecords
             * @return ???
             *
             * @description
             * define if a selected record in list (with checkbox) is allowed
             * for the sendevent action for current event type
             */
            filterUsableCrecords: function(event_type, crecords) {
                var selected = [];
                for(var i = 0, l = crecords.length; i < l; i++) {
                    var record = crecords[i];
                    var topush = this.processEvent(event_type, 'filter', [record]);
                    if(topush) {
                        selected.pushObject(record);
                    }
                }
                return selected;
            },
            //TODO refactor into sub classes
            /**
             * @property event_processors
             *
             * @description
             * Transform keys makes transformation on crecord in order to simulate server data is fetched.
             */
            event_processors: {
                ack: {
                    /**
                     * @method ack_extract
                     * @param record
                     * @param crecords
                     * @param formRecord
                     * @description:
                     */
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        // record.id = this.getRoutingKey(record);
                        if (formRecord !== undefined) {
                            record.ticket = get(formRecord, 'ticket');
                            record.output = get(formRecord, 'output');
                        }
                    },
                    /**
                     * @method ack_filter
                     * @param record
                     * @description:
                     */
                    filter: function(record) {
                        var BAGOT = 3,
                            OFF = 0;
                        return (get(record, 'state') && !get(record, 'ack.isAck') && !get(record, 'ack.isCancel')
                            || (get(record, 'state') === OFF && get(record, 'status') === BAGOT));
                    },
                    /**
                     * @method ack_handle
                     * @param crecords
                     * @description:
                     */
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('ack', crecords[0]);
                        this.getEventForm('ack', record, crecords, 'ackform');
                    },
                    /**
                     * @method ack_transform
                     * @param crecords
                     * @param record
                     * @description:
                     */
                    transform: function(crecord, record) {
                        console.log('transform method for ack -> crecords', crecord, 'record', record);
                        crecord.set('ack', {
                            comment: record.output,
                            timestamp: datesUtils.getNow(),
                            author: record.author,
                            isAck: true
                        });
                        if(!isNone(record.ticket)) {
                            crecord.set('ticket', record.ticket);
                            crecord.set('ticket_date', datesUtils.getNow());
                        }
                        crecord.set('ack_remove', undefined);
                    }
                },
                fastack: {
                    /**
                     * @method fastack_extract
                     * @param record
                     * @param crecords
                     * @param formRecord
                     * @description:
                     */
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        // record.id = this.getRoutingKey(record);
                        record.ticket = 'ticket';
                        record.output = 'auto ack by ' + get(record, 'author');
                    },
                    /**
                     * @method fastack_filter
                     * @param record
                     * @description:
                     */
                    filter: function(record) {
                        var BAGOT = 3,
                            OFF = 0;
                        return (get(record, 'state') && !get(record, 'ack.isAck') && !get(record, 'ack.isCancel')
                            || (get(record, 'state') === OFF && get(record, 'status') === BAGOT));
                    },
                    /**
                     * @method fastack_handle
                     * @param crecords
                     * @description:
                     */
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('ack', crecords[0]);
                        var fastackmsg = get(this, 'mixinOptions.sendevent.fastackmsg');
                        record.set('output', 'auto ack by ' + get(record, 'author'))
                        this.submitEvents(crecords, record, 'ack');
                    },
                    /**
                     * @method fastack_transform
                     * @param crecords
                     * @param record
                     * @description:
                     */
                    transform: function(crecord, record) {
                        console.log('transform method for ack -> crecords', crecord, 'record', record);
                        crecord.set('ack', {
                            comment: record.output,
                            timestamp: datesUtils.getNow(),
                            author: record.author,
                            isAck: true
                        });
                        if(!isNone(record.ticket)) {
                            crecord.set('ticket', record.ticket);
                            crecord.set('ticket_date', datesUtils.getNow());
                        }
                        crecord.set('ack_remove', undefined);
                    }
                },
                ackremove: {
                    extract: function(record, crecord, formRecord) {
                        if (!isNone(formRecord)) {
                            record.output = get(formRecord, 'output');
                        }
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        // record.id = this.getRoutingKey(record);
                    },
                    filter: function(record) {
                        var isCanceled = get(record, 'canceled') != undefined;
                        var isAcked = get(record, 'extra_details.ack') != undefined;
                        return !isCanceled && isAcked;
                        // return (get(record, 'ack.author') && get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('ackremove', crecords[0]);
                        this.getEventForm('ackremove', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for ack remove', crecord, record);
                        crecord.set('ack', undefined);
                        crecord.set('declare_ticket_author', undefined);
                        crecord.set('declare_ticket_date', undefined);
                        crecord.set('ticket', undefined);
                        crecord.set('ticket_date', undefined);
                        crecord.set('ack_remove', Ember.Object.create({
                            comment: record.output,
                            timestamp: datesUtils.getNow(),
                            author: record.author
                        }));
                    }
                },
                done: {
                    extract: function(record, crecord, formRecord) {
                        if (formRecord === undefined) {
                            record.output = '';
                        } else {
                            record.output = get(formRecord, 'output');
                        }
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('done', crecords[0]);
                        this.getEventForm('done', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for done', crecord, record);
                        crecord.set('output', record.get('output'));
                        //crecord.set('done_date', datesUtils.getNow());
                    }
                },
                pbehavior: {
                    extract: function(record, crecord, formRecord) {
                        if (!isNone(formRecord)) {
                            record.output = get(formRecord, 'output');
                            // record.dtstart = get(formRecord, 'dtstart');
                            // record.dtstop = get(formRecord, 'dtstop');
                            record.start = get(formRecord, 'start');
                            record.end = get(formRecord, 'end');
                            record.rrule = get(formRecord, 'rrule');
                            record.pbehavior_name = 'downtime';
                            record.action = 'create';
                            record.name = get(formRecord, 'name');
                        }
                        // record.pbehavior_name = 'downtime';

                        // record.filter = {};
                        // record.comments = '';
                        // record.enabled = true;
                        // record.eids = [];
                        // record.
                    },
                    filter: function(record) {
                        return (get(record, 'ack.author') && get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('pbehavior', crecords[0]);
                        this.getEventForm('pbehavior', record, crecords, 'pbehaviorform');
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for ack remove', crecord, record, 'pbehavior');
                        // crecord.set('ack', undefined);
                        // crecord.set('declare_ticket_author', undefined);
                        // crecord.set('declare_ticket_date', undefined);
                        // crecord.set('ticket', undefined);
                        // crecord.set('ticket_date', undefined);
                        // crecord.set('ack_remove', Ember.Object.create({
                        //     comment: record.output,
                        //     timestamp: datesUtils.getNow(),
                        //     author: record.author
                        // }));
                    }
                },
                declareticket: {
                    extract: function(record, crecord, formRecord) {
                        void(formRecord);
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        // record.id = this.getRoutingKey(record);
                        record.output = 'declare ticket';
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('declareticket', crecords[0]);
                        this.getEventForm('declareticket', record, crecords, 'ticketform');
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for declare ticket', crecord, record);
                        crecord.set('ticket_declared_author', record.author);
                        crecord.set('ticket_declared_date', datesUtils.getNow());
                        if(!isNone(record.done)) {
                            crecord.set('done', record.done);
                            crecord.set('done_date', datesUtils.getNow());
                        }
                        //crecord.set('done_remove', undefined);
                    }
                },
                assocticket: {
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        // record.id = this.getRoutingKey(record);
                        if (formRecord === undefined) {
                            record.output = __('Associated ticket number');
                        } else {
                            record.output = get(formRecord, 'output');
                            record.ticket = get(formRecord, 'ticket');
                        }
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('assocticket', crecords[0]);
                        this.getEventForm('assocticket', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for assoticket', crecord, record);
                        crecord.set('ticket', record.ticket);
                        crecord.set('ticket_date', datesUtils.getNow());
                    }
                },
                cancel: {
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        if (formRecord !== undefined) {
                            record.output = get(formRecord, 'output');
                        }
                        record.cancel = 1;
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('cancel', crecords[0]);
                        this.getEventForm('cancel', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for cancel -> crecords', crecord, 'record', record);
                        crecord.set('ack.isCancel',true);
                        crecord.set('ack.isAck',false);
                        crecord.set('status', 4);
                        crecord.set('cancel',{
                            comment: record.output,
                            timestamp: datesUtils.getNow(),
                            author: record.author,
                            previous_status: record.state
                        });
                    }
                },
                cancelack: {
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state = 0;
                        record.state_type = 1;
                        if (formRecord !== undefined) {
                            record.output = get(formRecord, 'output');
                        }
                        record.cancel = 1;
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isAck'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('cancel', crecords[0]);
                        this.getEventForm('cancel', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for cancel -> crecords', crecord, 'record', record);
                        crecord.set('ack.isCancel',true);
                        crecord.set('ack.isAck',false);
                        crecord.set('status', 4);
                        crecord.set('cancel',{
                            comment: record.output,
                            timestamp: datesUtils.getNow(),
                            author: record.author,
                            previous_status: record.state
                        });
                    }
                },
                recovery: {
                    extract: function(record, crecord, formRecord) {
                        void(formRecord);
                        void(record);
                        set(crecord, 'state', 0);
                    },
                    filter: function(record) {
                        var isCanceled = get(record, 'canceled') != undefined;
                        return isCanceled;
                    },
                    handle: function(crecords) {
                        var record = crecords[0];
                        this.submitEvents([record], record, 'recovery');
                    },
                    transform: function(crecord, record) {
                        void(crecord);
                        void(record);
                        //TODO
                    }
                },
                uncancel: {
                    extract: function(record, crecord, formRecord) {
                        record.ref_rk = get(crecord, 'entity_id');
                        record.state_type = 1;
                        record.state = 0;
                        record.cancel = 0;
                        if (formRecord !== undefined) {
                            var output = get(formRecord, 'output');
                            if (! output) {
                                output = ' ';
                            }
                            record.output = output;
                        }
                    },
                    filter: function(record) {
                        return (get(record, 'ack.isCancel'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('uncancel', crecords[0]);
                        this.getEventForm('uncancel', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for uncancel -> crecords', crecord, 'record', record);
                        crecord.set('ack.isCancel', false);
                        crecord.set('ack.isAck', true);
                        crecord.set('status', crecord.get('cancel.previous_status'));
                        //reset the ack is a hack if ack is not set in the event, but there is no choice and this is a temp information
                        if(isNone(crecord.get('ack'))) {
                            crecord.set('ack', {
                                comment: record.output,
                                timestamp: datesUtils.getNow(),
                                author: record.author,
                                isAck: true
                            });
                        } else {
                            crecord.set('ack', crecord.get('ack'));
                        }
                    }
                },
                changestate: {
                    extract: function(record, crecord, formRecord) {
                        if(!isNone(formRecord)) {
                            record.state = get(formRecord, 'state');
                            record.output = get(formRecord, 'output');
                        }
                        // diff with the previous widget
                        record.event_type = 'changestate';
                        record.keep_state = true;
                        record.state_type = 1;
                    },
                    filter: function(record) {
                        return (get(record, 'state'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('changestate', crecords[0]);
                        this.getEventForm('changestate', record, crecords);
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for ack changestate', crecord, record);
                        crecord.set('state', record.state);
                        if (record.state === 0) {
                            crecord.set('ack', undefined);
                        } else {
                            crecord.set('cancel', undefined);
                        }
                        set(crecord, 'change_state_output', get(record, 'output'));
                        set(crecord, 'keep_state', true);
                    }
                },
                // TODO change, currently it is copy of changestate
                snooze: {
                    extract: function(record, crecord, formRecord) {
                        if(!isNone(formRecord)) {
                            record.state = get(formRecord, 'state');
                            record.output = get(formRecord, 'output');
                            record.duration = get(formRecord, 'duration');
                        }

                        // record.event_type = 'check';
                        // record.keep_state = true;
                        // record.state_type = 1;
                        // record.duration = 3;
                    },
                    filter: function(record) {
                        return (get(record, 'state'));
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('snooze', crecords[0]);
                        this.getEventForm('snooze', record, crecords, 'snoozeform');
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for ack snooze', crecord, record);
                        // crecord.set('state', record.state);
                        // if (record.state === 0) {
                        //     crecord.set('ack', undefined);
                        // } else {
                        //     crecord.set('cancel', undefined);
                        // }
                        // set(crecord, 'change_state_output', get(record, 'output'));
                        // set(crecord, 'keep_state', true);
                    }
                },
                comment: {
                    extract: function(record, crecord, formRecord) {
                        void(formRecord);
                        record.referer = get(crecord, 'referer');
                        record.output = get(crecord, 'output');
                        record.display_name = get(this, 'login.firstname') + ' ' + get(this, 'login.lastname');
                    },
                    filter: function(crecords) {
                        void(crecords);
                        return false;
                    },
                    handle: function(crecords) {
                        var record = this.getDisplayRecord('comment', crecords[0]);
                        notificationUtils.info(__('event "comment" sent'));
                        this.submitEvents(crecords, record, 'comment');
                    },
                    transform: function(crecord, record) {
                        console.log('transform method for comment', crecord, record);
                        //TODO
                    }
                }
            },
            /**
             * @property hasEventProcessorForType
             * @argument event_type
             * @return {boolean}
             *
             * @description
             * Boolean return method that tells whether or not an event type has it 's own processing code
             */
            hasEventProcessorForType: function(event_type) {
                return (this.event_processors[event_type] !== undefined);
            },
            /**
             * @method processEvent
             * @argument event_type
             * @argument fname
             * @argument args
             * @return the processor method return value
             *
             * @description
             * Event processing code handler
             */
            processEvent: function(event_type, fname, args) {
                if(this.hasEventProcessorForType(event_type)) {
                    var callback = this.event_processors[event_type][fname];
                    return callback.apply(this, args);
                }
            },
            actions: {
                /**
                 * @method actions-sendEvent
                 * @argument event_type
                 * @argument crecord
                 *
                 * @description
                 * Entry point for this class
                 * Sends one or many event to the server depending on selected record and action performed.
                 */
                sendEvent: function(event_type, crecord) {
                    console.group('sendEvent:', arguments);
                    this.stopRefresh();
                    var crecords = [];
                    if (!isNone(crecord)) {
                        console.log('event:', event_type, crecord);
                        crecords.pushObject(crecord);
                    }
                    else {
                        var content = get(this, 'widgetData.content');
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
                }
            }
        });
        application.register('mixin:customsendevent', mixin);
    }
});
