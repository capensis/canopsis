#!/usr/bin/env node

require('./scenarios/alarms');

const { runBenchmarks, compareMetric } = require('./utils/runner');

// eslint-disable-next-line import/no-extraneous-dependencies
require('yargs')
  .scriptName('benchmarks')
  .usage('$0 <cmd> [args]')
  .command(
    /**
     * Run measure metrics for compare metrics
     *
     * @example
     *
     * yarn benchmark run --url=https://localhost:8080 --viewId=view-id --tabId=view-tab-id --jsonName=release-23.10
     */
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
    /**
     * Command for compare metrics
     *
     * @example
     * For compare all available metrics
     * yarn benchmark compare
     *
     * For compare two metrics
     * yarn benchmark compare --target=release-23.04 --source=release-23.10
     */
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
