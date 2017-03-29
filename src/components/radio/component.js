Ember.Application.initializer({
    name: 'component-radio',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({
            tagName: 'input',
            type : "radio",
            attributeBindings : [ "name", "type", "value", "checked:checked:" ],

            init: function() {
                this._super();
            },

            click : function() {
                            // console.error(this);
                this.set("selection", this.$().val())
            },
            
            checked : function() {
                return this.get("value") == this.get("selection");   
            }.property()

        });

        application.register('component:component-radio', component);
    }
});