/**
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
    name: 'DatesUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {

        var Utility = container.lookupFactory('class:utility');

        var __ = Ember.String.loc,
            isNone = Ember.isNone;

        var dates = Utility.create({

            name: 'dates',

            getNow: function() {
                return parseInt(new Date().getTime() / 1000);
            },

            getStringNow: function(format, shortDate) {
                return dates.timestamp2String(dates.getNow(), format, shortDate);
            },

            durationFromNow: function (timestamp) {
                var delta = dates.getNow() - timestamp;
                return dates.second2Duration(delta);
            },

            second2Duration: function (totalSec) {

                var days = parseInt( totalSec / (3600 * 24) );
                var hours = parseInt( totalSec / 3600 ) % 24;
                var minutes = parseInt( totalSec / 60 ) % 60;
                var seconds = parseInt(totalSec % 60);

                var displayHours = '';
                if (hours) {
                    displayHours = (hours < 10 ? "0" + hours : hours) + 'h ';
                }

                var displayMinutes = '';
                if (minutes) {
                    displayMinutes = (minutes < 10 ? "0" + minutes : minutes) + 'm ';
                }


                var result = displayHours +
                    displayMinutes +
                    (seconds  < 10 ? "0" + seconds : seconds) + 's';

                if (!isNaN(days) && days !== 0) {
                    result = days + 'd ' + result;
                }

                return result;
            },

            timestamp2String: function (value, format, shortDate) {
                function addZero(i) {
                    return (i < 10 ? '0'+ i : i.toString());
                }

                var a = new Date(value*1000);
                var months = [
                    __("January"),
                    __("February"),
                    __("March"),
                    __("April"),
                    __("May"),
                    __("June"),
                    __("July"),
                    __("August"),
                    __("September"),
                    __("October"),
                    __("November"),
                    __("December")
                ];
                if (shortDate === true) {
                    months = [
                        __("Jan"),
                        __("Feb"),
                        __("Mar"),
                        __("Apr"),
                        __("May"),
                        __("June"),
                        __("July"),
                        __("Aug"),
                        __("Sep"),
                        __("Oct"),
                        __("Nov"),
                        __("Dec")
                    ];

                }
                var year = a.getFullYear();
                var month = months[a.getMonth()];
                var date = addZero(a.getDate());
                var hour = addZero(a.getHours());
                var min = addZero(a.getMinutes());
                var sec = addZero(a.getSeconds());
                var time = "";

                switch(format) {
                    case 'f':
                        time = [date, month, year].join(' ') + ' ' + [hour, min, sec].join(':') ;
                        break;
                    case 'r':
                        time = [date, addZero(a.getMonth()+1), year].join('/') + ' <br>' + [hour, min, sec].join(':') ;
                        break;

                    case 'timeOnly':
                        time = [hour, min, sec].join(':') ;
                        break;

                    default:
                        time = [date, month, year].join(' ') + ' ' + [hour, min, sec].join(':') ;
                        break;
                }
                return time;
            },

            //bottom is not used yet
            locale: 'fr',
            months: ['January','February','March','April','May','June','July','August','September','October','November','December'],
            days: ['Monday', 'Thuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'],

            setLang: function (lang) {
                dates.locale = lang;
                if (lang === 'fr') {
                    dates.month = ['Janvier', 'Février', 'Mars', 'Avril', 'Mai', 'Juin', 'Juillet', 'Août', 'Octobre', 'Novembre', 'Décembre'];
                    dates.days = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'];
                }
            },

            /**
                Computes for a timestamp the timestamp at midnight of it's day.
                the computation day from timestamp depends on if a timestamp is given.
                when given start of the day is computed from timestamp, otherwise it is done from now timestamp
            **/
            startOfTheDay: function (aTimestamp) {
                if (isNone(aTimestamp)) {
                    aTimestamp = dates.getNow();
                    console.log('got date from now as no param given', aTimestamp);
                }
                aTimestamp *= 1000;
                var startDateOfTheDay = new Date(aTimestamp);
                startDateOfTheDay.setHours(0,0,0,0);
                return parseInt(startDateOfTheDay.getTime()/1000);
            },

            /**
                Boolean value determining wether the given date included in today
                Value depends on the client clock
            **/
            isToday: function (timestamp) {
                var startOfTheDay1 = dates.startOfTheDay(dates.getNow());
                var startOfTheDay2 = dates.startOfTheDay(timestamp);
                //console.error('coucou getNow ', dates.getNow())
                //console.error('coucou startOf', startOfTheDay2)
                // return dates.getNow() < startOfTheDay2 + 3600 && startOfTheDay2 < dates.getNow() + 3600;
                return startOfTheDay1 === startOfTheDay2;
            },

            dateFormat:'YYYY/MM/DD',

            diffDate: function(d1,d2,u) {
                div = 1;

                switch(u) {
                    case 's':
                        div=1000;
                        break;
                    case 'm':
                        div=1000*60;
                        break;
                    case 'h':
                        div=1000*60*60;
                        break;
                    case 'd':
                        div=1000*60*60*24;
                        break;
                    default:
                        break;
                }

                var Diff = d2 - d1;
                return Math.ceil((Diff/div));
            }
        });

        application.register('utility:dates', dates);
    }
});
