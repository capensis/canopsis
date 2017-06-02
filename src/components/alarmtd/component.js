Ember.Application.initializer({
    name: 'component-alarmtd',
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
            renderers: ['v_state', 'v_state_t', 'v_status', 'v_extra_details', ],
            // renderers: ['v_state_val'],

            init: function() {
                this._super();
              },

            click: function () {
                this.sendAction('action', this.get('alarm'), this.get('field'));
            },

            value: function() {
                var alarm = get(this, 'alarm');
                var field = get(this, 'field');
                var val = alarm[field.humanName];
                return val;
            }.property('alarm.changed', 'field'),

            hasRenderer: function () {
                return this.get('renderers').includes(this.get('field.name').replace(/\./g, "_"))
            }.property('alarm'),

            renderer: function () {
                return this.get('field.name').replace(/\./g, "_")
            }.property('alarm.name'),

            // vv: function() {
            //     return Ember.View.extend({
            //         // template: Ember.HTMLBars.compile("{{#component-rendererstate value=2}}{{/component-rendererstate}}")
            //         template: Ember.HTMLBars.compile("renderer")
                    
            //     });
            // }.property('field'),
        });

        application.register('component:component-alarmtd', component);
    }
});