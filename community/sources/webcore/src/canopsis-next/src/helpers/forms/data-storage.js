import { durationWithEnabledToForm, formToDurationWithEnabled } from '@/helpers/date/duration';

/**
 * @typedef {Object} DataStorageJunitConfig
 * @property {DurationWithEnabled} delete_after
 */

/**
 * @typedef {Object} DataStorageConfig
 * @property {DataStorageJunitConfig} junit
 */

/**
 * @typedef {Object} DataStorageHistory
 */

/**
 * @typedef {Object} DataStorage
 * @property {DataStorageConfig} config
 * @property {DataStorageHistory} history
 */

/**
 * @typedef {Object} DataStorageJunitConfigForm
 * @property {DurationWithEnabledForm} delete_after
 */

/**
 * @typedef {Object} DataStorageConfigForm
 * @property {DataStorageJunitConfigForm} junit
 */

/**
 * @typedef {Object} DataStorageRequest
 * @property {DataStorageJunitConfigForm} junit
 */

/**
 * @param {DataStorageJunitConfig} junitConfig
 * @return {DataStorageJunitConfigForm}
 */
export const dataStorageJunitSettingsToForm = (junitConfig = {}) => ({
  delete_after: durationWithEnabledToForm(junitConfig.delete_after || {}),
});

/**
 * Convert data storage object to data storage form
 *
 * @param {DataStorageConfig} dataStorage
 * @return {DataStorageConfigForm}
 */
export const dataStorageSettingsToForm = (dataStorage = {}) => ({
  junit: dataStorageJunitSettingsToForm(dataStorage.junit),
});

/**
 * Convert junit data storage form to junit data storage object
 *
 * @param {DataStorageJunitConfigForm} form
 * @return {DataStorageJunitConfig}
 */
export const formJunitToDataStorageSettings = (form = {}) => ({
  delete_after: formToDurationWithEnabled(form.delete_after),
});

/**
 * Convert data storage form to data storage object
 *
 * @param {DataStorageConfigForm} form
 * @return {DataStorageConfig}
 */
export const formToDataStorageSettings = (form = {}) => ({
  junit: formJunitToDataStorageSettings(form.junit),
});
