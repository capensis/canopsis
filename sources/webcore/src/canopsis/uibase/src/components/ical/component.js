/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'component-ical',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc,
            RRule = window.RRule;

        /**
         * @description Manage Ical format
         * @component ical
         * @example
         * <div class="well">
         *  <center>
         *       <h1>{{title}}</h1>
         *   </center>
         *   <form class="form-horizontal" role="form">
         *
         *       <div class="form-group">
         *           <label class="col-sm-2 control-label">
         *               {{#component-tooltip content=helpFrequency}}
         *                   {{tr 'Frequency'}}
         *               {{/component-tooltip}}
         *           </label>
         *           <div class="col-sm-10">
         *
         *               {{view Ember.Select
         *                   class="form-control"
         *                   content=frequencyForm
         *                   optionValuePath="content.value"
         *                   optionLabelPath="content.label"
         *                   value=frequencySelection
         *               }}
         *
         *           </div>
         *       </div>
         *
         *       <div class="form-group">
         *           <label class="col-sm-2 control-label">
         *               {{#component-tooltip content=helpRepetitionCount}}
         *                   {{tr 'Repetition count'}}
         *               {{/component-tooltip}}
         *           </label>
         *           <div class="col-sm-10">
         *               {{input class="form-control" value=repetitionCount}}
         *           </div>
         *       </div>
         *
         *       <div class="form-group">
         *           <label class="col-sm-2 control-label">
         *               {{#component-tooltip content=helpRepetitionInterval}}
         *                   {{tr 'Repetition interval'}}
         *               {{/component-tooltip}}
         *           </label>
         *           <div class="col-sm-10">
         *               {{input class="form-control" value=repetitionInterval}}
         *           </div>
         *       </div>
         *
         *       <div class="form-group">
         *           <label class="col-sm-2 control-label">
         *               {{#component-tooltip content=helpStartDate}}
         *                   {{tr 'Start date'}}
         *               {{/component-tooltip}}
         *           </label>
         *           <div class="col-sm-10">
         *               {{component-datetimepicker content=startDate}}
         *           </div>
         *       </div>
         *
         *       <div class="form-group">
         *           <label class="col-sm-2 control-label">
         *               {{#component-tooltip content=helpStopDate}}
         *                   {{tr 'Stop date'}}
         *               {{/component-tooltip}}
         *           </label>
         *           <div class="col-sm-10">
         *               {{component-datetimepicker content=stopDate}}
         *           </div>
         *       </div>
         *
         *       <center>
         *           <button class="btn btn-default" {{action "addRule"}}>{{tr 'Add rule'}}</button>
         *       </center>
         *
         *  </form>
         *
         *   <hr />
         *
         *
         *   <ul class="list-group">
         *       {{#each rule in rules}}
         *           <li  class="list-group-item list-group-item-success">
         *               <button {{action "removeRule" rule}} class="btn btn-default">
         *                   <span class="glyphicon glyphicon-minus"></span>
         *                    {{tr 'Remove'}}
         *               </button>
         *               {{rule.value}}
         *           </li>
         *       {{/each}}
         *   </ul>
         * </div>
         *
         */
        var component = Ember.Component.extend({

            //translated by tooltip
            helpFrequency: 'How often the rule is applied',
            helpRepetitionInterval: 'Time space between two repetitions',
            helpRepetitionCount: 'How many times action is repeated',
            helpStartDate: 'Since when to apply rule',
            helpStopDate: 'Until when to apply rule',


            frequencyForm: [
                {value: RRule.YEARLY, label: __('Yearly')},
                {value: RRule.MONTHLY, label: __('Monthly')},
                {value: RRule.WEEKLY, label: __('Weekly')},
                {value: RRule.DAILY, label: __('Daily')},
                {value: RRule.HOURLY, label: __('Hourly')},
                {value: RRule.MINUTELY, label: __('Minutely')},
                {value: RRule.SECONDLY, label: __('Secondly')}
            ],

            /*
            //not used yet
            months: [
                {id: 1, label: __('jan')},
                {id: 2, label: __('feb')},
                {id: 3, label: __('mar')},
                {id: 4, label: __('apr')},
                {id: 5, label: __('may')},
                {id: 6, label: __('jun')},
                {id: 7, label: __('jul')},
                {id: 8, label: __('aug')},
                {id: 9, label: __('sep')},
                {id: 10, label: __('oct')},
                {id: 11, label: __('nov')},
                {id: 12, label: __('dec')},
            ],

            proxiedMonth: Ember.computed.map('model', function(model){
                return Ember.ObjectProxy.create({
                  content: model,
                  checked: false
                });
            }),
            */

            /**
             * @description Initialize the component
             * @method init
             */
            init: function() {
                this._super();

                var monthmodel = DS.Model.extend({
                    label: DS.attr()
                });

                monthmodel.FIXTURES = get(this,'months');

                set(this, 'rules', []);
                set(this, 'frequencySelection', RRule.YEARLY);

                console.log('rrule information', get(this, 'frequencyForm'));

                //initialization from existing rrules content
                var content = get(this, 'content');
                if (!isNone(content) && Ember.isArray(content)) {
                    var contentlen = content.length;
                    for (var i=0; i< contentlen; i++) {
                        try {
                            var rule = RRule.fromString(content[i]);
                            var ruleObject = {value: rule.toText(), instance: rule};
                            get(this, 'rules').pushObject(ruleObject);
                        } catch (err) {
                            //error appends half the time form lib i dont know why
                            console.warn('Unable to parse rrule from given information', content[i], err);
                        }
                    }
                } else {
                    console.log('no content available for initialization');
                }
                this.updateContent();

            },

            actions: {
                /**
                 * @description remove a rule
                 * @method actions_removeRule
                 */
                removeRule: function (rule) {
                    get(this, 'rules').removeObject(rule);
                    this.updateContent();
                },

                /**
                 * @description Add a rule
                 * @method actions_addRule
                 */
                addRule: function () {

                    //Generating rrule options
                    var rrule = {};

                    var startDate = get(this, 'startDate');
                    if (!isNone(startDate) && !isNaN(startDate)) {
                        rrule.dtstart = new Date(startDate * 1000);
                    }

                    var stopDate = get(this, 'stopDate');
                    if (!isNone(stopDate) && !isNaN(stopDate)) {
                        rrule.until = new Date(stopDate * 1000);
                    }

                    var repetitionCount = get(this, 'repetitionCount');
                    if (!isNone(repetitionCount) && !isNaN(repetitionCount)) {
                        rrule.count = parseInt(repetitionCount);
                    }

                    var repetitionInterval = get(this, 'repetitionInterval');
                    if (!isNone(repetitionInterval) && !isNaN(repetitionInterval)) {
                        rrule.interval = parseInt(repetitionInterval);
                    } else {
                        //is set at least by default to 30
                        console.log('rrule inteval default value set to 30 as no value set');
                        rrule.interval = 30;
                    }

                    var frequencySelection = get(this, 'frequencySelection');
                    console.log('generated rrule options',  rrule, 'frequencySelection ', frequencySelection );
                    var rruleInstance =new RRule(frequencySelection, rrule);

                    var ruleObject = {value: rruleInstance.toText(), instance: rruleInstance};

                    get(this, 'rules').pushObject(ruleObject);

                    this.updateContent();
                }
            },

            /**
             * @description Update the component with new rules
             * @method actions_updateContent
             */
            updateContent: function () {
                var content = [];
                var rules = get(this,'rules');
                var ruleslen = rules.length;
                for (var i=0; i<ruleslen; i++) {
                    content.push(rules[i].instance.toString());
                }

                console.log('generated rules', content);
                set(this, 'content', content);
            },

            /**
             * @description Simply display the component
             * @method actions_didInsertElement
             */
            didInsertElement:function () {
                console.log('recurrence input loaded', this.$());
            }
        });
        application.register('component:component-rrulebak', component);
    }
});
