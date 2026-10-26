import { omit } from 'lodash';

import {
  EXTERNAL_DATA_CONDITION_TYPES,
  EXTERNAL_DATA_ON_MISMATCH_TYPES,
  EXTERNAL_DATA_TYPES,
} from '@/constants';

import { isApiExternalDataType } from '@/helpers/entities/shared/external-data/entity';
import { uid } from '@/helpers/uid';

import {
  formToRequest,
  formToRequestAuthToken,
  requestAuthTokenToForm,
  requestToForm,
} from '../request/form';

/**
 * @typedef {'mongo' | 'api'} ExternalDataType
 */

/**
 * @typedef {'select' | 'regexp'} ExternalDataConditionType
 */

/**
 * @typedef {Object<ExternalDataConditionType, Object<string, string>>} ExternalDataCondition
 */

/**
 * @typedef {ExternalDataCondition} ExternalDataItem
 * @property {ExternalDataType} type
 * @property {string} reference
 * @property {string} [table]
 * @property {string} [sort_by]
 * @property {string} [sort]
 * @property {boolean} [optional]
 * @property {number} [on_mismatch]
 * @property {Request} [request]
 * @property {RequestAuthToken} [auth_token]
 */

/**
 * @typedef {ExternalDataItem[]} ExternalData
 */

/**
 * @typedef {Object & ObjectKey} ExternalDataConditionForm
 * @property {ExternalDataConditionType} type
 * @property {string} attribute
 * @property {string} value
 */

/**
 * @typedef {Object & ObjectKey} ExternalDataItemForm
 * @property {RequestForm} request
 * @property {RequestAuthTokenForm} auth_token
 * @property {string} reference
 * @property {ExternalDataType} type
 * @property {string} table
 * @property {string} sort_by
 * @property {string} sort
 * @property {boolean} optional
 * @property {number} [on_mismatch]
 * @property {ExternalDataConditionForm[]} conditions
 */

/**
 * @typedef {ExternalDataItemForm[]} ExternalDataForm
 */

export const externalDataItemConditionAttributeToForm = (
  conditionType = EXTERNAL_DATA_CONDITION_TYPES.select,
  attribute = '',
  value = '',
) => ({
  key: uid(),
  type: conditionType,
  attribute,
  value,
});

/**
 * Convert external data item whole condition to form
 *
 * @param {ExternalDataConditionType} conditionType
 * @param {ExternalDataCondition} condition
 * @returns {ExternalDataConditionForm[]}
 */
export const externalDataItemConditionToForm = (conditionType, condition) =>
  Object.entries(condition).map(([attribute, value]) =>
    externalDataItemConditionAttributeToForm(conditionType, attribute, value),
  );

/**
 * Convert external data item conditions to form
 *
 * @param {ExternalDataCondition} [item = {}]
 * @returns {ExternalDataConditionForm[]}
 */
export const externalDataItemConditionsToForm = (item = {}) => {
  const conditions = Object.values(EXTERNAL_DATA_CONDITION_TYPES).reduce((acc, conditionType) => {
    const condition = item[conditionType];

    if (condition) {
      acc.push(...externalDataItemConditionToForm(conditionType, condition));
    }

    return acc;
  }, []);

  if (!conditions.length) {
    conditions.push(externalDataItemConditionAttributeToForm());
  }

  return conditions;
};

/**
 * Convert external data item to form
 *
 * @param {ExternalDataItem} item
 * @param {boolean} [withOnMismatch]
 * @returns {ExternalDataItemForm}
 */
export const externalDataItemToForm = (
  item = { type: EXTERNAL_DATA_TYPES.table },
  withOnMismatch = false,
) => ({
  key: uid(),
  reference: item.reference,
  type: item.type,
  request: requestToForm(item.request, item.auth_token),
  auth_token: requestAuthTokenToForm(item.auth_token),
  sort_by: item.sort_by,
  sort: item.sort,
  optional: item.optional ?? false,
  table: item.table?._id ?? '',
  conditions: externalDataItemConditionsToForm(item),
  ...(withOnMismatch
    ? { on_mismatch: item.on_mismatch ?? EXTERNAL_DATA_ON_MISMATCH_TYPES.throwError }
    : {}),
});

/**
 * Convert external data to form
 *
 * @param {ExternalData} [externalData]
 * @param {boolean} [withOnMismatch]
 * @returns {ExternalDataForm}
 */
export const externalDataToForm = (externalData, withOnMismatch = false) =>
  externalData
    ? externalData.map(item => externalDataItemToForm(item, withOnMismatch))
    : [];

/**
 * Convert form to external data conditions
 *
 * @param {ExternalDataConditionForm[]} form
 * @returns {ExternalDataCondition}
 */
export const formToExternalDataConditions = (form = []) =>
  form.reduce((acc, { type, attribute, value }) => {
    if (!acc[type]) {
      acc[type] = {};
    }

    acc[type][attribute] = value;

    return acc;
  }, {});

/**
 * Convert form to external data
 *
 * @param {ExternalDataForm} form
 * @param {boolean} [withOnMismatch]
 * @returns {ExternalData}
 */
export const formToExternalData = (form = [], withOnMismatch = false) =>
  form.map((externalData) => {
    const { type, reference } = externalData;

    return isApiExternalDataType(type)
      ? {
          reference,
          type,
          request: formToRequest(externalData.request),
          auth_token: formToRequestAuthToken(
            externalData.auth_token,
            externalData.request.auth?.type,
          ),
        }
      : {
          ...omit(externalData, ['request', 'auth_token', 'on_mismatch']),
          ...formToExternalDataConditions(externalData.conditions),
          ...(withOnMismatch && !externalData.optional
            ? {
                on_mismatch:
                  externalData.on_mismatch ?? EXTERNAL_DATA_ON_MISMATCH_TYPES.throwError,
              }
            : {}),
        };
  });

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsObject
 * @return {FlattenErrors}
 */
export const externalDataConditionTemplateVariablesErrorsToForm = (errorsObject) => {
  const { value } = errorsObject;
  const conditionErrors = {};

  if (!value.is_valid) {
    conditionErrors.value = value.err.message;
  }

  return conditionErrors;
};

/**
 * Convert template variables errors structure to form structure
 *
 * @param {Object} errorsArray
 * @param {ExternalDataConditionForm[]} conditions
 * @return {FlattenErrors}
 */
export const externalDataConditionsTemplateVariablesErrorsToForm = (errorsArray, conditions) =>
  errorsArray.reduce((acc, errors, index) => {
    const condition = conditions[index];

    acc[condition.key] = externalDataConditionTemplateVariablesErrorsToForm(errors);

    return acc;
  }, {});
