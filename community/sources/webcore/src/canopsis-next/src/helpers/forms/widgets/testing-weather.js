import { DEFAULT_PERIODIC_REFRESH, TEST_CASE_FILE_MASK } from '@/constants';

import { addKeyInEntities } from '@/helpers/entities';
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
 * @typedef {TestingWeatherWidgetParameters} TestingWeatherWidgetParametersForm
 * @property {StorageForm[]} screenshot_directories
 * @property {StorageForm[]} video_directories
 */

/**
 * Convert storages array to form array
 *
 * @param {Storage[]} storages
 * @return {StorageForm[]}
 */
const storagesToForm = (storages = []) => addKeyInEntities(storages.map(directory => ({ directory })));

/**
 * Convert testing weather widget parameters to form
 *
 * @param {TestingWeatherWidgetParameters} parameters
 * @return {TestingWeatherWidgetParametersForm}
 */
export const testingWeatherWidgetParametersToForm = (parameters = {}) => ({
  periodic_refresh: durationWithEnabledToForm(parameters.periodic_refresh ?? DEFAULT_PERIODIC_REFRESH),
  directory: parameters.directory ?? '',
  is_api: parameters.is_api ?? false,
  screenshot_directories: storagesToForm(parameters.screenshot_directories),
  video_directories: storagesToForm(parameters.video_directories),
  screenshot_filemask: parameters.screenshot_filemask ?? TEST_CASE_FILE_MASK,
  video_filemask: parameters.video_filemask ?? TEST_CASE_FILE_MASK,
  report_fileregexp: parameters.report_fileregexp ?? '',
});

/**
 * Convert form array storages array
 *
 * @param {StorageForm[]} storages
 * @return {Storage[]}
 */
const formToStorages = (storages = []) => storages.map(({ directory }) => directory);

/**
 * Convert form to testing weather widget parameters
 *
 * @param {TestingWeatherWidgetParametersForm} form
 * @return {TestingWeatherWidgetParameters}
 */
export const formToTestingWeatherWidgetParameters = form => ({
  ...form,

  screenshot_directories: formToStorages(form.screenshot_directories),
  video_directories: formToStorages(form.video_directories),
});
