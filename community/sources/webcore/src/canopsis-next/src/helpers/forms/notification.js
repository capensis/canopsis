import { durationToForm } from '../date/duration';

import { enabledToForm } from './shared/common';

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
  rate: enabledToForm(instructionNotificationsSettings.rate),
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
