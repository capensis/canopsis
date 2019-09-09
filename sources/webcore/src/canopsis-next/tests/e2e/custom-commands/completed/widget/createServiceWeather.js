// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createServiceWeather(
  {
    parameters: {
      filter,
      blockTemplate,
      modalTemplate,
      entityTemplate,
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

  if (filter) {
    common.clickCreateFilter();

    this.page.modals.view.createFilter()
      .verifyModalOpened()
      .clickCancelButton()
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
