/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/


define([], function() {

    var get = Ember.get,
        set = Ember.set,
        isNone = Ember.isNone,
        __ = Ember.String.loc;


    var component = Ember.Component.extend({
        field: Ember.Object.create({
            'second' : __('Second'),
            'minute' :__('Minute'),
            'hour' : __('Hour'),
            'day' : __('Day'),
            'week' : __('Week'),
            'month' : __('Month'),
            'year' : __('Year')
        }),

        convertDuration: Ember.Object.create({
            'second' : 1,
            'minute' : 60,
            'hour' : 3600,
            'day' : 3600 * 24,
            'week': 3600 * 24 * 7,
            'month': 3600 * 24 * 30,
            'year' : 3600 * 24 * 365
        }),

        durationType: [
            'second',
            'minute',
            'hour',
            'day',
            'week',
            'month',
            'year'
        ],

        template: Ember.HTMLBars.compile([
        '<div class="form-inline">',
            '<div class="form-group">',
                '{{input value=shownDuration class="form-control" type="number"}}',
            '</div>',
            '<div class="form-group">',
                '{{view Ember.Select content=durationType value=selectedDurationType class="form-control"}}',
            '</div>',
        '</div>'].join('')),

        init: function () {
            this._super.apply(this, arguments);



            if (get(this, 'useFullPeriodInformation')) {
                //initialization fromobject that contains full information
                set(this, 'selectedDurationType', get(this, 'content.durationType'));
                set(this, 'shownDuration', get(this, 'content.value'));
            } else {
                var durationType = 'second';

                var unformattedDuration = parseInt(get(this, 'content'), 10);

                if (isNaN(unformattedDuration)) {
                    unformattedDuration = 0;
                }
                var convert = get(this, 'convertDuration');
                var durationUnits = get(this, 'durationType');
                var durationUnitsLen = durationUnits.length;

                for (var i = 0; i < durationUnitsLen; i++) {

                    var durationUnit = durationUnits[i];
                    console.log('testing duration unit', durationUnit);
                    var unitValue = convert[durationUnit];
                    if (unitValue > unformattedDuration) {
                        break;
                    } else {
                        durationType = durationUnit;
                    }
                }

                console.log('selected unit value is', durationType);
                set(this, 'selectedDurationType', durationType);

                var conversionOperand = get(this, 'convertDuration').get(durationType);
                var res = parseInt(unformattedDuration / conversionOperand);

                console.log(
                    'conversion operand is', conversionOperand,
                    'unformattedDuration', unformattedDuration,
                    'conversionOperand', conversionOperand,
                    'res', res
                );

                if (isNone(res) || isNaN(res)) {
                    res = 0;
                }
                set(this, 'shownDuration', res);
                console.debug('shown duration initialized to ', this.get('shownDuration') ,res);
            }

            //initilizes content with computed value
            this.shownDurationChanged();
        },

        shownDurationChanged: function () {
            console.log('shownDurationChanged');
            var durationType = get(this, 'selectedDurationType');
            var conversionOperand = get(this, 'convertDuration.' + durationType);
            if (isNone(conversionOperand)) {
                //dynamic initialization when undefined
                conversionOperand = 0;
            }
            var value = get(this, 'shownDuration');

            var computedValue = parseInt(value * conversionOperand);
            if (isNaN(computedValue)) {
                computedValue = 0;
            }
            console.log('computed value in content is', computedValue);

            this.setContent(computedValue);
        }.observes('shownDuration'),

        selectedDurationTypeChanged: function () {
            var durationType = get(this, 'selectedDurationType');
            var conversionOperand = get(this, 'convertDuration').get(durationType);
            console.log('conversionOperand', durationType, conversionOperand, get(this, 'shownDuration'));

            var newValue = get(this, 'shownDuration') * conversionOperand;
            console.log('selectedDurationTypeChanged', durationType, conversionOperand, newValue);

            this.setContent(newValue);
        }.observes('selectedDurationType'),

        selectedDurationType: 'second',

        setContent: function (seconds) {
            var durationInformation = seconds;
            if (get(this, 'useFullPeriodInformation')) {
                durationInformation = {
                    value: parseInt(get(this, 'shownDuration', seconds)),
                    durationType: get(this, 'selectedDurationType'),
                    seconds: seconds
                };
            }
            set(this, 'content', durationInformation);
        },

        selectedDurationLabel: function () {
            console.log('selectedDurationType');
            var durationType = get(this, 'selectedDurationType');
            return get(this, 'field.' + durationType);
        }.property('selectedDurationType')
    });


    Ember.Application.initializer({
        name:"component-durationcombo",
        initialize: function(container, application) {
            application.register('component:component-durationcombo', component);
        }
    });

    return component;
});
