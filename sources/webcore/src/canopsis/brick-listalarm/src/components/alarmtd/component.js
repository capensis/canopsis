Ember.Application.initializer({
    name: 'component-alarmtd',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the alarmtd component for the widget listalarm
         *
         * @class eventcategories alarmtd
         */
        var component = Ember.Component.extend({
            tagName: 'td',

            /**
             * @property renderers
             */
            renderers: ['v_state', 'v_state_t', 'v_status', 'v_extra_details', 'v_creation_date', 'v_last_update_date', 'v_last_event_date', 'v_resolved', 'v_state_m', 'v_duration', 'v_current_state_duration', 'infos_hostgroups_value'],

            /**
             * @method init
             */
            init: function() {
                this._super();
              },

            /**
             * @mthod click
             */
            click: function () {
                this.sendAction('action', this.get('alarm'), this.get('field'));
            },

            /**
             * @property value
             */
            value: function() {
                var alarm = get(this, 'alarm');
                var field = get(this, 'field');
                var val = alarm[field.humanName];
                return val;
            }.property('alarm.changed', 'field'),

            /**
             * @property hasRenderer
             */
            hasRenderer: function () {
                return this.get('renderers').includes(this.get('field.name').replace(/\./g, "_")) || (this.get('field.name') === 'v.links')
            }.property('alarm'),

            /**
             * @property hasHelper
             */
            hasHelper: function () {
                return ['links'].includes(this.get('field.name'))
            }.property('alarm'),

            /**
             * @property renderer
             */
            renderer: function () {
                return this.get('field.name').replace(/\./g, "_")
            }.property('alarm.name'),
        });

        application.register('component:component-alarmtd', component);
    }
});
