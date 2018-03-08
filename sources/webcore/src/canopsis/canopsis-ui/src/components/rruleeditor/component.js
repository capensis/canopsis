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
    name: 'componentRruleEditor',
    initialize: function(container, application) {

        var __ = Ember.String.loc,
            get = Ember.get,
            set = Ember.set,
            isArray = Ember.isArray,
            RRule = window.RRule;

        /**
         * @method humanReadableRrule
         * @private
         * @param {String} correctly formated rrule
         * @returns human readable rrule or error message
         */
        var humanReadableRrule = function(){
            var text = '';
            if(get(this, 'rruleValue') == undefined){
				return __("Empty rrule");
			}
            try{
                var rruleObject = RRule.fromString(get(this, 'rruleValue'));
                text = rruleObject.toText();
            } catch(err) {
                text = __('Invalid RRule loaded');
            }
            return text;
        };

        var explodeArray = function(strToExplode){
            if(strToExplode === undefined)
                return undefined;

            var actual = strToExplode.split(',');

            //TODO move this in utils, somewhere ...
            //delete empty key
            var newArray = new Array();
            for (var i = 0; i < actual.length; i++) {
                if (actual[i]) {
                    newArray.push(actual[i]);
                }
            }

            return newArray;
        };


        /**
         * @component rruleeditor
         * @description A rrule editor component to write rrules
         */
        var component = Ember.Component.extend({
            /**
             * @property {String} ruleValue Rrule value, property binded in
             * editor
             */
            rruleValue: undefined,
            /**
             * @property {String} rruleText Property binded to tpl, hand typed
             * raw rrule
             */
            rruleText: undefined,
            /**
             * @property {String} rruleHuman Computed property converting rrule
             * Value to human sentence
             */
            rruleHuman: humanReadableRrule.property('rruleValue'),

            /**
             * @property {Boolean} isRruleValid Computed property that indicates if a rrule is valid or not
             **/
            isRruleValid: function(){
                var rruleHuman = '';
                try{
                    var rruleObject = RRule.fromString(get(this, 'rruleValue'));
                    rruleHuman = rruleObject.toText();
                } catch(err) {
                    rruleHuman = __('Invalid RRule loaded');
                }
                if (rruleHuman !== 'RRule error: Unable to fully convert this rrule to text'){
                    return true;
                } else {
                    return false;
                }
            }.property('rruleValue'),

            /**
             * @property {Object} tempRule Buffer between rrule updates
             */
            tempRule: {freq:3},

            /**
             * @method init
             * @description Rruleeditor init load rrule if given, set "input"
             * and load the rule in tempRule attribut
             */
            init: function () {
                this._super.apply(this, arguments);

                //loading rrule if we're in editing mode
                var rrule = get(this, 'rruleValue');
                if(rrule){
                    //we get the rrule parameters
                    var origOptions = RRule.fromString(rrule).origOptions;
                    var arrayToLoad = {};


                    //we loop on every key composing the loaded rrule
                    for(var key in origOptions){
                        //value of the rrule key
                        var value = origOptions[key];

                        switch(key){
                            case 'wkst':
                                arrayToLoad[key] = {'value': value.toString()};
                                set(this,'tempRule.'+key,value.toString());
                                break;
                            case 'freq':
                                arrayToLoad[key] = {'value': value};
                                set(this,'tempRule.'+key,value);
                                break;
                            case 'dtstart':
                            case 'until':
                                arrayToLoad[key] = value.getTime()/1000;
                                set(this,'tempRule.'+key,value);
                                break;
                            case 'byweekday':
                            case 'bymonth':
                                //get the element list list
                                var list = get(this,key + 'List');
                                //saving in tempRule
                                set(this,'tempRule.'+key, value);

                                if(! isArray(value))
                                    value = [value];

                                //we check box that need to be checked by
                                //iterating over the values of this rrule key
                                for (var i = 0; i < value.length; i++) {
                                    // value to load in form
                                    var dayOrMonth = value[i];

                                    // iterating over ****List attribut properties
                                    for(var listKey in list)
                                        // discard things added by ember (only want checkbox)
                                        if (list[listKey] && list[listKey].hasOwnProperty('isChecked'))
                                            //is it the value we searching ? (ex: Monday,ex: May)
                                            if (list[listKey].value === dayOrMonth)
                                                set(list[listKey], 'isChecked', true);
                                }
                                break;
                            case 'bysetpos':
                            case 'bymonthday':
                            case 'byyearday':
                            case 'byweekno':
                            case 'byhour':
                            case 'byminute':
                            case 'bysecond':
                                set(this,'tempRule.'+key,value);
                                if(isArray(value))
                                    arrayToLoad[key + 'Input']=  value.join();
                                else
                                    arrayToLoad[key + 'Input'] = String(value);
                                break;
                            default:
                                set(this,'tempRule.'+key,value);
                                arrayToLoad[key]= value;
                                break;
                        }
                    }
                    // load new properties, spare cpu of many triggers
                    this.setProperties(arrayToLoad);
                    set(this,'rruleText',rrule);
                }
            },

            /**
             * @method didInsertElement
             * @description contains Rrule-editor initialisation and data binding
             */
            didInsertElement: function() {
                // hugly hack, but it's the official way
                // https://guides.emberjs.com/v2.0.0/object-model/observers/
                // #toc_unconsumed-computed-properties-do-not-trigger-observers
                this.get('bysetpos');
                this.get('bymonthday');
                this.get('byyearday');
                this.get('byhour');
                this.get('byminute');
                this.get('bysecond');
                this.get('byweekno');
            },

            /**
             * @method  watchRruleText
             * @param  {Object} _obj Emberjs obj returned by Ember observer
             * @param  {String} name Name of var observed by Ember observer
             * @description Observe and update rruleValue according to RruleText typed value
             */
            watchRruleText: function(_obj,name){
                set(this, 'rruleValue', get(_obj,name));
            }.observes('rruleText'),

            /**
             * @method updateRrule
             * @param  {Object} _obj Emberjs obj returned by Ember observer
             * @param  {String} name Name of var observed by Ember observer
             * @description Observe and update rruleValue according to component form's
             */
            updateRrule: function (_obj,name) {
                var value = get(_obj,name);

                //case of freq.value and wkstart.value
                var keyName = name.replace('.value','');
                var tempRule = get(this,'tempRule');

                //if reset value or no value
                if(value === '' || value === '0'){
                    if(keyName in tempRule)
                        delete tempRule[keyName];
                }else{
                    //dispatch new value in temp dict
                    if(name ==='dtstart' || name ==='until')
                        //case for date
                        this.tempRule[keyName] = new Date(value*1000);
                    else
                        this.tempRule[keyName] = value;
                }

                //building rule
                var rule = new RRule(
                    Ember.copy(tempRule)
                );

                set(this, 'rruleValue', rule.toString());
                set(this, 'rruleText', rule.toString());
            }.observes(
                'wkst.value',
                'freq.value',
                'dtstart',
                'until',
                'count',
                'interval',
                'bysetpos',
                'bymonth',
                'bymonthday',
                'byyearday',
                'byweekday',
                'byweekno',
                'byhour',
                'byminute',
                'bysecond',
                'byeaster'
            ),

            /**
             * to icalendar spec
             * @property {Object} freq Default freq, according
             */
            freq: {value: RRule.DAILY},

            /**
             * @property {Object} frequencyList List of freq possible values
             */
            frequencyList: [
                {name:__('Secondly'),value: RRule.SECONDLY},
                {name:__('Minutely'),value: RRule.MINUTELY},
                {name:__('Hourly'),value: RRule.HOURLY},
                {name:__('Daily'),value: RRule.DAILY},
                {name:__('Weekly'),value: RRule.WEEKLY},
                {name:__('Monthly'),value: RRule.MONTHLY},
                {name:__('Yearly'),value: RRule.YEARLY}
            ],

            /**
             * @property {Integer} dtstart FrequencyList Starting date
             */
            dtstart: Date.now()/1000,

            /**
             * @property {Integer} until Stopping date
             */
            until: Date.now()/1000 + 24*60*60,

            /**
             * @property {Integer} count Number of occurence
             */
            count: undefined,
            /**
             * @property {Integer} interval Interval between occurence
             */
            interval: undefined,

            //************** Advance panel *************

            /**
             * @property {Object} wkstart Default wkstart, according
             * to icalendar spec
             */
            wkst: {value:'MO'},

            /**
             * @property {Object} wkstartList List of wkstart possible values
             */
            wkStartList: [
                {name:__('Monday'),value:'MO'},
                {name:__('Tuesday'),value:'TU'},
                {name:__('Wednesday'),value:'WE'},
                {name:__('Thursday'),value:'TH'},
                {name:__('Friday'),value:'FR'},
                {name:__('Saturday'),value:'SA'},
                {name:__('Sunday'),value:'SU'}
            ],

            /**
             * @property {Array} byweekday Day of the week where the rrule applies
             */
            byweekday: [],

            /**
             * @property {Object} byweekdayList List of correct values binded in
             * component hbs
             */
            byweekdayList: [
                {name:__('Monday'),value: RRule.MO, isChecked: false},
                {name:__('Tuesday'),value: RRule.TU, isChecked: false},
                {name:__('Wednesday'),value: RRule.WE, isChecked: false},
                {name:__('Thursday'),value: RRule.TH, isChecked: false},
                {name:__('Friday'),value: RRule.FR, isChecked: false},
                {name:__('Saturday'),value: RRule.SA, isChecked: false},
                {name:__('Sunday'),value: RRule.SU, isChecked: false}
            ],

            /**
             * @property {Function} byWeekDayChange Watch byweekdayList and change byweekday property to
             * the according value
             */
            byWeekDayChange: function(){
                var tempArray = this.get(
                        'byweekdayList'
                    ).filterBy(
                        'isChecked',
                         true
                    ).getEach('value');

                set(this,'byweekday',tempArray);

            }.observes('byweekdayList.@each.isChecked'),

            /**
             * @property {Array} bymonth Month where the rrule applies
             */
            bymonth:[],

            /**
             * @property {Array} bymonthList List of correct values binded in
             * component hbs
             */
            bymonthList: [
                {name:__('January'),value: 1, isChecked: false},
                {name:__('February'),value: 2, isChecked: false},
                {name:__('March'),value: 3, isChecked: false},
                {name:__('April'),value: 4, isChecked: false},
                {name:__('May'),value: 5, isChecked: false},
                {name:__('June'),value: 6, isChecked: false},
                {name:__('July'),value: 7, isChecked: false},
                {name:__('August'),value: 8, isChecked: false},
                {name:__('September'),value: 9, isChecked: false},
                {name:__('October'),value: 10, isChecked: false},
                {name:__('November'),value: 11, isChecked: false},
                {name:__('December'),value: 12, isChecked: false}
            ],

            /**
             * @property {Function} bymonthChange watch bymonthList and change bymonth
             * property to the according value
             */
            bymonthChange: function(){
                var tempArray = this.get(
                        'bymonthList'
                    ).filterBy(
                        'isChecked',
                         true
                    ).getEach('value');

                set(this,'bymonth',tempArray);

            }.observes('bymonthList.@each.isChecked'),

            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            bysetposInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            bymonthdayInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            byyeardayInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            byweeknoInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            byhourInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            byminuteInput: undefined,
            /**
             * @property { String } bysetposInput Property binded to tpl, raw
             * input values watched and parsed.
             */
            bysecondInput: undefined,


            /**
             * @property {Function} bysetpos Computed property from bysetposInput (string to array)
             */
            bysetpos: function(){
                return explodeArray(this.get('bysetposInput'));
            }.property('bysetposInput'),

            /**
             * @property {Function} bymonthday Computed property from bymonthdayInput (string to array)
             */
            bymonthday: function(){
                return explodeArray(this.get('bymonthdayInput'));
            }.property('bymonthdayInput'),

            /**
             * @property {Function} byyearday Computed property from byyeardayInput (string to array)
             */
            byyearday: function(){
                return explodeArray(this.get('byyeardayInput'));
            }.property('byyeardayInput'),

            /**
             * @property {Function} byhour Computed property from byhourInput (string to array)
             */
            byhour: function(){
                return explodeArray(this.get('byhourInput'));
            }.property('byhourInput'),

            /**
             * @property {Function} byweekno Computed property from byweeknoInput (string to array)
             */
            byweekno: function(){
                return explodeArray(this.get('byweeknoInput'));
            }.property('byweeknoInput'),

            /**
             * @property {Function} byminute Computed property from byminuteInput (string to array)
             */
            byminute: function(){
                return explodeArray(this.get('byminuteInput'));
            }.property('byminuteInput'),

            /**
             * @property {Function} bysecond Computed property from bysecondInput (string to array)
             */
            bysecond: function(){
                return explodeArray(this.get('bysecondInput'));
            }.property('bysecondInput')

        });

        application.register('component:component-rruleeditor', component);
    }
});
