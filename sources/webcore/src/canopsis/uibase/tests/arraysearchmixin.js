module('uibase', {
  beforeEach: function() {
    visit('/userview/view.event');
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Uibase: arraysearch mixin', function() {
    expect(2);

    createNewView('list_search_test');

    activateEditMode();
    click('.btn-add-widget');

    waitForElement('.form .ember-text-field').then(function(){
        fillIn('.form .ember-text-field', 'list');
        click('.form .panel-default:first a');
        click('.form .list-group-item a');

        click('.form #mixins_tab a');
        fillIn('.tab-pane.active input', 'arraysearch');
        click('.form .panel-default .panel-heading a');
        click('.form .panel-default .list-group a');

        click('.form .btn-submit');

        waitMilliseconds(1000).then(function(){
            fillIn('.widgetslot .search input', 'root');
            click('.widgetslot .search .fa-search');
            waitMilliseconds(200).then(function(){
                equal(find('.widgetslot table tr').length, 2, 'Only two (one for headers, one for results) rows displayed');
                click('.widgetslot .search .glyphicon-remove-circle');
                waitMilliseconds(200).then(function(){
                    equal(find('.widgetslot table tr').length > 2, true, 'more than two rows displayed');
                });
            });
        });
    });
});
