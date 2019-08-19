// http://nightwatchjs.org/guide#usage

const { CREDENTIALS } = require('../../constants');
const { command: loginCommand } = require('./login');

module.exports.command = function loginAsAdmin() {
  const { username, password } = CREDENTIALS.admin;

  return loginCommand.call(this, username, password);
};
