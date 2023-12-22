#!/usr/bin/env node

require('./scenarios/alarms');

const { runBenchmarks, compareMetric } = require('./utils/runner');

// eslint-disable-next-line import/no-extraneous-dependencies
require('yargs')
  .scriptName('benchmarks')
  .usage('$0 <cmd> [args]')
  .command(
    'run',
    'run measure metrics',
    (yargs) => {
      yargs.positional('url', {
        type: 'string',
        default: 'https://localhost:8080',
        describe: 'Application url',
      });
      yargs.positional('viewId', {
        type: 'string',
        describe: 'View id',
      });
      yargs.positional('tabId', {
        type: 'string',
        describe: 'View tab id',
      });
      yargs.positional('jsonName', {
        type: 'string',
        describe: 'Name of json report',
      });
    },
    runBenchmarks,
  )
  .command(
    'compare',
    'compare metrics and render images from metrics and export pdf',
    (yargs) => {
      yargs.positional('target', {
        type: 'string',
        describe: 'Target metrics',
      });
      yargs.positional('source', {
        type: 'string',
        describe: 'Source metrics',
      });
    },
    compareMetric,
  )
  .help()
  .argv;
