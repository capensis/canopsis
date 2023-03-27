const path = require('path');

const loadEnv = require('./tools/load-env');

const localEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env.local');
const baseEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env');

loadEnv(localEnvPath);
loadEnv(baseEnvPath);

module.exports = {
  moduleFileExtensions: [
    'js',
    'jsx',
    'json',
    'vue',
    'styl',
  ],
  testEnvironment: 'jest-environment-jsdom',
  transform: {
    '^.+\\.vue$': '<rootDir>/tests/unit/jest.vue',
    '^.+\\.jsx?$': '<rootDir>/tests/unit/jest.transform',
    '^.+\\.svg$': '<rootDir>/tests/unit/jest.svg',
  },
  transformIgnorePatterns: [
    '<rootDir>/node_modules/(?!(vue-tour|monaco-editor|dayspan-vuetify/src)/.*)',
  ],
  moduleNameMapper: {
    '^.+\\.styl(us)?$': '<rootDir>/tests/unit/mocks/styleMock.js',
    '^.+\\.css$': '<rootDir>/tests/unit/mocks/styleMock.js',
    '^@unit/(.*)$': '<rootDir>/tests/unit/$1',
    '^@/(.*)$': '<rootDir>/src/$1',
    '^vue$': 'vue/dist/vue.common.dev.js',
    '^mermaid$': '<rootDir>/node_modules/mermaid/dist/mermaid.js',
    '^monaco-mermaid$': '<rootDir>/node_modules/mermaid/dist/mermaid.js',
    './assets': '<rootDir>/tests/unit/mocks/flowchartAssets.js',
  },
  snapshotSerializers: [
    'jest-serializer-vue',
  ],
  setupFiles: ['jest-localstorage-mock'],
  setupFilesAfterEnv: ['<rootDir>/tests/unit/jest.setup-test-framework'],
  globalSetup: '<rootDir>/tests/unit/jest.global-setup',
  maxWorkers: process.env.JEST_MAX_WORKERS ?? '10%',
};
