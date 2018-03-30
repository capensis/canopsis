module('uibase', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Simple list creation', function() {
    expect(1);

    visit('/userview/view.event');

    createNewView('list_test');

    activateEditMode();
    click('.btn-add-widget');

    waitForElement('.form .ember-text-field').then(function(){
        fillIn('.form .ember-text-field', 'list');
        click('.form .panel-default:first a');
        click('.form .list-group-item a');
        click('.form .btn-submit');

        waitMilliseconds(1000).then(function(){
            equal(find('.table-responsive').length, 1, 'there is one table in the view');
        });
    });
});
