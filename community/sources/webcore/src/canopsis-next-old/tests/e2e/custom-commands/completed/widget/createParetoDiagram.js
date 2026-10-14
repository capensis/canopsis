// http://nightwatchjs.org/guide#usage

const { API_ROUTES } = require('@/config');
const { WAIT_FOR_FIRST_XHR_TIME } = require('../../../constants');

module.exports.command = function createParetoDiagram(fields, callback = () => {}) {
  const paretoDiagramWidget = this.page.widget.paretoDiagram();

  this.completed.widget.setCommonFields(fields);

  this.waitForFirstXHR(
    API_ROUTES.userPreferences,
    WAIT_FOR_FIRST_XHR_TIME,
    () => paretoDiagramWidget.clickSubmitParetoDiagram(),
    ({ responseData, requestData }) => callback({
      response: JSON.parse(responseData),
      request: JSON.parse(requestData),
    }),
  );
};
