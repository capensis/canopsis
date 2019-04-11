Ember.Application.initializer({
    name: 'component-rendererstatetimestamp',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the rendererstatetimestamp component for the widget listalarm
         *
         * @class rendererstatetimestamp component
         */

        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super();
            },

            /**
             * @property timestamp
             */
            timestamp: function() {
                return this.get('value');
                // return new Date(2017, 2, 20).getTime() / 1000; 
            }.property('value'),

            /**
             * @method dateFormat
             */
            dateFormat: function (date) {
                var dDate = new Date(date);
                return moment(dDate).format('MM/DD/YY hh:mm:ss');
            },

        });

        application.register('component:component-rendererstatetimestamp', component);
    }
});