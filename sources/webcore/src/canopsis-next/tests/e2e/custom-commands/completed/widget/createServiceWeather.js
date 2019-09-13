// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

module.exports.command = function createServiceWeather(
  {
    parameters: {
      blockTemplate,
      modalTemplate,
      entityTemplate,
      ...parameters
    },
    ...fields
  },
  callback = () => {},
) {
  const weather = this.page.widget.weather();
  const textEditorModal = this.page.modals.common.textEditorModal();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters,
  });

  if (blockTemplate) {
    weather.clickTemplateWeatherItem();

    textEditorModal.verifyModalOpened()
      .clickField()
      .setField(blockTemplate)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (modalTemplate) {
    weather.clickTemplateModal();

    textEditorModal
      .verifyModalOpened()
      .clickField()
      .setField(modalTemplate)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  if (entityTemplate) {
    weather.clickTemplateEntities();

    textEditorModal
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
