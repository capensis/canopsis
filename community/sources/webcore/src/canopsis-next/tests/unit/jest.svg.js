const path = require('path');

module.exports = {
  process: (src, filename, options) => {
    const rootDir = options.rootDir || options.config.rootDir;

    const exportedPath = JSON.stringify(path.relative(rootDir, filename));

    return `module.exports = ${exportedPath}`;
  },
};
