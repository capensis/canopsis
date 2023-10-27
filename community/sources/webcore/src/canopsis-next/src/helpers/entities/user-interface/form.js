import { DEFAULT_APP_TITLE, DEFAULT_LOCALE, POPUP_AUTO_CLOSE_DELAY } from '@/config';
import { TIME_UNITS } from '@/constants';

import { durationToForm } from '@/helpers/date/duration';

/**
 * @typedef {Object} PopupTimeout
 * @property {Duration} error
 * @property {Duration} info
 */

/**
 * @typedef {Object} UserInterface
 * @property {string} app_title
 * @property {string} language
 * @property {string} footer
 * @property {string} login_page_description
 * @property {string} timezone
 * @property {boolean} allow_change_severity_to_info
 * @property {boolean} show_header_on_kiosk_mode
 * @property {boolean} required_instruction_approve
 * @property {PopupTimeout} popup_timeout
 * @property {string} [logo]
 * @property {number} [max_matched_items]
 * @property {number} [check_count_request_timeout]
 */

/**
 * @typedef {UserInterface} UserInterfaceForm
 * @property {File} [logo]
 */

/**
 * @typedef {UserInterface} UserInterfaceRequest
 * @property {string} edition
 * @property {string} logo
 * @property {string} stack
 * @property {string} version
 */

/**
 * Convert user interface popupTimeout to form
 *
 * @param {PopupTimeout} [popupTimeout = {}]
 * @return {PopupTimeout}
 */
const userInterfacePopupTimeoutToForm = (popupTimeout = {}) => ({
  info: popupTimeout.info
    ? durationToForm(popupTimeout.info)
    : { value: POPUP_AUTO_CLOSE_DELAY, unit: TIME_UNITS.second },
  error: popupTimeout.error
    ? durationToForm(popupTimeout.error)
    : { value: POPUP_AUTO_CLOSE_DELAY, unit: TIME_UNITS.second },
});

/**
 * Convert userInterface object to form
 *
 * @param {UserInterface | {}} [userInterface = {}]
 * @returns {UserInterfaceForm}
 */
export const userInterfaceToForm = (userInterface = {}) => ({
  app_title: userInterface.app_title ?? DEFAULT_APP_TITLE,
  language: userInterface.language ?? DEFAULT_LOCALE,
  footer: userInterface.footer ?? '',
  login_page_description: userInterface.login_page_description ?? '',
  allow_change_severity_to_info: userInterface.allow_change_severity_to_info ?? false,
  show_header_on_kiosk_mode: userInterface.show_header_on_kiosk_mode ?? false,
  required_instruction_approve: userInterface.required_instruction_approve ?? false,
  timezone: userInterface.timezone ?? '',
  max_matched_items: userInterface.max_matched_items ?? '',
  check_count_request_timeout: userInterface.check_count_request_timeout ?? '',
  popup_timeout: userInterfacePopupTimeoutToForm(userInterface.popup_timeout),
});
