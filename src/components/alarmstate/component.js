Ember.Application.initializer({
    name: 'component-alarmstate',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the alarmstate component for the widget listalarm
         *
         * @class alarmstate component
         */
        var component = Ember.Component.extend({
            /**
             * @property state
             */
            state: undefined,

            /**
             * @property isSelected
             */
            isSelected: 0,

            CheckBoxAlarmStateOpen: false,
            CheckBoxAlarmStatusResolved: true,

            /**
             * @method init
             */
            init: function() {
                this._super();
                if (!isNone(this.get('content'))) {
                    var state =  get(this, 'content.state')
                    set(this, 'state',state);
                    if(state == "opened"){
                        this.set("CheckBoxAlarmStateOpen",true)
                        this.set("CheckBoxAlarmStatusResolved",false)
                    }
                }
              },

            /**
             * @method onUpdate
             */
            onUpdate: function() {
                this.set('content', {
                    state: get(this, 'state'),
                });
            }.observes('state')

        });

        application.register('component:component-alarmstate', component);
    }
});
