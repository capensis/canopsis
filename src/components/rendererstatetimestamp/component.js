Ember.Application.initializer({
    name: 'component-rendererstatetimestamp',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            init: function() {
                this._super();
            },

            timestamp: function() {
                return this.get('value');
                // return new Date(2017, 0, 30).getTime() / 1000; 
            }.property('value'),

            dateFormat: function (date) {
                var dDate = new Date(date);
                return moment(dDate).format('MM/DD/YY hh:mm:ss');
            },

        });

        application.register('component:component-rendererstatetimestamp', component);
    }
});