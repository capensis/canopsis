module('querybuilder editor');

test('Querybuilder crashing attempt', function() {
    visit('/userview/view.event');

    expect(2);

    click('.left-side .frontend-config');
    waitForElement('input[name=title]').then(function(){
        click('.modal-body a[href=#editors]');
        fillIn('.modal-body .tab-pane.active input:first', 'filter');
        fillIn('.modal-body .tab-pane.active input:last', 'querybuilder');
        click('.modal-body .btn-success');
        click('.modal-footer .btn-submit');
        waitForElementRemoval('.modal-backdrop').then(function() {
            click('.btn-add-customfilter');

            waitForElement('input[name=title]').then(function(){
                fillIn('input[name=title]', 'test filter 1');
                find('.rule-filter-container select').val('enable').change();
                click('.rule-value-container input[value=true]:first');
                click('.builder .btn[data-add=rule]');

                waitMilliseconds(100).then(function(){
                    find('.rule-filter-container select:last').val('resource').change();
                    fillIn('.builder .rule-value-container:last input', 'Engine_perfdata');
                    click('.builder .btn[data-add=group]');
                    waitMilliseconds(100).then(function(){
                        find('.rule-filter-container select:last').val('crecord_type').change();
                        fillIn('.builder .rule-value-container:last input', 'event');
                        click('.query-builder a[aria-controls=output]');
                        equal(find('.query-builder div[role=output] pre').html().replace(/\s/g, ''), '{"$and":[{"enable":true},{"resource":"Engine_perfdata"},{"$and":[{"crecord_type":"event"}]}]}', 'generated filter seems correct');
                        click('.query-builder .btn-reset');
                        waitMilliseconds(100).then(function(){
                            equal(find('.query-builder div[role=output] pre').html().replace(/\s/g, ''), '{}', 'After hitting reset button, the filter is empty');
                            click('.modal-footer .btn-submit');
                        });
                    });
                });
            });
        });
    });
});
