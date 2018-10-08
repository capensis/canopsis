Ember.Application.initializer({
    name: 'component-selectioncheckbox',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the selectioncheckbox component for the widget listalrms
         *
         * @class selectioncheckbox component
         */
        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super();
            },

        });

        application.register('component:component-selectioncheckbox', component);
    }
});