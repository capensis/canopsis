const resolveClientEnv = require('@vue/cli-service/lib/util/resolveClientEnv'); // eslint-disable-line import/no-extraneous-dependencies

module.exports = function loadEnv(path) {
  try {
    resolveClientEnv(path);
  } catch (err) {
    if (err.toString().indexOf('ENOENT') < 0) {
      console.warn(err);
    }
  }
};
