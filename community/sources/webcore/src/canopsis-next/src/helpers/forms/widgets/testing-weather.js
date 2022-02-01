import { DEFAULT_PERIODIC_REFRESH, TEST_CASE_FILE_MASK, WIDGET_TYPES } from '@/constants';

import { addKeyInEntities } from '@/helpers/entities';
import { widgetToForm } from '@/helpers/forms/widgets/common';
import { durationWithEnabledToForm } from '@/helpers/date/duration';

/**
 * @typedef {string} Storage
 */

/**
 * @typedef {Object} StorageForm
 * @property {Storage} directory
 * @property {string} key
 */

/**
 * @typedef {Object} TestingWeatherWidgetParameters
 * @property {DurationWithEnabled} periodic_refresh
 * @property {string} directory
 * @property {string} screenshot_filemask
 * @property {string} video_filemask
 * @property {boolean} is_api
 * @property {Storage[]} screenshot_directories
 * @property {Storage[]} video_directories
 * @property {string} report_fileregexp
 */

/**
 * @typedef {Widget} TestingWeatherWidget
 * @property {TestingWeatherWidgetParameters} parameters
 */

/**
 * @typedef {TestingWeatherWidgetParameters} TestingWeatherWidgetParametersForm
 * @property {StorageForm[]} screenshot_directories
 * @property {StorageForm[]} video_directories
 */

/**
 * @typedef {Widget} TestingWeatherWidgetForm
 * @property {TestingWeatherWidgetParametersForm} parameters
 */

/**
 * Convert storages array to form array
 *
 * @param {Storage[]} storages
 * @return {StorageForm[]}
 */
const storagesToFormStorages = (storages = []) => addKeyInEntities(storages.map(directory => ({ directory })));

/**
 * Convert storages array to form array
 *
 * @param {TestingWeatherWidgetParameters} parameters
 * @return {TestingWeatherWidgetParametersForm}
 */
const testingWeatherWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  directory: parameters.directory || '',
  is_api: !!parameters.is_api,
  screenshot_directories: storagesToFormStorages(parameters.screenshot_directories),
  video_directories: storagesToFormStorages(parameters.video_directories),
  screenshot_filemask: parameters.screenshot_filemask || TEST_CASE_FILE_MASK,
  video_filemask: parameters.video_filemask || TEST_CASE_FILE_MASK,
  report_fileregexp: parameters.report_fileregexp,
});

/**
 * Convert testing weather widget to form object
 *
 * @param {TestingWeatherWidget} [testingWeatherWidget = {}]
 * @returns {TestingWeatherWidgetForm}
 */
export const testingWeatherWidgetToForm = (testingWeatherWidget = {}) => {
  const widget = widgetToForm(testingWeatherWidget);

  return {
    ...widget,
    type: WIDGET_TYPES.testingWeather,
    parameters: testingWeatherWidgetParametersToForm(testingWeatherWidget.parameters),
  };
};

/**
 * Convert storages array to form array
 *
 * @param {StorageForm[]} storages
 * @return {Storage[]}
 */
const formStoragesToStorages = (storages = []) => storages.map(({ directory }) => directory);

/**
 * Convert role form to role object
 *
 * @param {TestingWeatherWidgetForm} form
 * @returns {TestingWeatherWidget}
 */
export const formToTestingWeatherWidget = (form) => {
  const { parameters } = form;

  return {
    ...form,
    parameters: {
      ...parameters,
      screenshot_directories: formStoragesToStorages(parameters.screenshot_directories),
      video_directories: formStoragesToStorages(parameters.video_directories),
    },
  };
};
