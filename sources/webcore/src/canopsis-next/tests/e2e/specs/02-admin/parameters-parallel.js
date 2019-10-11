// http://nightwatchjs.org/guide#usage
const uid = require('uid');

const { LANGUAGES_POSITIONS } = require('../../constants');
const { createUser, removeUser } = require('../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.globals.user = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    const { data: user } = await createUser({
      ui_language: undefined,
    });

    browser.globals.user = {
      ...user,
      username: user._id,
    };

    done();
  },

  async after(browser, done) {
    browser.end();

    await removeUser(browser.globals.user.username);
    delete browser.globals.user;

    done();
  },

  'Change parameters with some name': (browser) => {
    const hash = uid();

    browser.globals.hash = hash;
    browser.page.admin.parameters()
      .navigate()
      .verifyPageElementsBefore()
      .clearAppTitle()
      .setAppTitle(hash)
      .selectLanguage(LANGUAGES_POSITIONS.fr)
      .clearFooter()
      .setFooter(hash)
      .clearDescription()
      .setDescription(hash)
      .clickSubmitButton();
  },

  'Check parameters app title on page': (browser) => {
    const { hash } = browser.globals;

    browser.page.admin.parameters()
      .verifyAppTitle(hash);
  },

  'Check parameters login values on page': (browser) => {
    const { hash } = browser.globals;

    browser.completed.logout();

    browser.page.admin.parameters()
      .verifyLoginDescription(hash)
      .verifyLoginFooter(hash);
  },

  'Caching of temporary user first interface language': (browser) => {
    const { user } = browser.globals;

    browser.completed.login(user._id, user.password);

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickUserDropdown()
      .getLogoutButtonText(({ value }) => {
        browser.globals.logoutButtonText = value;
      })
      .clickLogoutButton();

    browser.page.auth.logout()
      .verifyPageElementsAfter();
  },

  'Change parameters global language': (browser) => {
    browser.completed.loginAsAdmin();

    browser.page.admin.parameters()
      .navigate()
      .selectLanguage(LANGUAGES_POSITIONS.en)
      .clickSubmitButton();

    browser.completed.logout();
  },

  'Check temporary user second interface language': (browser) => {
    const { user } = browser.globals;
    const { logoutButtonText } = browser.globals;

    browser.completed.login(user._id, user.password);

    browser.page.layout.popup()
      .clickOnEveryPopupsCloseIcons();

    browser.page.layout.topBar()
      .clickUserDropdown()
      .getLogoutButtonText(({ value }) => {
        browser.assert.notEqual(value, logoutButtonText);
      })
      .clickLogoutButton();

    browser.page.auth.logout()
      .verifyPageElementsAfter();
  },
};
