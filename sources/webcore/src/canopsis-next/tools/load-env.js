const loadEnvLib = require('@vue/cli-service/lib/util/loadEnv'); // eslint-disable-line import/no-extraneous-dependencies

module.exports = function loadEnv(path) {
  try {
    loadEnvLib(path);
  } catch (err) {
    if (err.toString().indexOf('ENOENT') < 0) {
      console.warn(err);
    }
  }
};
