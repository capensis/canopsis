
const ELEMENTS_WAITING_DELAY = 5000;

exports.command = function (browser, url, username, password) {
  browser
    .url(url)
    .waitForElementVisible('form', ELEMENTS_WAITING_DELAY)
    .setValue('input[name=username]', username)
    .setValue('input[type=password]', password)
    .click('button[type=submit]')
    .pause(ELEMENTS_WAITING_DELAY)
    .waitForElementVisible('.v-toolbar__content .v-btn__content', ELEMENTS_WAITING_DELAY)
    .assert.containsText('.v-toolbar__content .v-btn__content', 'menu')
    .end();
};

