// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createServiceWeather(
  {
    parameters: {
      filter,
      moreInfos,
      blockTemplate,
      modalTemplate,
      entityTemplate,
      newColumnNames,
      editColumnNames,
      moveColumnNames,
      deleteColumnNames,
      ...parameters
    },
    ...fields
  },
  callback = () => {},
) {
  const common = this.page.widget.common();
  const weather = this.page.widget.weather();
  const textEditor = this.page.modals.common.textEditor();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters,
  });

  if (newColumnNames || editColumnNames || moveColumnNames || deleteColumnNames) {
    common.clickColumnNames();
  }

  if (newColumnNames) {
    newColumnNames.forEach(({ index, data }) => {
      common
        .clickAddColumnName()
        .editColumnName(index, data);
    });
  }

  if (editColumnNames) {
    editColumnNames.forEach(({ index, data }) => {
      common.editColumnName(index, data);
    });
  }

  if (moveColumnNames) {
    moveColumnNames.forEach(({ index, up, down }) => {
      if (up) {
        common.clickColumnNameUpWard(index);
      }

      if (down) {
        common.clickColumnNameDownWard(index);
      }
    });
  }

  if (deleteColumnNames) {
    deleteColumnNames.forEach((index) => {
      common.clickDeleteColumnName(index);
    });
  }

  if (filter) {
    common.clickCreateFilter();

    this.page.modals.view.createFilter()
      .verifyModalOpened()
      .clickCancelButton()
      .verifyModalClosed();
  }

  if (moreInfos) {
    common.clickCreateMoreInfos();

    textEditor.verifyModalOpened()
      .clickField()
      .setField(moreInfos)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (blockTemplate) {
    weather.clickTemplateWeatherItem();

    textEditor.verifyModalOpened()
      .clickField()
      .setField(blockTemplate)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (modalTemplate) {
    weather.clickTemplateModal();

    textEditor
      .verifyModalOpened()
      .clickField()
      .setField(modalTemplate)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (entityTemplate) {
    weather.clickTemplateEntities();

    textEditor
      .verifyModalOpened()
      .clickField()
      .setField(entityTemplate)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => weather.clickSubmitWeather(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
