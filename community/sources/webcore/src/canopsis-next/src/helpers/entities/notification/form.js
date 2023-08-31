import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} InstructionNotificationsSettings
 * @property {boolean} rate
 * @property {Duration} rate_frequency
 */

/**
 * @typedef {Object} NotificationsSettings
 * @property {InstructionNotificationsSettings} instruction
 */

/**
 * Convert notification instruction settings to form object
 *
 * @param {InstructionNotificationsSettings} [instructionNotificationsSettings = {}]
 * @returns {InstructionNotificationsSettings}
 */
export const instructionNotificationsSettingsToForm = (instructionNotificationsSettings = {}) => ({
  rate: instructionNotificationsSettings.rate ?? true,
  rate_frequency: durationToForm(instructionNotificationsSettings.rate_frequency),
});

/**
 * Convert notification to form object
 *
 * @param {NotificationsSettings} [notificationsSettings = {}]
 * @returns {NotificationsSettings}
 */
export const notificationsSettingsToForm = (notificationsSettings = {}) => ({
  instruction: instructionNotificationsSettingsToForm(notificationsSettings.instruction),
});
