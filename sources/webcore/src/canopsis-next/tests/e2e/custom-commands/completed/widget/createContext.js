// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function createContext({
  parameters: {
    typeOfEntities,
    ...parameters
  } = {},
  ...fields
}, callback = () => {}) {
  const contextWidget = this.page.widget.context();

  this.completed.widget.setCommonFields({
    ...fields,
    parameters: {
      advanced: true,
      ...parameters,
    },
  });

  if (typeOfEntities) {
    contextWidget.clickContextTypeOfEntities();

    typeOfEntities.forEach(({ index, value }) => {
      contextWidget.selectEntitiesTypeCheckbox(index, value);
    });
  }

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    WAIT_FOR_FIRST_XHR_TIME,
    () => contextWidget.clickSubmitContext(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
