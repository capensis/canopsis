module.exports = {
  moduleFileExtensions: [
    'js',
    'jsx',
    'json',
    'vue',
  ],
  transform: {
    '^.+\\.vue$': 'vue-jest',
    '^.+\\.jsx?$': '<rootDir>/tests/unit/jest.transform',
  },
  moduleNameMapper: {
    '^@unit/(.*)$': '<rootDir>/tests/unit/$1',
    '^@/(.*)$': '<rootDir>/src/$1',
  },
  snapshotSerializers: [
    'jest-serializer-vue',
  ],
  setupFiles: ['jest-localstorage-mock'],
  setupFilesAfterEnv: ['<rootDir>/tests/unit/jest.setup-test-framework'],
};
