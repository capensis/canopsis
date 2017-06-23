Ember.Application.initializer({
    name: 'component-rendererstate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the rendererstate component for the widget listalarm
         *
         * @class rendererstate component
         */

        var component = Ember.Component.extend({

            /**
             * @property list
             */
            list: {
                0: {color: 'bg-green', name: 'Info'},
                1: {color: 'bg-yellow', name: 'Minor'},
                2: {color: 'bg-orange', name: 'Major'},
                3: {color: 'bg-red', name: 'Critical'}            
            },

            /**
             * @method init
             */
            init: function() {
                this._super();
              },

            /**
             * @property spanClass
             */
            spanClass: function() {
                return [this.get('list')[this.get('value.val')]['color'], 'badge'].join(' ')
            }.property('value.val'),

            /**
             * @property caption
             */
            caption: function() {
                return this.get('list')[this.get('value.val')]['name']
            }.property('value.val'),

            /**
             * @property isChangedByUser
             */
            isChangedByUser: function () {
                return this.get('value._t') == 'changestate'
            }.property('value._t'),

            /**
             * @property isCancelled
             */
            isCancelled: function () {
                return this.get('value.canceled') != undefined;
            }.property('value.canceled'),
        });

        application.register('component:component-rendererstate', component);
    }
});