// http://nightwatchjs.org/guide#usage
const { VIEW } = require('../../constants');
const { API_ROUTES } = require('../../../../src/config');
const fs = require('fs');

module.exports.command = function createView(
  view = VIEW,
  callback = result => fs.writeFile('view.json', JSON.stringify(result), 'utf8', () => {}),
) {
  const {
    name = VIEW.name,
    title = VIEW.title,
    description = VIEW.description,
    group = VIEW.group,
  } = view;

  const topBar = this.page.layout.topBar();
  const createUser = this.page.modals.admin.createUser();
  const groupsSideBar = this.page.layout.groupsSideBar();
  const navigation = this.page.layout.navigation();
  const modalViewCreate = this.page.modals.view.create();

  this.refresh();

  topBar.clickUserDropdown()
    .clickUserProfileButton();

  createUser.verifyModalOpened()
    .selectNavigationType(1)
    .clickSubmitButton()
    .verifyModalClosed();

  groupsSideBar.clickGroupsSideBarButton();

  navigation.verifySettingsWrapperBefore()
    .clickSettingsViewButton()
    .verifyControlsWrapperBefore()
    .clickAddViewButton()
    .defaultPause();

  modalViewCreate.verifyModalOpened()
    .setViewName(name)
    .setViewTitle(title)
    .setViewDescription(description)
    .clickViewEnabled()
    .setViewGroupTags(group)
    .setViewGroupIds(group);

  this.waitForXHR(
    API_ROUTES.view,
    10000,
    () => modalViewCreate.clickViewSubmitButton(),
    (xhr) => {
      const responseGroupId = JSON.parse(xhr[0].responseData)._id;
      const responseViewID = JSON.parse(xhr[1].responseData)._id;
      const responseViewTabId = JSON.parse(xhr[2].responseData).groups[responseGroupId].views
        .filter(item => item._id === responseViewID)[0].tabs[0].tabs._id;
      const responseData = {
        responseGroupId,
        responseViewID,
        responseViewTabId,
      };
      callback(responseData);
    },
  );

  modalViewCreate.verifyModalClosed();

  return this;
};
