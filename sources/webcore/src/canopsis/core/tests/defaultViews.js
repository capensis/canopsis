module('default views', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Test app_header', function() {

    visit('/userview/view.app_header');

    andThen(function() {
        equal(find('.tab-content .nav-tabs').length, 1, 'There is a bootstrap nav-tabs in the view');
    });
});

test('Test app_footer', function() {

    visit('/userview/view.app_footer');

    andThen(function() {
        equal(
            find('.tab-content img:first')[0].src.indexOf('static/canopsis/media/sakura.png') !== -1,
            true,
            'Canopsis logo is present'
        );
    });
});
