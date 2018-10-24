Ember.Application.initializer({
    name: 'component-popupinfo',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
        /**
         * This is the popupinfo component for the widget listalarm
         *
         * @class popupinfo component
         */
        var component = Ember.Component.extend({

            /**
             * @method inti
             */
            init: function() {
                this._super();
                this.set('columnTemplate', Ember.columnTemplates.findBy('columnName', this.get('columnName')).columnTemplate)
            },

            /**
             * @method upt
             */
            upd: function () {
                if (this.get('columnName') == this.get('clickedField.humanName')) {
                    $('.popupinfo').hide();                                                
                    this.$('.popupinfo').fadeIn(500);                    
                }
            }.observes('updater'),

            actions: {
                /**
                 * @method hide
                 */
                hide: function () {
                    this.$('.popupinfo').fadeOut(500);
                }
            }

        });

        application.register('component:component-popupinfo', component);
    }
});