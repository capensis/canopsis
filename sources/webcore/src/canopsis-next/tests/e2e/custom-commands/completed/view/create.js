// http://nightwatchjs.org/guide#usage
const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createView(view, callback = () => {}) {
  const {
    name,
    title,
    description,
    enabled,
    tags,
    group,
  } = view;

  const topBar = this.page.layout.topBar();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const navigation = this.page.layout.navigation();
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
    .clickAddViewButton()
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .setViewName(name)
    .setViewTitle(title)
    .setViewDescription(description)
    .setViewEnabled(enabled)
    .setViewGroupTags(tags)
    .setViewGroupId(group);

  this.waitForFirstXHR(
    new RegExp(`${API_ROUTES.view}$`),
    5000,
    () => modalViewCreate.clickViewSubmitButton(),
    ({ responseData, requestData }) => callback({ ...JSON.parse(requestData), ...JSON.parse(responseData) }),
  );

  modalViewCreate.verifyModalClosed();

  return this;
};
