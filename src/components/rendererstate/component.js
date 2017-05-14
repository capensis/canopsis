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
                return [this.get('list')[this.get('value.val')]['color'], 'badge'].join(' ')
            }.property('value.val'),

            caption: function() {
                return this.get('list')[this.get('value.val')]['name']
            }.property('value.val'),

            isChangedByUser: function () {
                return this.get('value._t') == 'changestate'
            }.property('value._t'),

            isCanceled: function () {
                return this.get('value.canceled') != undefined;
            }.property('value.canceled'),

        });

        application.register('component:component-rendererstate', component);
    }
});