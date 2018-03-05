module('canopsis-rights', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Test role creation', function() {
    expect(2);

    visit('/userview/view.permissions');

    click('.widget .glyphicon-plus-sign');

    waitForElement('.modal-dialog .ember-text-field').then(function() {
        fillIn('input[name=_id]', 'test_role');
        fillIn('input[name=description]', 'test description');
        click('.modal-dialog .btn-primary');
        waitForElementRemoval('.modal-backdrop').then(function() {
            fillIn('input[placeholder=Search]', 'test_role');
            click('.fa-search');
            waitMilliseconds(1500).then(function() {
                equal(find('td._id').html(), 'test_role', 'The role has successfully been created');
                click('.widget td .glyphicon-trash');
                waitForElement('.modal-dialog .btn-primary').then(function() {
                    click('.modal-dialog .btn-primary');
                });
                waitForElementRemoval('.modal-backdrop').then(function() {
                    waitMilliseconds(1500).then(function() {
                        equal(find('td._id').length, 0, 'The account has successfully been deleted');
                    });
                });
            });
        });
    });
});
