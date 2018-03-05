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

(function() {

    var set = Ember.set,
        isNone = Ember.isNone,
        isArray = Ember.isArray;


    var helper = function(json, settableObject) {
        if (typeof json === 'string') {

            try {
                json = JSON.parse(json);
            } catch (exception) {
                console.error('unable to deserialize json');
                json = {'invalid': 'json string'};
            }
        }

        function parseJson(object) {
            var html = '';
            if (isArray(object)) {
                html += '<ul class="jsonUl">';
                var len = object.length;
                for (var i=0; i<len; i++) {
                    html += '<li class="jsonLi">' + parseJson(object[i]) + '</li>';
                }
                html += '</ul>';
            } else if (typeof object === 'object') {
                for (var key in object) {
                    html += '<ul class="jsonUl"><li class="jsonLi">';
                    html += '<span class="label label-primary">'+ key +'</span>';
                    if (typeof object[key] === 'object' || isArray(object[key])) {
                        html += parseJson(object[key]);
                    } else {
                        html +=  '&nbsp;<span class="glyphicon glyphicon-arrow-right" style="display:inline"></span><span class="label label-warning">'+ object[key] +'</span>' ;
                    }
                    html += '</li></ul>';
                }
            } else {
                html +=  '&nbsp;<span class="glyphicon glyphicon-arrow-right" style="display:inline"></span><span class="label label-warning">'+ object +'</span>' ;
            }
            return html;
        }

        var html = parseJson(json);

        console.info('json2html', {'data': json, 'settable object:': settableObject, 'html': html, 'arguments': arguments});

        //argument with json param and options only will have a length of 2
        if(!isNone(settableObject) && arguments.length > 2) {
            //Set html to object
            console.log('processing set html to the settable item');
            set(settableObject, 'json2html', html);
            //Do not print html
            html = '';
        }

        return new Ember.Handlebars.SafeString(html);
    };

    //declaring helper this way allow it to be used as simple function somewhere else.
    Ember.Handlebars.helper('json2html', helper);
    window.json2html = helper;
})();
