Ember.Application.initializer({
    name: 'component-customtimeline',
    after: ['DataUtils', 'HashUtils'],
    initialize: function(container, application) {
        var dataUtils = container.lookupFactory('utility:data');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * This is the eventcategories component for the widget calendar
         *
         * @class eventcategories component
         * @memberOf canopsis.frontend.brick-calendar
         */
        var component = Ember.Component.extend({
            tagName: 'tr',
            

            classNameBindings: ['isHidden:invisible'],

            isHidden: function() {
                return !this.get('alarm.isExpanded');
            }.property('alarm.isExpanded'),

            timelineData: undefined,

            statusToName: {
                'ack': 'Acknowledged by ',
                'ackremove': 'Ack removed by ',
                'assocticket': 'Ticket association by ',
                'declareticket': 'Ticket declared by ',
                'cancel': 'Canceled by ',
                'uncancel': 'Restored by ',
                'statusinc': 'Status increased',
                'statusdec': 'Status decreased',
                'stateinc': 'State increased',
                'statedec': 'State decreased',
                'changestate': 'State changed',
                'snooze': 'Snoozed by '
            },

            stateArray: [
                'Ok',
                'Minor',
                'Major',
                'Critical'
            ],

            statusArray: [
                'off',
                'ongoing',
                'stealthy',
                'bagot',
                'canceled'
            ],

            colorArray:[
                'bg-green',
                'bg-yellow',
                'bg-orange',
                'bg-red'
            ],

            iconsAndColors: {
                'ack': {'icon': 'fa-check', 'color': 'bg-purple'},
                'ackremove': {'icon': 'glyphicon glyphicon-ban-circle', 'color': 'bg-purple'},
                'assocticket': {'icon': 'fa-ticket', 'color': 'bg-blue'},
                'declareticket': {'icon': 'fa-ticket', 'color': 'bg-blue'},
                'cancel': {'icon': 'glyphicon glyphicon-trash', 'color': 'bg-gray'},
                'uncancel': {'icon': 'glyphicon glyphicon-share', 'color': 'bg-gray'},
                'statusinc': {'icon': 'fa-chevron-up', 'color': 'bg-gray'},
                'statusdec': {'icon': 'fa-chevron-down', 'color': 'bg-gray'},
                'stateinc': {'icon': 'fa-flag', 'color': undefined},
                'statedec': {'icon': 'fa-flag', 'color': undefined},
                'changestate': {'icon': 'fa-flag', 'color': undefined},
                'snooze': {'icon': 'fa-clock-o', 'color': 'bg-fuchsia'}
            },



            init: function() {
                this._super();
            },

            stepsLoader: function() {

                if (!this.get('isHidden')) {

                
                    var component = this;

                    var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alarm');
                    var query = {'entity_id': this.get('alarm.entity_id')};

                    adapter.findQuery('alarm', 'get-current-alarm', query).then(function (result) {
                        // onfullfillment

                        var test = [
                            {
                                t: new Date().getTime(),
                                _t: 'ackremove',
                                val: 1
                            },
                            {
                                t: new Date().getTime(),
                                _t: 'statusdec',
                                val: 2
                            },
                            {
                                t: new Date().getTime(),
                                _t: 'snooze',
                                val: 3
                            }
                        ];
                        var steps = [];
                        if (result.data) {
                        // if (true) {
                            for (var i = result.data[0].value.steps.length - 1 ; i >= 0 ; i--) {
                            // for (var i = test.length - 1; i >= 0 ; i--) {
                                
                                var step = result.data[0].value.steps[i];
                                // var step = test[i];
                                

                                //build time related information
                                var date = new Date(step.t*1000);
                                step.date = moment(date).format('LL');
                                step.time = moment(date).format('h:mm:ss a');

                                //build color class
                                step.color = get(component,'iconsAndColors')[step._t].color;
                                //set icon
                                step.icon = get(component,'iconsAndColors')[step._t].icon;

                                //if no color, it's a state/value, so color
                                if(!step.color)
                                    step.color = get(component,'colorArray')[step.val];

                                if(step._t.indexOf('state') > -1)
                                    step.state = get(component,'stateArray')[step.val];

                                if(step._t.indexOf('status') > -1)
                                    step.status = get(component,'statusArray')[step.val];

                                if(step._t === 'snooze') {
                                    var until = new Date(step.val * 1000);
                                    step.until = moment(until).format('h:mm:ss a');
                                }

                                step.name = get(component,'statusToName')[step._t];

                                steps.push(step);
                            }
                        }

                            /* steps can be null if entity has no current alarm. */
                            set(component, 'steps', steps);
                    }, function (reason) {
                        // onrejection
                        console.error('ERROR in the adapter: ', reason);
                    });
                }
            }.observes('isHidden', 'alarm.id'),

        });

        application.register('component:component-customtimeline', component);
    }
});