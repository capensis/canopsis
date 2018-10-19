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
    name:'CustomfilterlistMixin',
    after: ['FormsUtils', 'MixinFactory', 'DataUtils', 'NotificationUtils'],
    initialize: function(container, application) {

        var Mixin = container.lookupFactory('factory:mixin'),
            formsUtils = container.lookupFactory('utility:forms'),
            dataUtils = container.lookupFactory('utility:data'),
            notificationUtils = container.lookupFactory('utility:notification');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @mixin customfilterlist
         * @description
         * Implements Custom filter management for list
         * A filter is a combination of a cfilter and a title.
         * Custom cfilter allow perform selelection on a list with custom filter information.
         *
         * ![Mixin preview](../screenshots/mixin-customfilterlist.png)
         */
        var mixin = Mixin('customfilterlist', {
            partials: {
                subHeader: ['customfilters']
            },

            /**
            * Load custom filter if any form the widget model
            **/

            init: function() {
                this._super();

                if(!get(this, 'model.user_filters')) {
                    set(this, 'model.user_filters', Ember.A());
                }

                set(this, 'multiFiltersTitle', []);
                set(this, 'customFilterOperators', [
                    {
                        value: '$and',
                        display: '&',
                        active: true
                    },
                    {
                        value: '$or',
                        display: '||',
                        active: false
                    }
                ]);

            },

            /**
             * @property canAddCustomFiltersInUserPreferences
             * @type boolean
             * @description whether it is possible or not to add custom filters stored in userpreferences
             */
            canAddCustomFiltersInUserPreferences: true,

            /**
            * Builds the list of active filters
            *
            * @param filterlist : a list of filter that contains
            * a string mongo filter, a title and an active boolean
            **/

            isSelectedFilter: function (filterList) {
                if(!filterList || !filterList.length) {
                    return false;
                }

                if (!get(this, 'mixinOptions.customfilterlist.can_mix_filters')) {
                    var filterLen = filterList.length;
                    var currentTitle = get(this, 'model.selected_filter.title');
                    for (var i = 0; i < filterLen; i++) {

                        var compareTitle = get(filterList[i], 'title');

                        console.log('compare filters',currentTitle, compareTitle);

                        if (currentTitle === compareTitle) {
                            set(filterList[i], 'isActive', true);
                        } else {
                            set(filterList[i], 'isActive', false);
                        }
                    }
                }

                return filterList;
            },

            /**
            * Get filter list from the mixin parameters
            **/

            filters_list: function () {
                return this.isSelectedFilter(get(this, 'mixinOptions.customfilterlist.filters'));
            }.property('mixinOptions.customfilterlist.filters', 'model.selected_filter'),

            /**
            * Get filter list from the user parameters
            **/

            user_filters: function () {
                return this.isSelectedFilter(get(this, 'model.user_filters'));
            }.property('model.user_filters', 'model.selected_filter'),

            /**
            * Compute the complete filter from user selected filters.
            * It depends on what user selected in the UI
            **/

            computeFilterFragmentsList: function() {
                var list = this._super(),
                    mixinOptions = get(this, 'model.mixins').findBy('name', 'customfilterlist'),
                    userFilter;

                var selectedFilterFilter = get(this, 'model.selected_filter.filter');

                if(!isNone(selectedFilterFilter)) {
                    userFilter = selectedFilterFilter;
                } else if(get(this, 'model.selected_filter') && !selectedFilterFilter) {

                    userFilter = {};
                } else if(mixinOptions && mixinOptions.default_filter) {
                    userFilter = mixinOptions.default_filter;
                    userFilter = JSON.parse(userFilter);
                } else {
                    userFilter = {};
                }

                list.pushObject(userFilter);

                return list;
            },

            /**
            * Generate the complete filter from multi selection to be inserted in the computeFilterFragmentList
            * Filters elements are the whole user filter selection from user filter and widget filter.
            **/

            generateMultiFilter: function () {

                var widgetFilters = get(this, 'mixinOptions.customfilterlist.filters');
                var userFilters = get(this, 'model.user_filters');

                //get all filter available
                var allFilters = [];
                allFilters.addObjects(widgetFilters);
                allFilters.addObjects(userFilters);

                var length = allFilters.length,
                    titles = get(this, 'multiFiltersTitle'),
                    filter = {},
                    i;

                //get selector operator value
                var operators = get(this, 'customFilterOperators');
                var operatorLength = operators.length,
                    operatorValue;

                for(i=0; i<operatorLength; i++) {
                    if (get(operators[i], 'active')) {
                        operatorValue = get(operators[i], 'value');
                    }
                }

                filter[operatorValue] = [];

                //generate user selected filter list
                for (i=0; i<length; i++) {
                    var title = get(allFilters[i], 'title');
                    if ($.inArray(title, titles) !== -1) {
                        var subFilter = JSON.parse(get(allFilters[i], 'filter'));
                        filter[operatorValue].pushObject(subFilter);
                    }
                }

                //Test if there is sub filters.
                if (filter[operatorValue].length === 0) {
                    filter = {};
                }

                var multiFilter = JSON.stringify(filter);

                set(this, 'model.selected_filter', {filter: multiFilter});

            },

            actions: {

                /**
                * Set an operator active by desactivating all other operator and set active the one passed as parameter.
                *
                * @param operator: an operator dict that have to be set as active.
                **/

                setActiveOperator: function (operator) {

                    var operators = get(this, 'customFilterOperators');
                    var length = operators.length;

                    for(var i=0; i<length; i++) {
                        set(operators[i], 'active', false);
                    }

                    set(operator, 'active', true);

                    //reload widget with new filter when operator is changed
                    this.generateMultiFilter();

                    this.refreshContent();

                },


                /**
                * The entry point of the customfilterlist system.
                * By selecting a filter in the list, the user trigger this function
                * and make available a filter fragment into the widget.
                *
                * @param filter: a dict that contains a title and a string mongo filter and an active status.
                **/
                setFilter: function (filter) {

                    var canMixFilter = get(this, 'mixinOptions.customfilterlist.can_mix_filters');

                    if (canMixFilter) {
                        console.log('multi filter mode');

                        var filters = get(this, 'multiFiltersTitle'),
                            title = get(filter, 'title'),
                            active = !get(filter, 'isActive');
                        //set active filter
                        set(filter, 'isActive', active);

                        //when filter is active and not
                        if (active) {
                            filters.pushObject(title);
                        } else {
                            filters.removeObject(title);
                        }

                        this.generateMultiFilter();

                    } else {
                        console.log('single filter mode');

                        var query = get(filter, 'filter');
                        //current user filter set for list management
                        set(this, 'model.selected_filter', filter);
                        //user filter for list be able to reload properly
                        set(this, 'model.userFilter', query);
                        var attributesList = get(this, "userPreferencesModel.attributes.list")
                        // I do not know where and how to add the 'selected_filter'
                        // in order to no be filtered by the mixin UserConfigurationMixin
                        // So I add it in the model here.
                        attributesList.push({
                            isAttribute: true,
                            name: "selected_filter",
                            options: {
                                isUserPreference: true,
                                type: "boolean",
                            },
                            type: "boolean"
                        })
                        //push userparams to db
                        this.saveUserConfiguration();

                        //changing pagination information
                        if (!isNone(get(this, 'currentPage'))) {
                            set(this, 'currentPage', 1);
                        }
                    }

                    this.refreshContent();
                },

                /**
                * The action triggered when the user add a filter.
                * This leads to a filter save into the user preference storage
                **/

                addUserFilter: function () {
                    var widgetController = this;
                    var record = dataUtils.getStore().createRecord('customfilter', {
                        crecord_type: 'customfilter'
                    });

                    var recordWizard = formsUtils.showNew('modelform', record, {
                        title: __('Create a custom filter for current list')
                    });

                    recordWizard.submit.then(function(form) {
                        record = form.get('formContext');

                        get(widgetController, 'model.user_filters').pushObject(record);

                        widgetController.notifyPropertyChange('model.user_filters');


                        console.log('Custom filter created', record, form);
                        notificationUtils.info(__('Custom filter created'));

                        // get(widgetController, 'viewController').get('content').save();

                        widgetController.saveUserConfiguration();
                    });
                },

                /**
                * Action allowing filter edition from one existing user filter.
                *
                * @param filter: a dict with filter information to update into the user preferences
                **/

                editFilter: function (filter) {

                    var widgetController = this;

                    //rebuild a crecord as data may be simple js object saved to userpreferences
                    var record = dataUtils.getStore().createRecord('customfilter', {
                        crecord_type: 'customfilter',
                        filter: get(filter, 'filter'),
                        title: get(filter, 'title')
                    });

                    var recordWizard = formsUtils.showNew('modelform', record, {
                        title: __('Edit filter for current list')
                    });

                    recordWizard.submit.then(function(form) {
                        widgetController.get('model.user_filters').removeObject(filter);
                        record = form.get('formContext');
                        widgetController.get('model.user_filters').pushObject(record);
                        console.log('Custom filter created', record, form);
                        notificationUtils.info(__('Custom filter created'));

                        widgetController.saveUserConfiguration();
                    });
                },

                /**
                * The action to remove a filter from the user filter list.
                *
                * @param filter: dict that contains a filter information.
                **/

                removeFilter: function (filter) {

                    var title = get(filter, 'title');

                    var recordWizard = formsUtils.showNew('confirmform', {}, {
                        title: __('Are you sure to delete filter') + ' ' + title + '?'
                    });

                    var widgetController = this;

                    recordWizard.submit.then(function(form) {
                        void(form);

                        widgetController.get('model.user_filters').removeObject(filter);
                        notificationUtils.info(__('Custom filter removed'));

                        // get(widgetController, 'viewController').get('content').save();

                        widgetController.saveUserConfiguration();
                    });
                }
            }
        });

        application.register('mixin:customfilterlist', mixin);
    }
});
