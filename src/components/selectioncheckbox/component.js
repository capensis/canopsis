Ember.Application.initializer({
    name: 'component-selectioncheckbox',
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

            init: function() {
                this._super();
            },

        });

        application.register('component:component-selectioncheckbox', component);
    }
});