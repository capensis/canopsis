// For authoring Nightwatch tests, see
// http://nightwatchjs.org/guide#usage

const ELEMENTS_WAITING_DELAY = 5000;

module.exports.command = function (username, password) {
  this
    .url(process.env.VUE_DEV_SERVER_URL)
    .waitForElementVisible('form', ELEMENTS_WAITING_DELAY)
    .setValue('input[name=username]', username)
    .setValue('input[type=password]', password)
    .click('button[type=submit]')
    .pause(ELEMENTS_WAITING_DELAY)
    .waitForElementVisible('.v-toolbar__content .v-btn__content', ELEMENTS_WAITING_DELAY)
    .assert.containsText('.v-toolbar__content .v-btn__content', 'menu');
  return this;
};

