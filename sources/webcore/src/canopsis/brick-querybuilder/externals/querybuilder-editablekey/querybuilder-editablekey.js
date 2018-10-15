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
define([
    'query-builder'
], function (QueryBuilder) {
    var QueryBuilder = $.fn.queryBuilder;

    /*!
     * jQuery QueryBuilder Bootstrap Selectpicker
     * Applies Bootstrap Select on filters and operators combo-boxes.
     * Copyright 2014-2015 Damien "Mistic" Sorel (http://www.strangeplanet.fr)
     */
    QueryBuilder.define('key-editable-select', function(options) {
        var queryBuilder = this;

        if (!$.fn.editableSelect || !$.fn.editableSelect.constructor) {
            console.error('editableSelect jquery plugin is required to use "key-editable-select" plugin.');
        }

        // // init selectpicker
        this.on('afterCreateRuleFilters', function(e, rule) {
            // rule.$el.find('.rule-filter-container select').editableSelect();
        });

        this.on('afterCreateRuleOperators', function(e, rule) {
            if(rule.$el.find('.rule-filter-container input').length === 0) {
                var realSelect = rule.$el.find('.rule-filter-container select');

                var changeFunction = function(selectValParam) {
                    var selectVal = selectValParam;

                    if(typeof selectVal !== 'string') {
                        selectVal = select.val();
                    }

                    if (realSelect.find("option[value='" + selectVal + "']").length === 0) {

                        var r = {
                            field: selectVal,
                            id: selectVal,
                            input: "text",
                            label: selectVal,
                            type: "string"
                        };

                        queryBuilder.addFilter(r);

                        realSelect.append($('<option>', {
                            value: selectVal,
                            text: selectVal
                        }));
                    }

                    rule.$el.find('.rule-filter-container select')
                        .val(selectVal)
                        .trigger('change');
                };

                var select = realSelect.editableSelect({
                    onSelect: function () {
                        select.val($(this).val()).change();
                        changeFunction($(this).val());
                    }
                });
                select.css('width', '150px');

                select.change(changeFunction);
            }
        });

        // update selectpicker on change
        this.on('afterUpdateRuleFilter', function(e, rule) {
            // rule.$el.find('.rule-filter-container select').selectpicker('render');
        });

        this.on('afterUpdateRuleOperator', function(e, rule) {
            // rule.$el.find('.rule-filter-container select').selectpicker('render');
        });
    }, {
        container: 'body',
        style: 'btn-inverse btn-xs',
        width: 'auto',
        showIcon: false
    });
});
