// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function createText(
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
    WAIT_FOR_FIRST_XHR_TIME,
    () => text.clickSubmitText(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
