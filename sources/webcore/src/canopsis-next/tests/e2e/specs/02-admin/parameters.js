// http://nightwatchjs.org/guide#usage
const uid = require('uid');

const { LANGUAGES_POSITIONS } = require('../../constants');
const { generateTemporaryUser } = require('../../helpers/entities');

module.exports = {
  async before(browser, done) {
    browser.globals.user = {};

    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.completed.logout()
      .end(done);
  },

  'Create temporary user without interface language': (browser) => {
    const user = generateTemporaryUser();

    browser.completed.user.create(user, (createdUser) => {
      browser.globals.user = {
        ...createdUser,

        password: user.password,
      };
    });
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

  'Delete temporary user': (browser) => {
    const { user } = browser.globals;

    browser.completed.loginAsAdmin();
    browser.completed.user.delete(user._id, () => {
      delete browser.globals.user;
    });
  },
};
