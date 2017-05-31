Ember.Application.initializer({
    name: 'component-alarmactions',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the eventcategories component for the widget calendar
         *
         * @class eventcategories component
         * @memberOf canopsis.frontend.brick-calendar
         */
        var component = Ember.Component.extend({
            tagName: 'td',
            classNames: ['action-cell'],
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
                // teeeeeeeeeeeeeeeeeeeest
                {
                    class: 'fa fa-calendar-o',
                    internal_states: ['immutable' ,'unacked', 'acked', 'cancelled'],
                    name: 'pbehavior',
                    mixin_name: 'pbehavior'
                },
                // {
                //     class: 'glyphicon glyphicon-trash',
                //     internal_states: ['acked'],
                //     name: 'cancelalarm',
                //     mixin_name: 'cancel'
                // },
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
                // TODO declare incident
            ]),
            init: function() {
                this._super();
              },
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
            isAcked: function() {
                return this.get('alarm.extra_details.ack') != undefined;
            }.property('alarm.changed', 'alarm.extra_details'),
//             isRemoved: function() {
//                 return false;
// // TODO temporary solution
//             }.property('alarm.ack._t'),
            isCanceled: function () {
                return this.get('alarm.canceled') != undefined;
            }.property('alarm.canceled'),
//             isRestored: function () {
//                 return false;
// // TODO temporary solution
//             }.property('alarm.ack._t'),
            isSnoozed: function () {
                return this.get('alarm.extra_details.snooze') != undefined; 
            }.property('alarm.extra_details.snooze'),
            hasLinks: function() {
                return this.get('alarm.linklist.event_links.length') > 0;
            }.property('alarm.linklist.event_links'),

            isChangedByUser: function () {
                return this.get('alarm.state._t') == 'changestate'
            }.property('alarm.state._t'),

            actions: {
                sendAction: function (action) {
                    this.sendAction('action', action, this.get('alarm'));
                }
            }
        });
        application.register('component:component-alarmactions', component);
    }
});