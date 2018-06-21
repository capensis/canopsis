Ember.Application.initializer({
    name: 'component-alarmraw',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the alarmraw component for the widget listalarm
         *
         * @class alarmraw component
         */
        var component = Ember.Component.extend({
            tagName: 'tr',

            /**
             * @method init
             */
            init: function() {
                this._super();
                set(this, 'alarm', get(this, 'alarm'));
                set(this, 'fields', get(this, 'fields'));
            },

            /**
             * @property expandClass
             */

            expandClass: function () {
                var cl = '';
                if (this.get('alarm.isExpanded')) {
                    cl = 'glyphicon-minus';
                } else {
                    cl = 'glyphicon-plus';
                }
                return 'glyphicon ' + cl;
            }.property('alarm.isExpanded'),

            actions: {
                /**
                 * @method tdClick
                 */
                tdClick: function (alarm, field) {
                    this.sendAction('action', alarm, field);
                },

                /**
                 * @method sendAction
                 */
                sendAction: function (action, alarm) {
                    this.sendAction('saction', action, alarm);
                },

                /**
                 * @method expand
                 */
                expand: function () {
                    this.toggleProperty('alarm.isExpanded');
                }
            }

        });

        application.register('component:component-alarmraw', component);
    }
});