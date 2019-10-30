// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function editView(groupId, viewId, view, callback = () => {}) {
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
    .clickEditViewButton(viewId)
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
      .setViewGroupTags(tags);
  }

  if (group) {
    modalViewCreate.clearViewGroupId()
      .setViewGroupId(group);
  }

  modalViewCreate.setViewEnabled(enabled);

  this.waitForFirstXHR(
    `${API_ROUTES.view}/${viewId}`,
    WAIT_FOR_FIRST_XHR_TIME,
    () => modalViewCreate.clickViewSubmitButton(),
    ({ responseData, requestData }) => {
      modalViewCreate.verifyModalClosed();

      callback({ ...JSON.parse(requestData), ...JSON.parse(responseData) });
    },
  );

  return this;
};
