Ember.Application.initializer({
    name: 'component-pbehaviorslist',
    initialize: function(container, application) {

        /**
         * This is the pbehaviorslist component for the widget listalarm
         *
         * @class pbehaviorslist component
         */
        var component = Ember.Component.extend({

            /**
             * @method init
             */
            init: function() {
                this._super();
                
                this.set('pbehaviors', this.get('alarm.extra_details.pbehaviors'));
                this.set('extraDetailsComponent', this.get('alarm.extra_details.extraDetailsComponent'));
            },

            /**
             * @method rerenderExtraDetailsComponent
             * @description Rerenders the extraDetailed component in order to sync
             * a representation of pbehaviors and an actual state
             */

            rerenderExtraDetailsComponent: function() {
                var component = this.get('extraDetailsComponent');

                if (component) {
                    component.rerender();
                }
            },

            /**
             * @method deletePbehavior
             * @description Sends a request to the server to remove a pbehavior
             */

            deletePbehavior: function(pBehavior) {
                var me = this;
                $.ajax({
                    type: 'DELETE',
                    url: '/api/v2/pbehavior/' + pBehavior._id,
                    contentType: 'application/json',
                    dataType: 'json',
                    success: function () {
                        console.log('pbehavior removed');

                        me.get('pbehaviors').removeObject(pBehavior);
                        me.rerenderExtraDetailsComponent();
                    },
                    error: function () { console.error('Failure to send remove pbehavior'); }
                });
            },

            actions: {
                deletePbehavior: function(pBehavior) {
                    this.deletePbehavior(pBehavior);
                }
            }
        });
        application.register('component:component-pbehaviorslist', component);
    }
});
 