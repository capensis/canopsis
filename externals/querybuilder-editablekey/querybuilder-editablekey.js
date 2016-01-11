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
                var select = realSelect.editableSelect();

                select.change(function() {
                    var selectVal = select.val();

                    if (realSelect.find("option[value='" + selectVal + "']").length === 0) {

                        var r = {
                            field: select.val(),
                            id: select.val(),
                            input: "text",
                            label: select.val(),
                            type: "string"
                        };

                        queryBuilder.addFilter(r);

                        console.error('append new option', select.val());
                        realSelect.append($('<option>', {
                            value: select.val(),
                            text: select.val()
                        }));
                    }

                    console.error(arguments);
                    rule.$el.find('.rule-filter-container select')
                        .val(select.val())
                        .trigger('change');
                });
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
