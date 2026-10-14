const path = require('path');

const loadEnv = require('../../tools/load-env');

const localEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env.local');
const baseEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env');

loadEnv(localEnvPath);
loadEnv(baseEnvPath);

module.exports = {
  fileName: 'test-result',
  nameAfterTest: true,
  format: 'mp4',
  enabled: process.env.TEST_VIDEOS_ENABLED === 'true',
  deleteOnSuccess: true,
  path: path.resolve('tests', 'e2e', 'records'),
  resolution: '1440x900',
  fps: 15,
  input: '',
  videoCodec: 'libx264',
};
