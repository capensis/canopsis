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

  'Browse view by id': (browser) => {
    browser.page.layout.topBar()
      .clickUserDropdown()
      .clickUserProfileButton();

    browser.page.modals.admin.createUser()
      .verifyModalOpened()
      .selectNavigationType(1)
      .clickSubmitButton()
      .verifyModalClosed();

    browser.page.layout.groupsSideBar()
      .clickGroupsSideBarButton()
      .browseGroupById('05b2e049-b3c4-4c5b-94a5-6e7ff142b28c') // TODO: use from some constants file when we will use fixtures
      .browseViewById('da7ac9b9-db1c-4435-a1f2-edb4d6be4db8')
      .defaultPause(); // TODO: put verification
  },

  'Add view with some name from constants': (browser) => {
    const { text, create: { prefix } } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    browser.completed.createView(TEMPORARY_DATA[prefix], ({ viewResponseData }) => {
      TEMPORARY_DATA[prefix] = {
        ...TEMPORARY_DATA[prefix],
        viewResponseData,
      };
    });
  },

  'Checking view copy with name from constants': (browser) => {
    const { text, copy: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description,
    } = TEMPORARY_DATA[prefix];

    browser.page.layout.navigation()
      .verifyControlsWrapperBefore()
      .clickEditViewButton()
      .defaultPause();

    browser.page.layout.groupsSideBar()
      .clickPanelHeader(TEMPORARY_DATA[create.prefix].group)
      .verifyPanelBody(TEMPORARY_DATA[create.prefix].group)
      .clickCopyViewButton(TEMPORARY_DATA[create.prefix].title)
      .defaultPause();

    browser.page.modals.view.create()
      .verifyModalOpened()
      .setViewName(name)
      .setViewTitle(title)
      .clearViewDescription()
      .setViewDescription(description)
      .clickViewSubmitButton()
      .verifyModalClosed();
  },

  'Editing test view with name from constants': (browser) => {
    const { text, edit: { prefix }, create } = groups;

    TEMPORARY_DATA[prefix] = createTemporaryObject({ prefix, text });

    const {
      name, title, description, group,
    } = TEMPORARY_DATA[prefix];

    const r = Math.random().toString(36).substring(7);

    browser.page.layout.groupsSideBar()
      .clickEditGroupButton(TEMPORARY_DATA[create.prefix].group)
      .defaultPause();

    TEMPORARY_DATA[create.prefix].group = `${create.prefix}-${text}-group-${r}`;

    browser.page.modals.view.createGroup()
      .verifyModalOpened()
      .clearGroupName()
      .setGroupName(TEMPORARY_DATA[create.prefix].group)
      .clickSubmitButton()
      .verifyModalClosed();


    browser.page.layout.groupsSideBar()
      .verifyPanelBody(TEMPORARY_DATA[create.prefix].group)
      .clickEditViewButton(TEMPORARY_DATA[create.prefix].title)
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
  },

  'Deleting all test items view with name from constants': (browser) => {
    const { create, edit, copy } = groups;

    browser.completed.deleteView({
      tags: TEMPORARY_DATA[create.prefix].group,
      title: TEMPORARY_DATA[copy.prefix].title,
    });

    browser.completed.deleteView({
      tags: TEMPORARY_DATA[edit.prefix].group,
      title: TEMPORARY_DATA[edit.prefix].title,
    });
  },
};
