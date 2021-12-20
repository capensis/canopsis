const { merge } = require('lodash');
const { process, getCacheKey } = require('vue-jest');

const updateFieldDirective = require('../../tools/update-field-directive');

module.exports = {
  process: (src, filename, config) => {
    const configWithDirectives = merge(config, {
      globals: {
        'vue-jest': {
          templateCompiler: {
            compilerOptions: {
              directives: {
                field: updateFieldDirective,
              },
            },
          },
        },
      },
    });

    return process(src, filename, configWithDirectives);
  },
  getCacheKey,
};
