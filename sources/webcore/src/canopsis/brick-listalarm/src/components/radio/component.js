Ember.Application.initializer({
    name: 'component-radio',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the radio component for the widget listalarm
         *
         * @class radio component
         */
        var component = Ember.Component.extend({
            tagName: 'input',
            type : "radio",
            attributeBindings : [ "name", "type", "value", "checked:checked:" ],

            /**
             * @method init
             */
            init: function() {
                this._super();
            },

            /**
             * @method click
             */
            click : function() {
                this.set("selection", this.$().val())
            },
            
            /**
             * @property checked
             */
            checked : function() {
                return this.get("value") == this.get("selection");   
            }.property()

        });

        application.register('component:component-radio', component);
    }
});