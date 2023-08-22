import { pick } from 'lodash';

import { EXTERNAL_DATA_CONDITION_TYPES, EXTERNAL_DATA_TYPES } from '@/constants';

import { uid } from '@/helpers/uid';

import { formToRequest, requestToForm } from '../request/form';

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
 * @property {string} [sort_by]
 * @property {string} [sort]
 * @property {string} [collection]
 * @property {Request} [request]
 */

/**
 * @typedef {Object<string, ExternalDataItem>} ExternalData
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
 * @property {string} reference
 * @property {ExternalDataType} type
 * @property {string} collection
 * @property {string} sort_by
 * @property {string} sort
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
export const externalDataItemConditionToForm = (conditionType, condition) => (
  Object.entries(condition)
    .map(([attribute, value]) => externalDataItemConditionAttributeToForm(conditionType, attribute, value))
);

/**
 * Convert external data item conditions to form
 *
 * @param {ExternalDataCondition} [item = {}]
 * @returns {ExternalDataConditionForm[]}
 */
export const externalDataItemConditionsToForm = (item = {}) => {
  const conditions = Object.values(EXTERNAL_DATA_CONDITION_TYPES)
    .reduce((acc, conditionType) => {
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
 * @param {string} reference
 * @param {ExternalDataItem} item
 * @returns {ExternalDataItemForm}
 */
export const externalDataItemToForm = (reference = '', item = { type: EXTERNAL_DATA_TYPES.mongo }) => ({
  key: uid(),
  reference,
  type: item.type,
  request: requestToForm(item.request),
  sort_by: item.sort_by,
  sort: item.sort,
  collection: item.collection ?? '',
  conditions: externalDataItemConditionsToForm(item),
});

/**
 * Convert external data to form
 *
 * @param {ExternalData} [externalData]
 * @returns {ExternalDataForm}
 */
export const externalDataToForm = externalData => (
  externalData
    ? Object.entries(externalData).map(([reference, item]) => externalDataItemToForm(reference, item))
    : []
);

/**
 * Convert form to external data conditions
 *
 * @param {ExternalDataConditionForm[]} form
 * @returns {ExternalDataCondition}
 */
export const formToExternalDataConditions = (form = []) => (
  form.reduce((acc, { type, attribute, value }) => {
    if (!acc[type]) {
      acc[type] = {};
    }

    acc[type][attribute] = value;

    return acc;
  }, {})
);

/**
 * Convert form to external data
 *
 * @param {ExternalDataForm} form
 * @returns {ExternalData}
 */
export const formToExternalData = (form = []) => (
  form.reduce((acc, externalData) => {
    const { type, reference } = externalData;

    const additionalFields = type === EXTERNAL_DATA_TYPES.api
      ? { request: formToRequest(externalData.request) }
      : {
        ...pick(externalData, ['sort', 'sort_by', 'collection']),
        ...formToExternalDataConditions(externalData.conditions),
      };

    acc[reference] = {
      type,

      ...additionalFields,
    };

    return acc;
  }, {})
);
