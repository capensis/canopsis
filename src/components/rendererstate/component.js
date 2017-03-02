Ember.Application.initializer({
    name: 'component-rendererstate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            list: {
                0: {color: 'bg-green', name: 'Info'},
                1: {color: 'bg-yellow', name: 'Minor'},
                2: {color: 'bg-orange', name: 'Major'},
                3: {color: 'bg-red', name: 'Critical'}            
            },

            init: function() {
                this._super();
              },

            spanClass: function() {
                return [this.get('list')[this.get('value')]['color'], 'badge'].join(' ')
            }.property('value'),

            caption: function() {
                return this.get('list')[this.get('value')]['name']
            }.property('value'),

        });

        application.register('component:component-rendererstate', component);
    }
});