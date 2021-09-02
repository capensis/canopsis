import { durationToForm, formToDuration } from '../date/duration';

import { enabledToForm } from './shared/common';

/**
 * @typedef {Object} InstructionNotificationsSettings
 * @property {boolean} rate
 * @property {Duration} rate_frequency
 */

/**
 * @typedef {InstructionNotificationsSettings} InstructionNotificationsSettingsForm
 * @property {DurationForm} rate_frequency
 */

/**
 * @typedef {Object} NotificationsSettings
 * @property {InstructionNotificationsSettings} instruction
 */

/**
 * @typedef {Object} NotificationsSettingsForm
 * @property {InstructionNotificationsSettingsForm} instruction
 */

/**
 * Convert notification instruction settings to form object
 *
 * @param {InstructionNotificationsSettings} [instructionNotificationsSettings = {}]
 * @returns {InstructionNotificationsSettingsForm}
 */
export const instructionNotificationsSettingsToForm = (instructionNotificationsSettings = {}) => ({
  rate: enabledToForm(instructionNotificationsSettings.rate),
  rate_frequency: durationToForm(instructionNotificationsSettings.rate_frequency),
});

/**
 * Convert notification to form object
 *
 * @param {NotificationsSettings} [notificationsSettings = {}]
 * @returns {NotificationsSettingsForm}
 */
export const notificationsSettingsToForm = (notificationsSettings = {}) => ({
  instruction: instructionNotificationsSettingsToForm(notificationsSettings.instruction),
});

/**
 * Convert form to notification settings
 *
 * @param {InstructionNotificationsSettingsForm} instructionForm
 * @returns {InstructionNotificationsSettings}
 */
export const formToInstructionsNotificationSettings = instructionForm => ({
  ...instructionForm,
  rate_frequency: formToDuration(instructionForm.rate_frequency),
});

/**
 * Convert form to notification settings
 *
 * @param {NotificationsSettingsForm} form
 * @returns {NotificationsSettings}
 */
export const formToNotificationsSettings = form => ({
  instruction: formToInstructionsNotificationSettings(form.instruction),
});
