// http://nightwatchjs.org/guide#usage

module.exports.command = function createView(
  view = {
    name: 'test-name',
    title: 'test-title',
    description: 'test-description',
    tags: 'test-tags',
  },
  callback = result => result,
) {
  const {
    name = 'test-name',
    title = 'test-title',
    description = 'test-description',
    tags = 'test-tags',
  } = view;

  const viewPath = `${process.env.VUE_DEV_SERVER_URL}view/`;

  const topBar = this.page.layout.topBar();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const leftSideBar = this.page.layout.leftSideBar();
  const modalViewCreate = this.page.modals.view.create();

  this.url(viewPath);

  topBar.clickUserDropdown()
    .clickUserProfileButton();

  createUser.verifyModalOpened()
    .selectNavigationType(1)
    .clickSubmitButton()
    .verifyModalClosed();

  groupsSideBar.clickGroupsSideBarButton();

  leftSideBar.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickAddViewButton()
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .setViewName(name)
    .setViewTitle(title)
    .setViewDescription(description)
    .clickViewEnabled()
    .setViewGroupTags(tags)
    .setViewGroupIds(tags);

  this.waitForFirstXHR(
    '/api/v2/views',
    1000,
    () => {
      modalViewCreate.clickViewSubmitButton();
    },
    (xhr) => {
      const viewResponseData = JSON.parse(xhr.responseData);
      callback({ viewResponseData, view });
    },
  );

  modalViewCreate.verifyModalClosed();

  return this;
};
