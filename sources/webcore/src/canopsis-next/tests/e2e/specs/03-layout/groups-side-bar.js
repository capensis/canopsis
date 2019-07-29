// http://nightwatchjs.org/guide#usage

const { NAVIGATION: { groups } } = require('../../constants');

const TEMPORARY_DATA = {};

const createTemporaryObject = ({ prefix, text, index }) => {
  const i = typeof index === 'number' ? `-${index}` : '';
  const r = Math.random().toString(36).substring(7);
  return {
    name: `${prefix}-${text}-name${i}-${r}`,
    title: `${prefix}-${text}-title${i}-${r}`,
    description: `${prefix}-${text}-description${i}-${r}`,
    group: `${prefix}-${text}-group${i}-${r}`,
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

  'Add view with some name from constants': (browser) => {
    const { text, create: { prefix } } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    browser.completed.createView(TEMPORARY_DATA[prefix], (view) => {
      TEMPORARY_DATA[prefix] = view;
    });
  },

  'Checking view copy with name from constants': (browser) => {
    const { text, copy: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    browser.completed.copyView(
      TEMPORARY_DATA[create.prefix].group_id,
      TEMPORARY_DATA[create.prefix]._id,
      TEMPORARY_DATA[prefix],
      (view) => {
        TEMPORARY_DATA[prefix] = view;
      },
    );
  },

  'Editing test view with name from constants': (browser) => {
    const { text, edit: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name,
      title,
      description,
      group,
    } = TEMPORARY_DATA[prefix];

    const r = Math.random().toString(36).substring(7);

    browser.page.layout.navigation()
      .verifySettingsWrapperBefore()
      .clickSettingsViewButton()
      .verifyControlsWrapperBefore()
      .clickEditModeButton()
      .defaultPause();

    browser.page.layout.groupsSideBar()
      .clickEditGroupButton(TEMPORARY_DATA[create.prefix].group_id)
      .defaultPause();

    TEMPORARY_DATA[create.prefix].group = `${create.prefix}-${text}-group-${r}`;

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clearGroupName()
      .setGroupName(`${create.prefix}-${text}-group-${r}`)
      .clickSubmitButton()
      .verifyModalClosed();


    browser.page.layout.groupsSideBar()
      .verifyPanelBody(TEMPORARY_DATA[create.prefix].group_id)
      .clickEditViewButton(TEMPORARY_DATA[create.prefix]._id)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .clearViewName()
      .setViewName(name)
      .clearViewTitle()
      .setViewTitle(title)
      .clearViewDescription()
      .setViewDescription(description)
      .clickViewEnabled()
      .clearViewGroupTags()
      .setViewGroupTags(group)
      .clearViewGroupIds()
      .setViewGroupIds(group)
      .clickViewSubmitButton()
      .verifyModalClosed();

    browser.page.layout.navigation()
      .clickEditModeButton()
      .clickSettingsViewButton();
  },

  'Deleting all test items view with name from constants': (browser) => {
    const { create, copy } = groups;

    browser.completed.deleteView(TEMPORARY_DATA[create.prefix].group_id, TEMPORARY_DATA[create.prefix]._id);
    browser.completed.deleteView(TEMPORARY_DATA[copy.prefix].group_id, TEMPORARY_DATA[copy.prefix]._id);
  },
};
