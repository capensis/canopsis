import { camelCase } from 'lodash';

const requireModule = require.context('.', true, /\.js$/);

const messages = requireModule.keys().reduce((acc, fileName) => {
  if (fileName.includes('index.js')) {
    return acc;
  }

  const [, language, moduleName] = fileName.match(/\/(.+)\/(.+).js$/);

  if (!acc[language]) {
    acc[language] = {};
  }

  acc[language][camelCase(moduleName)] = requireModule(fileName).default;

  return acc;
}, {});

export default messages;
