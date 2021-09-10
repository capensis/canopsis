const { spawn } = require('child_process');
const { default: PQueue } = require('p-queue');

const queue = new PQueue({ concurrency: 1 });

const requestStartType = 'startRequest';
const requestFinishType = 'requestFinish';
const initQueueType = 'initQueue';

function nightwatchRunWithQueue(colors, done) {
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

    let resolveFunction = () => {};

    const childQueueHandler = (message) => {
      const data = JSON.parse(message);
      this.emit('result', data);

      if (data.type === initQueueType) {
        queue.add(() => new Promise((resolve) => {
          this.child.send(JSON.stringify({ type: requestStartType }));
          resolveFunction = resolve;
        }));
      }

      if (data.type === requestFinishType) {
        resolveFunction();
      }
    };

    this.child.on('message', childQueueHandler);

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
}

/**
 * Function return a promise, that resolve, when the turn of the process comes.
 *
 * ... performed anyway
 * await onNextQueueFunction();
 * ... performed when the turn of the process comes
 *
 * @returns {Promise}
 */
const onNextQueueFunction = () => new Promise((resolve) => {
  const eventHandler = (message) => {
    const data = JSON.parse(message);
    if (data.type === requestStartType) {
      resolve(eventHandler);
    }
  };

  if (!process.send) {
    resolve();
  } else {
    process.send(JSON.stringify({ type: initQueueType }));
    process.on('message', eventHandler);
  }
}).then(handler => handler && process.removeListener('message', handler));

/**
 * Call and await callback, after resolve send a message to parent process
 *
 * await queueFunction(async () => {
 * ... The code that should be run alternately in all processes.
 * });
 *
 * @param request
 * @returns {Promise<*|Promise<unknown>>}
 */
const queueFunction = async (request) => {
  const response = await request();

  if (process.send) {
    process.send(JSON.stringify({ type: requestFinishType }));
  }

  return response;
};

module.exports.nightwatchRunWithQueue = nightwatchRunWithQueue;
module.exports.requestStartType = requestStartType;
module.exports.requestFinishType = requestFinishType;
module.exports.initQueueType = initQueueType;
module.exports.onNextQueueFunction = onNextQueueFunction;
module.exports.queueFunction = queueFunction;
