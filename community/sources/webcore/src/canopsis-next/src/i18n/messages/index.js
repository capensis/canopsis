import { camelCase } from 'lodash';

const modules = import.meta.glob('./**/*.js', {
  eager: true,
});

const messages = Object.entries(modules).reduce((acc, [fileName, module]) => {
  if (fileName.includes('index.js')) {
    return acc;
  }

  const [, language, moduleName] = fileName.match(/\/(.+)\/(.+).js$/);

  if (!acc[language]) {
    acc[language] = {};
  }

  acc[language][camelCase(moduleName)] = module.default;

  return acc;
}, {});

export default messages;
