// http://nightwatchjs.org/guide#usage

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Browse view by id': (browser) => {
    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .browseGroupById('05b2e049-b3c4-4c5b-94a5-6e7ff142b28c')
      .browseViewById('da7ac9b9-db1c-4435-a1f2-edb4d6be4db8'); // TODO: put verification
  },
};
