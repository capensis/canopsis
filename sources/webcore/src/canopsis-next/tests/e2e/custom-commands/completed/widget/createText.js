// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');

module.exports.command = function createServiceWeather(
  {
    parameters: {
      template,
      ...parameters
    },
    ...fields
  },
  callback = () => {},
) {
  const text = this.page.widget.text();
  const textEditorModal = this.page.modals.common.textEditor();

  text.clickStats();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters,
  });

  if (template) {
    text.clickCreateTemplate();

    textEditorModal.verifyModalOpened()
      .clickField()
      .setField(template)
      .clickSubmitButton()
      .verifyModalClosed();
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    5000,
    () => text.clickSubmitText(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
