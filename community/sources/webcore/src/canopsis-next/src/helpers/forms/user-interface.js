import { cloneDeep, pick, isNumber } from 'lodash';

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
 * @property {string} language
 * @property {string} footer
 * @property {string} login_page_description
 * @property {string} timezone
 * @property {boolean} allow_change_severity_to_info
 * @property {PopupTimeout} popup_timeout
 * @property {number} [max_matched_items]
 * @property {number} [check_count_request_timeout]
 */
/**
 * @typedef {UserInterface} UserInterfaceForm
 * @property {number | string} [max_matched_items]
 * @property {number | string} [check_count_request_timeout]
 */

/**
 * @typedef {UserInterface} UserInterfaceRequest
 * @property {string} edition
 * @property {string | null} logo
 * @property {string} stack
 * @property {string} version
 */

/**
 * Convert userInterface object to form
 *
 * @param {UserInterface | {}} [userInterface = {}]
 * @returns {UserInterfaceForm}
 */
export const userInterfaceToForm = (userInterface = {}) => ({
  app_title: userInterface.app_title || DEFAULT_APP_TITLE,
  language: userInterface.language || DEFAULT_LOCALE,
  footer: userInterface.footer || '',
  login_page_description: userInterface.login_page_description || '',
  allow_change_severity_to_info: userInterface.allow_change_severity_to_info || false,
  timezone: userInterface.timezone || '',
  max_matched_items: userInterface.max_matched_items || '',
  check_count_request_timeout: userInterface.check_count_request_timeout || '',
  popup_timeout: userInterface.popup_timeout ? cloneDeep(userInterface.popup_timeout) : {},
});

/**
 * Convert form to userInterface object
 *
 * @param {UserInterfaceForm | {}} [form = {}]
 * @returns {UserInterface}
 */
export const formToUserInterface = (form = {}) => {
  const userInterface = pick(form, [
    'app_title',
    'language',
    'footer',
    'login_page_description',
    'allow_change_severity_to_info',
    'timezone',
    'popup_timeout',
  ]);

  if (isNumber(form.max_matched_items)) {
    userInterface.max_matched_items = form.max_matched_items;
  }

  if (isNumber(form.check_count_request_timeout)) {
    userInterface.check_count_request_timeout = form.check_count_request_timeout;
  }

  return userInterface;
};
