/** http://nightwatchjs.org/gettingstarted#settings-file */

const parh = require('path');
const deepmerge = require('deepmerge');

/* eslint-disable import/no-extraneous-dependencies */
const seleniumServer = require('selenium-server');
const chromeDriver = require('chromedriver');
/* eslint-enable import/no-extraneous-dependencies */

const sel = require('./helpers/sel');

const userOptions = JSON.parse(process.env.VUE_NIGHTWATCH_USER_OPTIONS || '{}');

const seleniumConfig = {
  start_process: true,
  server_path: seleniumServer.path,
  host: '127.0.0.1',
  port: 4444,
  cli_args: {
    'webdriver.chrome.driver': chromeDriver.path,
  },
};

/**
 * Put sel helper method into global object
 */
global.sel = sel;

module.exports = deepmerge({
  src_folders: [parh.resolve('tests', 'e2e', 'specs')],
  output_folder: parh.resolve('tests', 'e2e', 'reports'),
  custom_assertions_path: [parh.resolve('tests', 'e2e', 'custom-assertions')],
  custom_commands_path: [parh.resolve('tests', 'e2e', 'custom-commands')],
  page_objects_path: [parh.resolve('tests', 'e2e', 'page-objects')],
  globals_path: parh.resolve('tests', 'e2e', 'globals.js'),

  selenium: seleniumConfig,

  test_settings: {
    default: {
      selenium_host: seleniumConfig.host,
      selenium_port: seleniumConfig.port,
      silent: true,

      videos: {
        fileName: 'test-result', // Required field
        nameAfterTest: true,
        format: 'mp4',
        enabled: true,
        deleteOnSuccess: true,
        path: parh.resolve('tests', 'e2e', 'records'),
        resolution: '1440x900',
        fps: 15,
        input: '',
        videoCodec: 'libx264',
      },
    },

    chrome: {
      desiredCapabilities: {
        browserName: 'chrome',
        javascriptEnabled: true,
        acceptSslCerts: true,
      },
    },
  },
}, userOptions);
