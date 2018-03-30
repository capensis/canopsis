module('canopsis-rights', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Test account creation', function() {
    expect(2);

    visit('/userview/view.accounts');

    click('.widget .glyphicon-plus-sign');

    waitForElement('.modal-dialog .ember-text-field').then(function() {
        fillIn('input[name=_id]', 'test_account');
        fillIn('input[name=firstname]', 'test');
        fillIn('input[name=lastname]', 'account');
        fillIn('input[name=mail]', 'account@test.com');
        fillIn('input[type=password]', 'aabbcc');
        click('.modal-dialog .btn-primary');
        waitForElementRemoval('.modal-backdrop').then(function() {
            fillIn('input[placeholder=Search]', 'test');
            click('.fa-search');
            waitMilliseconds(1500).then(function() {
                equal(find('td._id').html(), 'test_account', 'The account has successfully been created');
                click('.widget td .glyphicon-trash');
                waitForElement('.modal-dialog .btn-primary').then(function() {
                    click('.modal-dialog .btn-primary');
                });
                waitForElementRemoval('.modal-backdrop').then(function() {
                    waitForElementRemoval('.widget .fa-refresh').then(function() {
                        equal(find('td._id').html().contains('test'), false, 'The account has successfully been deleted');
                    });
                });
            });
        });
    });
});
