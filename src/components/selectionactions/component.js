Ember.Application.initializer({
    name: 'component-selectionactions',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the selectionactions component for the widget listalarm
         *
         * @class selectionactions component
         */
        var component = Ember.Component.extend({
            tagName: 'td',
            classNames: ['action-cell'],

            /**
             * @property actionsMap
             */
            actionsMap: Ember.A([
                {
                    class: 'glyphicon glyphicon-ok',
                    mixin_name: 'fastack',
                    caption: 'Fast Ack'
                },
                {
                    class: 'glyphicon glyphicon-saved',
                    mixin_name: 'ack',
                    caption: 'Ack'
                },
                {
                    class: 'glyphicon glyphicon-ban-circle',
                    mixin_name: 'ackremove',
                    caption: 'Cancel ack'
                },
                {
                    class: 'glyphicon glyphicon-share-alt',
                    mixin_name: 'recovery',
                    caption: 'Restore alarm'
                }
            ]),

            /**
             * @method init
             */
            init: function() {
                this._super();
            },

            actions: {
                /**
                 * @property sendAction
                 */
                sendAction: function (action) {
                    this.sendAction('action', action);
                }
            }

        });

        application.register('component:component-selectionactions', component);
    }
});