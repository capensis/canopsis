// http://nightwatchjs.org/guide#usage

module.exports.command = function refreshPage(route, action, responseCallback, timeout = 5000) {
  this.waitForFirstXHR(
    route,
    timeout,
    action,
    responseCallback || (({ status, method, httpResponseCode }) => {
      this.assert.equal(status, 'success');
      this.assert.equal(method, 'GET');
      this.assert.equal(httpResponseCode, '200');
    }),
  );

  return this;
};
