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
    '^.+\\.html?$': 'html-loader-jest',
    '^.+\\.(jpg|jpeg|png)$': '<rootDir>/tests/unit/jest.assets.js',
  },
  transformIgnorePatterns: [
    '<rootDir>/node_modules/(?!(vue-tour|monaco-editor|vuetify/lib)/.*)',
  ],
  moduleNameMapper: {
    '^.+\\.styl(us)?$': '<rootDir>/tests/unit/mocks/styleMock.js',
    '^.+\\.(scss|sass|css)$': '<rootDir>/tests/unit/mocks/styleMock.js',
    '^@unit/(.*)$': '<rootDir>/tests/unit/$1',
    '^@/(.*)$': '<rootDir>/src/$1',
    '^vue$': 'vue/dist/vue.common.dev.js',
    '^mermaid$': '<rootDir>/node_modules/mermaid/dist/mermaid.js',
    '^monaco-mermaid$': '<rootDir>/node_modules/mermaid/dist/mermaid.js',
    '^./assets$': '<rootDir>/tests/unit/mocks/flowchartAssets.js',
    '@/assets/images/engineering.svg': '<rootDir>/tests/unit/mocks/flowchartAssets.js',
    './components/icons': '<rootDir>/tests/unit/mocks/vuetifyIcons.js',
  },
  snapshotSerializers: ['<rootDir>/tests/unit/jest.serializer-vue'],
  collectCoverageFrom: ['<rootDir>/src/**/*.{js,vue}'],
  setupFiles: ['jest-localstorage-mock'],
  setupFilesAfterEnv: ['<rootDir>/tests/unit/jest.setup-test-framework'],
  globalSetup: '<rootDir>/tests/unit/jest.global-setup',
  maxWorkers: process.env.JEST_MAX_WORKERS ?? '10%',
};
