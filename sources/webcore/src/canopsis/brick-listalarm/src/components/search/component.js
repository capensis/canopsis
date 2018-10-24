Ember.Application.initializer({
    name: 'component-search',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the search component for the widget listalarm
         *
         * @class search component
         */

        var component = Ember.Component.extend({
            /**
             * @property classNames
             */
            classNames: ['col-xs-3', 'search'],

            /**
             * @method init
             */
            init: function() {
                this._super();
              },

            /**
             * @property removeInvalidSearchTextNotification
             */
            removeInvalidSearchTextNotification: function () {
                this.set('isValid', true)
            }.observes('value'),

            actions: {
                /**
                 * @method search
                 */
                search: function () {
                    if (this.get('value').length > 0) {
                        this.sendAction('action', this.get('value'));
                    }
                },

                /**
                 * @method resetValue
                 */
                resetValue: function () {
                    this.set('value', '');
                    this.sendAction('action', '');
                },
				openSearchDoc: function (){
					window.open("/api/v2/documentation/listalarms-search", "_blank")
				}
            }

        });

        application.register('component:component-search', component);
    }
});
