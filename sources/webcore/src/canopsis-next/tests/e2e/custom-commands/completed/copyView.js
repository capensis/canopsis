// http://nightwatchjs.org/guide#usage
const { VIEW } = require('../../constants');
const { API_ROUTES } = require('../../../../src/config');

module.exports.command = function copyView(groupId, viewId, view = { ...VIEW }, callback = () => {}) {
  const {
    name,
    title,
    description,
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

  groupsSideBar.clickPanelHeader(groupId)
    .verifyPanelBody(groupId)
    .clickCopyViewButton(viewId)
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .setViewName(name)
    .setViewTitle(title)
    .setViewDescription(description)
    .clickViewEnabled()
    .setViewGroupTags(group)
    .setViewGroupIds(group);

  this.waitForFirstXHR(
    new RegExp(`${API_ROUTES.view}$`),
    5000,
    () => modalViewCreate.clickViewSubmitButton(),
    ({ responseData, requestData }) => callback({ ...JSON.parse(requestData), ...JSON.parse(responseData) }),
  );

  modalViewCreate.verifyModalClosed();

  navigation.clickEditModeButton()
    .clickSettingsViewButton();

  return this;
};
