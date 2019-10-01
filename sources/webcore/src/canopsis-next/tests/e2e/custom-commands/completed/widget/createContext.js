// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('../../../../../src/config');

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
    5000,
    () => contextWidget.clickSubmitContext(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
