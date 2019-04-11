Ember.Application.initializer({
    name: 'component-rendererstatus',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the rendererstatus component for the widget listalarm
         *
         * @class rendererstatus component
         */
        var component = Ember.Component.extend({

            /**
             * @property list
             */
            list: {
                0: 'Off',
                1: 'On going',
                2: 'Stealthy',
                3: 'Bagot',
                4: 'Canceled'
            },

            /**
             * @method init
             */
            init: function() {
                this._super();
            },

            /**
             * @property status
             */
            status: function() {
                return this.get('list')[this.get('value.val')];
            }.property('value.val'),

            /**
             * @property isCanceled
             */
            isCanceled: function() {
                return this.get('value.val') == 4
            }.property('value.val'),
        });

        application.register('component:component-rendererstatus', component);
    }
});