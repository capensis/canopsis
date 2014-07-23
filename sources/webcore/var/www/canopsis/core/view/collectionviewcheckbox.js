/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

define([
    'ember',
    'app/application',
    'app/view/crecords'
], function(Ember, Application) {


Application.CollectionViewCheckBox = Ember.CollectionView.extend({

    /**
     * Create a deep Copy.(no change are made on account instance until save boutton is pressed)
     */
     //TODO : Make a base class for array from this one
    init: function() {
        var value = this.get("value");
        var valueRef = this.get("valueRef");

        console.log("valueRef", valueRef);
        //FIXME @Momo does not work on record creation
        if (valueRef === undefined) {
            valueRef = [false];
        }

        value = valueRef.slice(0);

        this.set("value",value);
        this._super();
    },

    registerFieldWithController: function() {
        var ArrayFields = this.get('controller.ArrayFields');
        if (ArrayFields) {
            ArrayFields.pushObject(this);
        }
    }.on('didInsertElement'),

    onUpdate: function() {
        var value = this.get("value");
        var valueRef = this.get("valueRef");

        while(valueRef.length > 0) {
            valueRef.pop();
        }

        for (var key in value) {
            valueRef[key] = value[key] ;
        }
        this.set("valueRef", valueRef);
    },

    valueRef: "",
    value: "",
    templateName: 'collection_template',
    itemViewClass: Ember.View.extend({
        tagName: 'li',
        template: Ember.Handlebars.compile("<label> {{view Canopsis.Application.customCheckBoxView contentBinding='title'}}{{title}}</label>")
    }),
    content: [
        { title: 'Read' },
        { title: 'Write' }
    ]
});
    return Application.CollectionViewCheckBox;

});