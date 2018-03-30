module('core', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Creating two widgets and reorder them', function() {
    $('.modal-backdrop').remove();

    visit('/userview/view.event');

    expect(4);

    createNewView('widget reorder test');

    activateEditMode();

    click('.btn-add-widget');


    waitForElement('.form .ember-text-field').then(function(){
        fillIn('.form .ember-text-field', 'text');

        click('.form .panel-default:first a');
        click('.form .list-group-item a');
        fillIn('input[name=title]', 'widget1');
        click('.form .btn-submit');

        waitForElementRemoval('.modal-backdrop').then(function() {
            click('.btn-add-widget');

            waitForElement('.form .ember-text-field').then(function(){
                fillIn('.form .ember-text-field', 'text');
                click('.form .panel-default:first a');
                click('.form .list-group-item a');
                fillIn('input[name=title]', 'widget2');
                click('.form .btn-submit');

                waitForElementRemoval('.modal-backdrop').then(function() {
                    equal(find('.box-title').text(), "widget1widget2", 'widget1 then widget2');

                    click('.btn-move-down:first');
                    equal(find('.box-title').text(), "widget2widget1", 'widget2 then widget1');
                    click('.btn-move-up:last');
                    waitMilliseconds(500).then(function() {
                        equal(find('.box-title').text(), "widget1widget2", 'widget1 then widget2');

                        click('.btn-move-up:first');
                        click('.btn-move-down:last');
                        equal(find('.box-title').text(), "widget1widget2", 'widget1 then widget2');
                    });
                });
            });
        });
    });

});
