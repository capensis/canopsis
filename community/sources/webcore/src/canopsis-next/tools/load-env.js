const { config } = require('dotenv'); // eslint-disable-line import/no-extraneous-dependencies

module.exports = function loadEnv(path) {
  try {
    config({ path });
  } catch (err) {
    if (err.toString().indexOf('ENOENT') < 0) {
      console.warn(err);
    }
  }
};
