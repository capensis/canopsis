import { DEFAULT_APP_TITLE, DEFAULT_LOCALE, POPUP_AUTO_CLOSE_DELAY } from '@/config';
import { durationToForm, formToDuration } from '@/helpers/date/duration';
import { TIME_UNITS } from '@/constants';

/**
 * @typedef {Object} PopupTimeout
 * @property {Duration} error
 * @property {Duration} info
 */

/**
 * @typedef {Object} PopupTimeoutForm
 * @property {DurationForm} error
 * @property {DurationForm} info
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
 * @property {string} [logo]
 * @property {number} [max_matched_items]
 * @property {number} [check_count_request_timeout]
 */
/**
 * @typedef {UserInterface} UserInterfaceForm
 * @property {PopupTimeoutForm} popup_timeout
 * @property {!File} [logo]
 */

/**
 * @typedef {UserInterface} UserInterfaceRequest
 * @property {string} edition
 * @property {string} logo
 * @property {string} stack
 * @property {string} version
 */

/**
 *
 *
 * @param {PopupTimeout} [popupTimeout = {}]
 * @return {PopupTimeoutForm}
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
  timezone: userInterface.timezone ?? '',
  max_matched_items: userInterface.max_matched_items ?? '',
  check_count_request_timeout: userInterface.check_count_request_timeout ?? '',
  popup_timeout: userInterfacePopupTimeoutToForm(userInterface.popup_timeout),
});

/**
 *
 *
 * @param {PopupTimeoutForm | {}} [popupTimeout = {}]
 * @return {PopupTimeout}
 */
const formPopupTimeoutToUserInterfacePopupTimeout = (popupTimeout = {}) => ({
  info: formToDuration(popupTimeout.info),
  error: formToDuration(popupTimeout.error),
});

/**
 * Convert form to userInterface object
 *
 * @param {UserInterfaceForm | {}} [form = {}]
 * @returns {UserInterface}
 */
export const formToUserInterface = (form = {}) => ({
  ...form,
  popup_timeout: formPopupTimeoutToUserInterfacePopupTimeout(form.popup_timeout),
});
