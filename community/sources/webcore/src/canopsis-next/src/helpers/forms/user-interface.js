import { cloneDeep } from 'lodash';

import { DEFAULT_APP_TITLE, DEFAULT_LOCALE } from '@/config';

// TODO: change to duration
/**
 * @typedef {Object} PopupTimeoutItem
 * @property {number} interval
 * @property {DurationUnit} unit
 */

/**
 * @typedef {Object} PopupTimeout
 * @property {PopupTimeoutItem} error
 * @property {PopupTimeoutItem} info
 */

/**
 * @typedef {Object} UserInterface
 * @property {string} app_title
 * @property {string} edition
 * @property {string} language
 * @property {string | null} logo
 * @property {string} footer
 * @property {string} login_page_description
 * @property {string} stack
 * @property {string} timezone
 * @property {string} version
 * @property {boolean} allow_change_severity_to_info
 * @property {number} max_matched_items
 * @property {number} check_count_request_timeout
 * @property {PopupTimeout} popup_timeout
 */

export const userInterfaceToForm = (userInterface = {}) => ({
  app_title: userInterface.app_title || DEFAULT_APP_TITLE,
  language: userInterface.language || DEFAULT_LOCALE,
  footer: userInterface.footer || '',
  login_page_description: userInterface.login_page_description || '',
  allow_change_severity_to_info: userInterface.allow_change_severity_to_info || false,
  timezone: userInterface.timezone || '',
  max_matched_items: userInterface.max_matched_items || 10000,
  check_count_request_timeout: userInterface.check_count_request_timeout || 30,
  popup_timeout: userInterface.popup_timeout ? cloneDeep(userInterface.popup_timeout) : {},
});
