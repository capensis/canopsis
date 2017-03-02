Ember.Application.initializer({
    name: 'component-rendererpbehaviors',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({

            propertyMap: {
                'tstop': 'stop date',
                'enabled': 'enabled',
                'name': 'name',
                'tstart': 'start date',
                'rrule': 'rule'
            },

            properties: ['tstop', 'enabled', 'name', 'tstart', 'rrule'],

            pbehMap: function () {
                var propertyMap = this.get('propertyMap');
                var val = this.get('value');
                var self = this;

                return val.map(function(pbeh) {
                    return pbeh.name + ' - start: ' + self.dateFormat(pbeh.tstart) + ' - stop: ' + self.dateFormat(pbeh.tstop) + ' - freq: ' + pbeh.rrule;
                })
            }.property('properties', 'propertyMap', 'value'),

            dateFormat: function (date) {
                var dDate = new Date(date);
                return moment(dDate).format('MM/DD/YY hh:mm:ss');
            },

            init: function() {
                this._super();
              },
            
        });

        application.register('component:component-rendererpbehaviors', component);
    }
});