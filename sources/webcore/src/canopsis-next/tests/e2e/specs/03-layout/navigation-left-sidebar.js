// http://nightwatchjs.org/guide#usage

const { NAVIGATION: { LEFT_SIDEBAR } } = require('../../constants');

const TEMPORARY_DATA = {};

const onCreateTemporaryObject = ({ prefix, text, index }) => {
  const i = typeof index === 'number' ? `-${index}` : '';
  return {
    name: `${prefix}-${text}-name${i}`,
    title: `${prefix}-${text}-title${i}`,
    description: `${prefix}-${text}-description${i}`,
    tags: `${prefix}-${text}-tags${i}`,
  };
};

module.exports = {
  async before(browser, done) {
    await browser.maximizeWindow()
      .completed.loginAsAdmin();

    done();
  },

  after(browser, done) {
    browser.end(done);
  },

  'Open sidebar': (browser) => {
    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .defaultPause();
  },

  'Add view with some name': (browser) => {
    const { text, create: { prefix } } = LEFT_SIDEBAR;

    TEMPORARY_DATA[prefix] = onCreateTemporaryObject({ prefix, text });

    const {
      name, title, description, tags,
    } = TEMPORARY_DATA[prefix];

    browser.page.layout.leftSideBar()
      .verifySettingsWrapperBefore()
      .clickSettingsViewButton()
      .verifyControlsWrapperBefore()
      .clickAddViewButton()
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .setViewName(name)
      .setViewTitle(title)
      .setViewDescription(description)
      .clickViewEnabled()
      .setViewGroupTags(tags)
      .setViewGroupIds(tags)
      .clickViewSubmitButton()
      .verifyModalClosed();
  },
};
