// http://nightwatchjs.org/guide#usage

const { CREDENTIALS } = require('../../constants');

module.exports.command = function loginAsAdmin() {
  const { username, password } = CREDENTIALS.admin;

  return this.completed.login(username, password);
};
