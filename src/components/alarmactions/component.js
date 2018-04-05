Ember.Application.initializer({
    name: 'component-alarmactions',
    initialize: function(container, application) {
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
                    }

					if (name != null) {
						toRemove.push(name)
					}
                }
				for(i = 0; i < toRemove.length; i++){
						console.debug("Cleaning ", name)
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
             * @property isSnoozed
             */
            isSnoozed: function () {
                return this.get('alarm.extra_details.snooze') != undefined;
            }.property('alarm.extra_details.snooze'),

            /**
             * @property hasLinks
             */
            hasLinks: function() {
                return this.get('alarm.linklist.event_links.length') > 0;
            }.property('alarm.linklist.event_links'),

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

            actions: {
                /**
                 * @method sendAction
                 */
                sendAction: function (action) {
                    this.sendAction('action', action, this.get('alarm'));
                }
            }
        });
        application.register('component:component-alarmactions', component);
    }
});
