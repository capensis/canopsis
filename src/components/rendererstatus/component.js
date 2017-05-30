Ember.Application.initializer({
    name: 'component-rendererstatus',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            list: {
                0: 'Off',
                1: 'On going',
                2: 'Stealthy',
                3: 'Bagot',
                4: 'Canceled'
            },

            init: function() {
                this._super();
            },

            status: function() {
                return this.get('list')[this.get('value.val')];
            }.property('value.val'),

            isCanceled: function() {
                return this.get('value.val') == 4
            }.property('value.val'),
        });

        application.register('component:component-rendererstatus', component);
    }
});