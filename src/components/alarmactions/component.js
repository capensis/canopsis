Ember.Application.initializer({
    name: 'component-alarmactions',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
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
                    class: 'fa fa-calendar-o',
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
              },

            /**
             * @property availableActions
             */
            availableActions: function() {
                var intState = this.get('internalState');
                // return this.get('actionsMap');
                var actions = this.get('actionsMap').filter(function(item, index, enumerable) {
                    return item.internal_states.includes(intState)
                });
                if (this.get('isSnoozed')) {
                    actions.removeObject(actions.findBy('mixin_name', 'snooze'))
                };
                if (this.get('isChangedByUser')) {
                    actions.removeObject(actions.findBy('mixin_name', 'cancelack'))                    
                }
                return actions;
            }.property('internalState', 'isSnoozed', 'isChangedByUser'),

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