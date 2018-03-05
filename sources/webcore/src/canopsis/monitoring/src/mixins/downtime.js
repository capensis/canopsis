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
    name: 'DowntimeMixin',
    after: ['MixinFactory', 'FormsUtils', 'HashUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');
        var formsUtils = container.lookupFactory('utility:forms');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @mixin downtime
         *
         * @description
         * Implement downtime handling for widgets that manages collections
         *
         * Useful in lists for example, where it adds buttons to donwtime list elements
         */
        var mixin = Mixin('Downtime', {
            partials: {
                selectionToolbarButtons: [],
                actionToolbarButtons: [],
                itemactionbuttons: [],
                header: [],
                subHeader: [],
                footer: []
            },

            mixinsOptionsReady: function () {
                this._super();

                if (!get(this, 'mixinOptions.downtime.hideRemove')) {
                    get(this,'partials.selectionToolbarButtons').push('actionbutton-removeselection');
                    get(this,'partials.itemactionbuttons').push('actionbutton-remove');
                }
                if (!get(this, 'mixinOptions.downtime.hideEdit')) {
                    get(this,'partials.itemactionbuttons').push('actionbutton-edit');
                }
                if (!get(this, 'mixinOptions.downtime.hideCreate')) {
                    get(this,'partials.actionToolbarButtons').push('actionbutton-create');
                }

                set(this, 'itemsPerPagePropositionSelected', get(this, 'itemsPerPage'));
            },

            userCanReadRecord: function() {
                if(get(this, 'user') === 'root') {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');

                return get(this, 'rights.' + crecord_type + '_read.checksum');
            }.property('config.listed_crecord_type'),

            userCanCreateRecord: function() {
                if(get(this, 'user') === 'root') {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');

                return get(this, 'rights.' + crecord_type + '_create.checksum');
            }.property('config.listed_crecord_type'),

            userCanUpdateRecord: function() {
                if(get(this, 'user') === 'root') {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');

                return get(this, 'rights.' + crecord_type + '_update.checksum');
            }.property('config.listed_crecord_type'),

            userCanDeleteRecord: function() {
                if(get(this, 'user') === 'root') {
                    return true;
                }

                var crecord_type = get(this, 'listed_crecord_type');

                return get(this, 'rights.' + crecord_type + '_delete.checksum');
            }.property('config.listed_crecord_type'),

            actions: {

                edit: function (record) {
                    console.log('edit', record);

                    var extraoptions = get(this, 'mixinOptions.downtime.formoptions'),
                        formclass = get(this, 'mixinOptions.downtime.form');
                    var formoptions = {
                        title: 'Edit ' + get(record, 'crecord_type')
                    };

                    if(!isNone(extraoptions)) {
                        $.extend(formoptions, extraoptions);
                    }

                    if(isNone(formclass)) {
                        formclass = 'modelform';
                    }

                    console.log('open form:', formclass, formoptions);

                    var listController = this;
                    var recordWizard = formsUtils.showNew(formclass, record, formoptions);

                    recordWizard.submit.then(function(form) {
                        console.log('record going to be saved', record, form);

                        record = get(form, 'formContext');

                        record.save();

                        listController.trigger('refresh');
                    });
                }
            }
        });

        application.register('mixin:downtime', mixin);
    }
});
