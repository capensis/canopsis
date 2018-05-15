module('core', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Creating a view with an empty text widget', function() {
    visit('/userview/view.event');

    expect(2);

    click('.main-tabs a.dropdown-toggle');
    click('.main-tabs .fa.fa-plus');

    waitForElement('input[name=crecord_name]').then(function(){
        fillIn('input[name=crecord_name]', 'test');
        click('.modal-dialog .btn-primary');
    });

    activateEditMode();
    click('.btn-add-widget');

    waitForElement('.modal-dialog .ember-text-field').then(function(){
        equal(find('.box-title').length, 0, 'No widget on the view');
        fillIn('.modal-dialog .ember-text-field', 'text');
        click('.modal-dialog .panel-default:first a');
        click('.modal-dialog .list-group-item a');
        click('.modal-dialog .btn-primary');
        click('.modal-dialog .btn-primary');
        waitForElement('.box-title').then(function(){
            equal(find('.box-title').text(), "< Untitled text widget >", 'an untitled text widget is present');
        });
    });
});
