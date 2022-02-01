// http://nightwatchjs.org/gettingstarted#settings-file

/**
 * The following code helps us for requiring es6 modules
 */
process.env.VUE_CLI_BABEL_TARGET_NODE = true;
process.env.VUE_CLI_BABEL_TRANSPILE_MODULES = true;

require('@babel/register')({
  plugins: [
    'require-context-hook',
    ['babel-plugin-module-resolver', {
      root: ['.'],
      alias: {
        '@': './src',
      },
    }],
  ],
});

require('babel-plugin-require-context-hook/register')();

const path = require('path');
const deepmerge = require('deepmerge');

const seleniumServer = require('selenium-server');
const chromedriver = require('chromedriver');
const ChildProcess = require('nightwatch/lib/runner/concurrency/child-process');

const { nightwatchRunWithQueue } = require('./helpers/nightwatch-child-process');

const loadEnv = require('../../tools/load-env'); // eslint-disable-line import/no-extraneous-dependencies
const nightWatchRecordConfig = require('./nightwatch-record.config');

const localEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env.local');
const baseEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env');
const testMode = process.env.E2E_TESTS_MODE;
const isParallelMode = testMode === 'parallel';
const isConsistentlyMode = testMode === 'consistently';

loadEnv(localEnvPath);
loadEnv(baseEnvPath);

const sel = require('./helpers/sel');

ChildProcess.prototype.run = nightwatchRunWithQueue;

const userOptions = JSON.parse(process.env.VUE_NIGHTWATCH_USER_OPTIONS || '{}');

const seleniumConfig = {
  start_process: true,
  server_path: seleniumServer.path,
  host: '127.0.0.1',
  port: 4444,
  cli_args: {
    'webdriver.chrome.driver': chromedriver.path,
  },
};

/**
 * Put sel helper method into global object
 */
global.sel = sel;
global.window = {
  location: {
    href: process.env.VUE_DEV_SERVER_URL,
  },
};

module.exports = deepmerge({
  src_folders: [path.resolve('tests', 'e2e', 'specs')],
  output_folder: path.resolve('tests', 'e2e', 'reports'),
  custom_assertions_path: [
    path.resolve('node_modules', 'nightwatch-xhr', 'es5', 'assertions'),
    path.resolve('tests', 'e2e', 'custom-assertions'),
  ],
  custom_commands_path: [
    path.resolve('node_modules', 'nightwatch-xhr', 'es5', 'commands'),
    path.resolve('tests', 'e2e', 'custom-commands'),
  ],
  page_objects_path: [path.resolve('tests', 'e2e', 'page-objects')],
  globals_path: path.resolve('tests', 'e2e', 'globals.js'),

  selenium: seleniumConfig,

  test_workers: {
    enabled: isParallelMode,
    workers: Number(process.env.TEST_WORKERS_COUNT),
  },
  live_output: process.env.TEST_WORKERS_LIVE_OUTPUT_ENABLED === 'true',

  test_settings: {
    default: {
      selenium_host: seleniumConfig.host,
      selenium_port: seleniumConfig.port,
      silent: true,

      videos: nightWatchRecordConfig,

      exclude: isParallelMode && ['**/*.consistently.js'],
      filter: isConsistentlyMode && '**/*.consistently.js',
    },

    chrome: {
      desiredCapabilities: {
        browserName: 'chrome',
        javascriptEnabled: true,
        acceptSslCerts: true,
        chromeOptions: {
          args: ['--no-sandbox', '--ignore-certificate-errors'],
        },
      },
    },
  },
}, userOptions);
