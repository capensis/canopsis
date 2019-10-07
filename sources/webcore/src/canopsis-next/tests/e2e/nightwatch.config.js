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
const { spawn } = require('child_process');
const PQueue = require('p-queue');
const axios = require('axios');

/* eslint-disable import/no-extraneous-dependencies */
const seleniumServer = require('selenium-server');
const ChildProcess = require('nightwatch/lib/runner/cli/child-process');
/* eslint-enable import/no-extraneous-dependencies */

const loadEnv = require('../../tools/load-env'); // eslint-disable-line import/no-extraneous-dependencies

const localEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env.local');
const baseEnvPath = path.resolve(process.cwd(), 'tests', 'e2e', '.env');

loadEnv(localEnvPath);
loadEnv(baseEnvPath);

const sel = require('./helpers/sel');

const queue = new PQueue({ concurrency: 1 });

ChildProcess.prototype.run = function run(colors, done) {
  this.availColors = colors;

  const cliArgs = this.getArgs();
  const env = {};

  Object.keys(process.env).forEach((key) => {
    env[key] = process.env[key];
  });

  setTimeout(() => {
    /* eslint-disable no-underscore-dangle */
    env.__NIGHTWATCH_PARALLEL_MODE = 1;
    env.__NIGHTWATCH_ENV = this.environment;
    env.__NIGHTWATCH_ENV_KEY = this.itemKey;
    env.__NIGHTWATCH_ENV_LABEL = this.env_itemKey;
    /* eslint-enable no-underscore-dangle */

    this.child = spawn(process.execPath, cliArgs, {
      env,

      cwd: process.cwd(),
      encoding: 'utf8',
      stdio: [null, null, null, 'ipc'],
    });

    this.child.on('message', (data) => {
      const result = JSON.parse(data);

      this.emit('result', result);

      if (result.type === 'request') {
        queue.push(() => axios(result.data))
          .then(({ data: responseData }) => this.child.send(JSON.stringify({ type: 'response', data: responseData })))
          .catch(err => console.error(err));
      }
    });

    this.processRunning = true;

    if (this.settings.output) {
      // eslint-disable-next-line no-console
      console.log(`Started child process for:${this.env_label}`);
    }

    this.child.stdout.on('data', (data) => {
      this.writeToStdout(data);
    });

    this.child.stderr.on('data', (data) => {
      this.writeToStdout(data);
    });

    this.child.on('exit', (code) => {
      if (this.settings.output) {
        // eslint-disable-next-line no-console
        console.log(`\n  >>${this.env_label}finished. \n`);
      }

      if (code) {
        this.globalExitCode = 2;
      }
      this.processRunning = false;
      done(this.env_output, code);
    });
  }, this.index * this.startDelay);
};

const userOptions = JSON.parse(process.env.VUE_NIGHTWATCH_USER_OPTIONS || '{}');

const seleniumConfig = {
  start_process: true,
  server_path: seleniumServer.path,
  host: '127.0.0.1',
  port: 4444,
  cli_args: {
    'webdriver.chrome.driver': process.env.CHROME_DRIVER_PATH,
  },
};

/**
 * Put sel helper method into global object
 */
global.sel = sel;

/* process.on('message', (...args) => {
  console.log(args);
}); */

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

  test_settings: {
    default: {
      selenium_host: seleniumConfig.host,
      selenium_port: seleniumConfig.port,
      silent: true,

      videos: {
        fileName: 'test-result', // Required field
        nameAfterTest: true,
        format: 'mp4',
        enabled: false,
        deleteOnSuccess: true,
        path: path.resolve('tests', 'e2e', 'records'),
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
        chromeOptions: {
          args: ['--no-sandbox'],
        },
      },
    },
  },
  test_workers: {
    enabled: true,
    workers: 2,
  },
  live_output: true,
}, userOptions);
