Ember.Application.initializer({
    name: 'component-rendererpbehaviors',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;
            rRule = window.RRule;
        /**
         * This is the rendererpbehaviors component for the widget listalarm
         *
         * @class rendererpbehaviors component
         */

        var component = Ember.Component.extend({
            
            /**
             * @property propertyMap
             */
            propertyMap: {
                'tstop': 'stop date',
                'enabled': 'enabled',
                'name': 'name',
                'tstart': 'start date',
                'rrule': 'rule'
            },

            /**
             * @property properties
             */
            properties: ['tstop', 'enabled', 'name', 'tstart', 'rrule'],

            /**
             * @property pbehMap
             */
            pbehMap: function () {
                var propertyMap = this.get('propertyMap');
                var val = this.get('value');
                var self = this;

                return val.map(function(pbeh) {
                    return pbeh.name + ' - start: ' + self.dateFormat(pbeh.tstart) + ' - stop: ' + self.dateFormat(pbeh.tstop) + ' - freq: ' + self.rruleParse(pbeh.rrule);
                })
            }.property('properties', 'propertyMap', 'value'),

            /**
             * @method dateFormat
             */
            dateFormat: function (date) {
                var dDate = new Date(date);
                return moment(dDate).format('MM/DD/YY hh:mm:ss');
            },

            /**
             * @method rruleParse
             */
            rruleParse: function (rrule) {
                var params = rRule.parseString(rrule);
                var rule = new rRule(params || {});
                return rule.toText();
            },

            /**
             * @method init
             */
            init: function() {
                this._super();
            },
            
        }); 

        application.register('component:component-rendererpbehaviors', component);
    }
});