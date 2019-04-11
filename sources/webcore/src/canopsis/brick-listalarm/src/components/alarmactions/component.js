Ember.Application.initializer({
    name: 'component-alarmactions',
    initialize: function(container, application) {
        var formsUtils = container.lookupFactory('utility:forms');
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        __ = Ember.String.loc;

        /**
         * This is the alarmactions component for the widget listalarm
         *
         * @class alarmactions component
         */
        var component = Ember.Component.extend({
            tagName: 'td',
            classNames: ['action-cell'],

            /**
             * @property actionsMap
             */
            actionsMap: Ember.A([
                {
                    class: 'glyphicon glyphicon-saved',
                    internal_states: ['unacked'],
                    name: 'ack',
                    mixin_name: 'ack'
                },
                {
                    class: 'glyphicon glyphicon-ok',
                    internal_states: ['unacked'],
                    name: 'fastack',
                    mixin_name: 'fastack'
                },
                {
                    class: 'glyphicon glyphicon-ban-circle',
                    internal_states: ['acked'],
                    name: 'cancelack',
                    mixin_name: 'ackremove'
                },
                {
                    class: 'fa fa-ticket',
                    internal_states: ['acked', 'cancelled'],
                    name: 'declareanincident',
                    mixin_name: 'declareticket'
                },
                {
                    class: 'fa fa-thumb-tack',
                    internal_states: ['acked', 'cancelled'],
                    name: 'assignticketnumber',
                    mixin_name: 'assocticket'
                },
                {
                    class: 'fa fa-pause',
                    internal_states: ['immutable' ,'unacked', 'acked', 'cancelled'],
                    name: 'pbehavior',
                    mixin_name: 'pbehavior'
                },
                {
                    class: 'glyphicon glyphicon-trash',
                    internal_states: ['acked'],
                    name: 'removealarm',
                    mixin_name: 'cancelack'
                },
                {
                    class: 'fa fa-exclamation-triangle',
                    internal_states: ['acked'],
                    name: 'changecriticity',
                    mixin_name: 'changestate'
                },
                {
                    class: 'glyphicon glyphicon-share-alt',
                    internal_states: ['cancelled'],
                    name: 'restorealarm',
                    mixin_name: 'recovery'
                },
                {
                    class: 'fa fa-clock-o',
                    internal_states: ['unacked', 'acked', 'cancelled'],
                    name: 'snoozealarm',
                    mixin_name: 'snooze'
                },
                {
                    class: 'fa fa-check-square',
                    internal_states: ['acked', 'cancelled'],
                    name: 'donealarm',
                    mixin_name: 'done'
                },
                {
                    class: 'fa fa-list',
                    internal_states: ['immutable' ,'unacked', 'acked', 'cancelled'],
                    name: 'listpbehavior',
                    mixin_name: 'listpbehavior',
                    isPbhList: true
                }
            ]),

            /**
             * @method init
             */
            init: function() {
                this._super();
                this.set("rights", this.get("_parentView._parentView._parentView._parentView._parentView._controller.login.rights"))
                // Translating tooltips
                this.get('actionsMap').filter(function(item, index, enumerable) {
                        Ember.set(item, 'translation', __(item.name));
                });
            },

            genRemoveList : function(actions, rights, name, right){
                if (rights.hasOwnProperty(right)) {
                    if (rights.get(right).checksum){
                        return null
                    }
                }
				return name
            },

            /**
             * @property availableActions
             */
            availableActions: function() {
                var intState = this.get('internalState');
                var actions = this.get('actionsMap').filter(function(item, index, enumerable) {
                    return item.internal_states.includes(intState)
                });
                if (this.get('isSnoozed')) {
                    actions.removeObject(actions.findBy('mixin_name', 'snooze'))
                };
                if (this.get('isClosed')) {
                    actions.removeObject(actions.findBy('mixin_name', 'pbehavior'))
                }

                var rights = this.get("rights");
				var toRemove = []
                for (var i = 0; i < actions.length; i++) {
                    var func = this.get("genRemoveList")
					var name = null
                    switch(actions[i]["name"]){
                    case "ack":
                        name = func(actions, rights, actions[i]["name"], "listalarm_ack")
                        break;
                    case "fastack":
                        name = func(actions, rights, actions[i]["name"], "listalarm_fastAck")
                        break;
                    case "cancelack":
                        name = func(actions, rights, actions[i]["name"], "listalarm_cancelAck")
                        break;
                    case "declareanincident":
                        name = func(actions, rights, actions[i]["name"], "listalarm_declareanIncident")
                        break;
                    case "assignticketnumber":
                        name = func(actions, rights, actions[i]["name"], "listalarm_assignTicketNumber")
                        break;
                    case "donealarm":
                        name = func(actions, rights, actions[i]["name"], "listalarm_done")
                        break;
                    case "pbehavior":
                        name = func(actions, rights, actions[i]["name"], "listalarm_pbehavior")
                        break;
                    case "removealarm":
                        name = func(actions, rights, actions[i]["name"], "listalarm_removeAlarm")
                        break;
                    case "changestate":
                        name = func(actions, rights, actions[i]["name"], "listalarm_changeState")
                        break;
                    case "changecriticity":
                        name = func(actions, rights, actions[i]["name"], "listalarm_changeCriticity")
                        break;
                    case "restorealarm":
                        name = func(actions, rights, actions[i]["name"], "listalarm_restoreAlarm")
                        break;
                    case "snoozealarm":
                        name = func(actions, rights, actions[i]["name"], "listalarm_snoozeAlarm")
                        break;
                    case "listpbehavior":
                        name = func(actions, rights, actions[i]["name"], "listalarm_listPbehavior")
                        break;
                    }

                    if (name != null) {
                         toRemove.push(name)
                    }
                }
                for(i = 0; i < toRemove.length; i++){
                    actions.removeObject(actions.findBy('name', toRemove[i]))
                }
                return actions;
            }.property('internalState', 'isSnoozed', 'isChangedByUser', 'isClosed', "isAcked"),

            /**
             * @property internalState
             */
            internalState: function() {
                if (this.get('alarm.state.val') == 0) {
                    return 'immutable';
                }
                if (this.get('isCanceled')) {
                    return 'cancelled';
                }
                if (this.get('isAcked')) {
                    return 'acked';
                }
                if (this.get('isDone')) {
                    return 'done';
                }
                return 'unacked';
            }.property('alarm.state.val', 'isCanceled', 'isAcked'),

            /**
             * @property isAcked
             */
            isAcked: function() {
                return this.get('alarm.extra_details.ack') != undefined;
            }.property('alarm.changed', 'alarm.extra_details'),

            /**
             * @property isCanceled
             */
            isCanceled: function () {
                return this.get('alarm.canceled') != undefined;
            }.property('alarm.canceled'),

            /**
             * @property isDone
             */
            isDone: function () {
                return this.get('alarm.done') != undefined;
            }.property('alarm.done'),

            /**
             * @property isSnoozed
             */
            isSnoozed: function () {
                return this.get('alarm.extra_details.snooze') != undefined;
            }.property('alarm.extra_details.snooze'),

            /**
             * @property isChangedByUser
             */
            isChangedByUser: function () {
                return this.get('alarm.state._t') == 'changestate'
            }.property('alarm.state._t'),

            /**
             * @property isClosed
             */
            isClosed: function () {
                return this.get('alarm.state.val') == 0;
            }.property('alarm.state.val'),


            /**
             * @property pBehaviorListActionId
             * @description Forms a unique id for a modal form
             */
            pBehaviorListActionId: function() {
                var alarmId = this.get('alarm.id').replaceAll("-", "");

                return 'pBehaviorListAction' + alarmId;
            }.property('alarm.id'),

            /**
             * @property pBehaviorListActionIdHash
             * @description Adds a leading hash sing to pBehaviorListActionId
             */

            pBehaviorListActionIdHash: function() {
                return '#' + this.get('pBehaviorListActionId');
            }.property('pBehaviorListActionId'),

            actions: {
                /**
                 * @method sendAction
                 */
                sendAction: function (action) {
                    var me = this;
                    var alarm = Object.assign({}, me.get('alarm'));
                    if (action.name === 'pbehavior'){

                        var obj = Ember.Object.create({ 'crecord_type': 'pbehaviorform' });
                        var confirmform = formsUtils.showNew('modelform', obj, {
                            title: 'Put a pbehavior on these elements ?'
                        });

                        confirmform.submit.then(function (form) {

                            var payload = {
                                'tstart': form.get('formContext.start'),
                                'tstop': form.get('formContext.type_') === 'Pause' ? 2147483647 : form.get('formContext.end'),
                                'rrule': form.get('formContext.rrule'),
                                'name': form.get('formContext.name'),
                                'type_': form.get('formContext.type_'),
                                'reason': form.get('formContext.reason'),
                                'author': window.username,
                                'filter': {}
                            };

                            if (!payload.rrule) {
                                delete (payload.rrule);
                            }
                            payload.filter = {
                                '_id': {
                                    '$in': [alarm.d]
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
                        return
                    }

                    this.sendAction('action', action, alarm);
                },
            }
        });
        application.register('component:component-alarmactions', component);
    }
});
