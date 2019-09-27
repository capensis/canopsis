// http://nightwatchjs.org/guide#usage

const { createAdminUser, removeUser } = require('../../helpers/api');

module.exports = {
  async before(browser, done) {
    browser.maximizeWindow();

    const { data } = await createAdminUser();

    browser.globals.credentials = {
      password: data.password,
      username: data._id,
    };

    done();
  },

  async after(browser, done) {
    await removeUser(browser.globals.credentials.username);
    await browser.end();

    delete browser.globals.credentials;

    done();
  },

  'Correct user credentials login': (browser) => {
    const { username, password } = browser.globals.credentials;

    browser.completed.login(username, password);
  },

  'Authorized user logout': (browser) => {
    browser.completed.logout();
  },
};
