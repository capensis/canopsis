import { TIME_UNITS } from '@/constants';

import { durationWithEnabledToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} DataStorageJunitConfig
 * @property {DurationWithEnabled} delete_after
 */

/**
 * @typedef {Object} DataStorageRemediationConfig
 * @property {DurationWithEnabled} delete_after
 * @property {DurationWithEnabled} delete_stats_after
 * @property {DurationWithEnabled} delete_mod_stats_after
 */

/**
 * @typedef {Object} DataStorageAlarmConfig
 * @property {DurationWithEnabled} archive_after
 * @property {DurationWithEnabled} delete_after
 */

/**
 * @typedef {Object} DataStorageEntityConfig
 * @property {boolean} archive
 * @property {boolean} archive_dependencies
 */

/**
 * @typedef {Object} DataStoragePbehaviorConfig
 * @property {DurationWithEnabled} delete_after
 */

/**
 * @typedef {Object} DataStorageHealthCheckConfig
 * @property {DurationWithEnabled} delete_after
 */

/**
 * @typedef {Object} DataStorageWebhookConfig
 * @property {DurationWithEnabled} delete_after
 * @property {boolean} log_credentials
 */

/**
 * @typedef {Object} DataStorageConfig
 * @property {DataStorageJunitConfig} junit
 * @property {DataStorageRemediationConfig} remediation
 * @property {DataStorageAlarmConfig} alarm
 * @property {DataStorageEntityConfig} [entity]
 * @property {DataStoragePbehaviorConfig} pbehavior
 * @property {DataStorageHealthCheckConfig} health_check
 * @property {DataStorageWebhookConfig} webhook
 */

/**
 * @typedef {Object} HistoryWithCount
 * @property {number} archived
 * @property {number} deleted
 * @property {number} time
 */

/**
 * @typedef {Object} DataStorageHistory
 * @property {number} junit
 * @property {number} remediation
 * @property {HistoryWithCount} alarm
 * @property {HistoryWithCount} entity
 * @property {number} health_check
 */

/**
 * @typedef {Object} DataStorage
 * @property {DataStorageConfig} config
 * @property {DataStorageHistory} history
 */

/**
 * @typedef {Object} DataStorageRequest
 * @property {DataStorageJunitConfig} junit
 */

/**
 * Convert data storage junit config to junit form object
 *
 * @param {DataStorageJunitConfig} [junitConfig = {}]
 * @return {DataStorageJunitConfig}
 */
export const dataStorageJunitSettingsToForm = (junitConfig = {}) => ({
  delete_after: junitConfig.delete_after
    ? durationWithEnabledToForm(junitConfig.delete_after)
    : { value: 1, unit: TIME_UNITS.day, enabled: false },
});

/**
 * Convert data storage remediation config to remediation form object
 *
 * @param {DataStorageRemediationConfig} remediationConfig
 * @return {DataStorageRemediationConfig}
 */
export const dataStorageRemediationSettingsToForm = (remediationConfig = {}) => ({
  delete_after: remediationConfig.delete_after
    ? durationWithEnabledToForm(remediationConfig.delete_after)
    : { value: 2, unit: TIME_UNITS.day, enabled: false },
  delete_stats_after: remediationConfig.delete_stats_after
    ? durationWithEnabledToForm(remediationConfig.delete_stats_after)
    : { value: 2, unit: TIME_UNITS.day, enabled: false },
  delete_mod_stats_after: remediationConfig.delete_mod_stats_after
    ? durationWithEnabledToForm(remediationConfig.delete_mod_stats_after)
    : { value: 2, unit: TIME_UNITS.day, enabled: false },
});

/**
 * Convert data storage alarm config to alarm form object
 *
 * @param {DataStorageAlarmConfig} alarmConfig
 * @return {DataStorageAlarmConfig}
 */
export const dataStorageAlarmSettingsToForm = (alarmConfig = {}) => ({
  archive_after: alarmConfig.archive_after
    ? durationWithEnabledToForm(alarmConfig.archive_after)
    : { value: 1, unit: TIME_UNITS.year, enabled: false },
  delete_after: alarmConfig.delete_after
    ? durationWithEnabledToForm(alarmConfig.delete_after)
    : { value: 2, unit: TIME_UNITS.year, enabled: false },
});

/**
 * Convert data storage pbehavior config to pbehavior form object
 *
 * @param {DataStoragePbehaviorConfig} pbehaviorConfig
 * @return {DataStoragePbehaviorConfig}
 */
export const dataStoragePbehaviorSettingsToForm = (pbehaviorConfig = {}) => ({
  delete_after: pbehaviorConfig.delete_after
    ? durationWithEnabledToForm(pbehaviorConfig.delete_after)
    : { value: 1, unit: TIME_UNITS.year, enabled: false },
});

/**
 * Convert data storage entity config to entity form object
 *
 * @param {DataStorageEntityConfig} entityConfig
 * @return {DataStorageEntityConfig}
 */
export const dataStorageEntitySettingsToForm = (entityConfig = {}) => ({
  archive: entityConfig.archive || false,
  archive_dependencies: entityConfig.archive_dependencies || false,
});

/**
 * Convert data storage health check config to health check form object
 *
 * @param {DataStorageHealthCheckConfig} healthCheckConfig
 * @return {DataStorageHealthCheckConfig}
 */
export const dataStorageHealthCheckSettingsToForm = (healthCheckConfig = {}) => ({
  delete_after: healthCheckConfig.delete_after
    ? durationWithEnabledToForm(healthCheckConfig.delete_after)
    : { value: 6, unit: TIME_UNITS.month, enabled: false },
});

/**
 * Convert data storage health check config to health check form object
 *
 * @param {DataStorageWebhookConfig} webhook
 * @return {DataStorageWebhookConfig}
 */
export const dataStorageWebhookSettingsToForm = (webhook = {}) => ({
  delete_after: webhook.delete_after
    ? durationWithEnabledToForm(webhook.delete_after)
    : { value: 60, unit: TIME_UNITS.day, enabled: false },
  log_credentials: webhook.log_credentials ?? false,
});

/**
 * Convert data storage object to data storage form
 *
 * @param {DataStorageConfig} dataStorage
 * @return {DataStorageConfig}
 */
export const dataStorageSettingsToForm = (dataStorage = {}) => ({
  junit: dataStorageJunitSettingsToForm(dataStorage.junit),
  remediation: dataStorageRemediationSettingsToForm(dataStorage.remediation),
  alarm: dataStorageAlarmSettingsToForm(dataStorage.alarm),
  entity: dataStorageEntitySettingsToForm(dataStorage.entity),
  pbehavior: dataStoragePbehaviorSettingsToForm(dataStorage.pbehavior),
  health_check: dataStorageHealthCheckSettingsToForm(dataStorage.health_check),
  webhook: dataStorageWebhookSettingsToForm(dataStorage.webhook),
});
