import { setField } from '@/helpers/immutable';
import { durationToForm } from '@/helpers/date/duration';
import { requestToForm, formToRequest } from '@/helpers/entities/shared/request/form';

/**
 * @typedef {Object} ExternalAuthToken
 * @property {string} [name] - The name of the external auth token
 * @property {string} [description] - The description of the external auth token
 * @property {string} [template] - The template for the external auth token
 * @property {Duration} [expiration_duration] - The expiration duration
 * @property {Request} [request] - The request configuration object
 */

/**
 * @typedef {Object} ExternalAuthTokenForm
 * @property {string} name - The name field for the form
 * @property {string} description - The description field for the form
 * @property {string} template - The template field for the form
 * @property {Duration} expiration_duration - The formatted expiration duration for the form
 * @property {RequestForm} request - The formatted request configuration for the form
 */

/**
 * Converts an external auth token entity to form data
 *
 * @param {ExternalAuthToken} [externalAuthToken={}] - The external auth token entity
 * @returns {ExternalAuthTokenForm} The form data object
 */
export const externalAuthTokenToForm = (externalAuthToken = {}) => ({
  name: externalAuthToken.name ?? '',
  description: externalAuthToken.description ?? '',
  template: externalAuthToken.template ?? '',
  expiration_duration: durationToForm(externalAuthToken.expiration_duration),
  request: requestToForm(externalAuthToken.request ?? {}),
});

/**
 * Converts form data to an external auth token entity
 *
 * @param {ExternalAuthTokenForm} [form={}] - The form data object
 * @returns {ExternalAuthToken} The external auth token entity
 */
export const formToExternalAuthToken = (form = {}) => ({
  ...form,

  request: formToRequest(setField(form.request ?? {}, 'auth.enabled', true)),
});
