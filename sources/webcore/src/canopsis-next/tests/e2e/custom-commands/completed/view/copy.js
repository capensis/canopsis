// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function copyView(groupId, viewId, view, callback = () => {}) {
  const {
    name,
    title,
    description,
    enabled,
    tags,
    group,
  } = view;

  const topBar = this.page.layout.topBar();
  const navigation = this.page.layout.navigation();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const modalViewCreate = this.page.modals.view.create();

  groupsSideBar.groupsSideBarButtonElement(({ status }) => {
    if (status === -1) {
      topBar.clickUserDropdown()
        .clickUserProfileButton();

      createUser.verifyModalOpened()
        .selectNavigationType(1)
        .clickSubmitButton()
        .verifyModalClosed();
    }

    groupsSideBar.clickGroupsSideBarButton();
  });

  navigation.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickEditModeButton()
    .defaultPause();

  groupsSideBar.browseGroupById(groupId)
    .verifyPanelBody(groupId)
    .clickCopyViewButton(viewId)
    .defaultPause();

  modalViewCreate.verifyModalOpened();

  if (name) {
    modalViewCreate.clearViewName()
      .setViewName(name);
  }

  if (title) {
    modalViewCreate.clearViewTitle()
      .setViewTitle(title);
  }

  if (description) {
    modalViewCreate.clearViewDescription()
      .setViewDescription(description);
  }

  if (tags) {
    modalViewCreate.clearViewGroupTags()
      .setViewGroupTags(description);
  }

  if (group) {
    modalViewCreate.clearViewGroupId()
      .setViewGroupId(group);
  }

  modalViewCreate.setViewEnabled(enabled);

  this.waitForFirstXHR(
    new RegExp(`${API_ROUTES.view}$`),
    5000,
    () => modalViewCreate.clickViewSubmitButton(),
    ({ responseData, requestData }) => {
      modalViewCreate.verifyModalClosed();

      callback({ ...JSON.parse(requestData), ...JSON.parse(responseData) });
    },
  );

  return this;
};
