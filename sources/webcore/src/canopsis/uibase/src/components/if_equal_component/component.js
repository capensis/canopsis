Ember.Application.initializer({
    name: 'if_equal_component',
    initialize: function(container, application){
        var get = Ember.get, 
            set = Ember.set;

        var component = Ember.Component.extend({
            
            init: function() {
                this._super();
            },

            isEqual: function() {
                return (this.get('param1') == this.get('param2'));
            }.property('param1', 'param2')

        });

        application.register('component:component-if_equal_component', component)

    }
});

