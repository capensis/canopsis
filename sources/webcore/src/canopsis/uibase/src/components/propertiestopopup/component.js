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
    name: 'component-propertiestopopup',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            json2html = window.json2html;

        /**
         * @component propertiestopopup
         */
        var component = Ember.Component.extend({
            init: function() {
                this._super();

                var propertieslist;
                try{
                    propertieslist = JSON.parse(get(this, 'properties'));
                    if (!Ember.isArray(propertieslist)) {
                        throw 'Not an array';
                    }
                }catch(err) {
                    console.warn('Unable to parse properties list');
                    propertieslist = [];
                }
                set(this, 'propertieslist', propertieslist);
            },

            propertiesAsHtml: function(){
                //Generate a html rendering for choosen data in the properties list
                var propertieslist = get(this, 'propertieslist');
                var length = propertieslist.length;
                var source = get(this, 'source');

                if (isNone(source)) {
                    return '';
                } else {
                    var data,
                        i;
                    if (get(this, 'propertiesOnly')) {
                        data = [];
                    } else {
                        data = {};
                    }
                    for (i=0; i<length; i++) {
                        var value = get(source, propertieslist[i]);
                        //When data found
                        if (!isNone(value)) {
                            if (get(this, 'propertiesOnly')) {
                                data.push(value);
                            } else {
                                data[propertieslist[i]] = value;
                            }
                        }
                    }
                    var html = json2html(data);
                    if(get(this, 'htmlOnly')) {
                        return html;
                    } else {
                        return html.toString();
                    }
                }

            }.property(),

            icon: function () {
                // allow set custom icon in the span display
                var icon = get(this, 'customIcon');
                if (isNone(icon)) {
                    return 'glyphicon-eye-open';
                } else {
                    return icon;
                }
            }.property()
        });

        application.register('component:component-propertiestopopup', component);
    }
});
