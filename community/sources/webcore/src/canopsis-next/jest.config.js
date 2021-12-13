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
  },
  transformIgnorePatterns: [
    '<rootDir>/node_modules/(?!vue-tour/.*)',
  ],
  moduleNameMapper: {
    '^.+\\.styl(us)?$': '<rootDir>/tests/unit/mocks/styleMock.js',
    '^@unit/(.*)$': '<rootDir>/tests/unit/$1',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  snapshotSerializers: [
    'jest-serializer-vue',
  ],
  setupFiles: ['jest-localstorage-mock'],
  setupFilesAfterEnv: ['<rootDir>/tests/unit/jest.setup-test-framework'],
};
