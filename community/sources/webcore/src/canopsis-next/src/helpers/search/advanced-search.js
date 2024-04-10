import {
  ADVANCED_SEARCH_ITEM_TYPES,
  ADVANCED_SEARCH_NEXT_ITEM_TYPES,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ADVANCED_SEARCH_UNION_REGEXP_PATTERN,
  ADVANCED_SEARCH_REGEXP_PATTERN,
} from '@/constants';

import { uid } from '@/helpers/uid';

/**
 * @typedef { 'field' | 'condition' | 'value' | 'union' } AdvancedSearchItemType
 */

/**
 * @typedef {Object} AdvancedSearchItem
 * @property {string} key
 * @property {AdvancedSearchItemType} type
 * @property {string} value
 * @property {string} text
 * @property {boolean} [not]
 */

/**
 * Checks if the provided type corresponds to the `value` type in the ADVANCED_SEARCH_ITEM_TYPES enumeration.
 *
 * @param {AdvancedSearchItemType} type - The type to check against the `value` type.
 * @returns {boolean} Returns `true` if the type is equal to the `value` type, otherwise returns `false`.
 */
export const isValueType = type => type === ADVANCED_SEARCH_ITEM_TYPES.value;

/**
 * Checks if the provided type corresponds to the `value` type in the ADVANCED_SEARCH_ITEM_TYPES enumeration.
 *
 * @param {AdvancedSearchItemType} type - The type to check against the `value` type.
 * @returns {boolean} Returns `true` if the type is equal to the `field` type, otherwise returns `false`.
 */
export const isFieldType = type => type === ADVANCED_SEARCH_ITEM_TYPES.field;

/**
 * Determines the next advanced search item type based on the current type.
 * If the current type is `union`, it returns `field`. Otherwise, it increments the type by 1.
 * This function is useful for cycling through the advanced search item types in a predefined order.
 *
 * @param {AdvancedSearchItemType} [type = ADVANCED_SEARCH_ITEM_TYPES.union] - The current type of the advanced search
 * item. Defaults to `union` if not specified.
 * @returns {number} The next advanced search item type.
 */
export const getNextAdvancedSearchType = (type = ADVANCED_SEARCH_ITEM_TYPES.union) => (
  ADVANCED_SEARCH_NEXT_ITEM_TYPES[type] ?? ADVANCED_SEARCH_ITEM_TYPES.field
);

/**
 * Parses a string representing an advanced search query into an array of search items.
 * Each search item is an object that may represent a field, condition, value, or a logical union (AND/OR).
 * The function splits the input string based on logical unions, then further analyzes each part
 * to classify it into one of the item types. It supports negation, fields, conditions, and values.
 *
 * @param {string} [search = ``] - The advanced search query string to be parsed.
 * @returns {Array<AdvancedSearchItem>} An array of objects, each representing a part of the parsed search query.
 * Each object includes a unique key, the item's value, its type (field, condition, value, union),
 * the original text, and a flag indicating negation (for fields).
 *
 * @example
 * // Example usage of parseAdvancedSearch
 * const searchQuery = "-field1 = value1 AND field2 != `value 2`";
 * const parsedQuery = parseAdvancedSearch(searchQuery);
 * console.log(parsedQuery);
 * // Output might look like:
 * // [
 * //   { key: `uid1`, value: `field1`, type: 'field', text: `field1`, not: true },
 * //   { key: `uid2`, value: `=`, type: 'condition', text: `=`, not: false },
 * //   { key: `uid3`, value: `value1`, type: 'value', text: `value1`, not: false },
 * //   { key: `uid4`, value: `AND`, type: 'union', text: `AND`, not: false },
 * //   { key: `uid5`, value: `field2`, type: 'field', text: `field2`, not: false },
 * //   { key: `uid6`, value: `!=`, type: 'condition', text: `!=`, not: false },
 * //   { key: `uid7`, value: "'value 2`", type: 'value', text: "'value 2`", not: false }
 * // ]
 */
export const parseAdvancedSearch = (search = '') => search.replace(/^\s*-\s*/, '')
  .split(ADVANCED_SEARCH_UNION_REGEXP_PATTERN)
  .reduce((acc, item) => {
    const trimmedItem = item.trim();

    if (!trimmedItem) {
      return acc;
    }

    const unionItem = ADVANCED_SEARCH_UNION_CONDITIONS[trimmedItem.toLocaleLowerCase()];

    if (unionItem) {
      acc.push({
        key: uid(),
        value: unionItem,
        type: ADVANCED_SEARCH_ITEM_TYPES.union,
        text: unionItem,
      });

      return acc;
    }

    const { groups } = trimmedItem.match(ADVANCED_SEARCH_REGEXP_PATTERN);

    if (groups) {
      const {
        not,
        field,
        condition,
        value = '\'\'',
      } = groups;

      if (!field || !condition) {
        throw new Error('Incorrect search');
      }

      acc.push(
        {
          key: uid(),
          value: field,
          type: ADVANCED_SEARCH_ITEM_TYPES.field,
          text: field,
          not: !!not,
        },
        {
          key: uid(),
          value: condition,
          type: ADVANCED_SEARCH_ITEM_TYPES.condition,
          text: condition,
        },
        {
          key: uid(),
          value,
          type: ADVANCED_SEARCH_ITEM_TYPES.value,
          text: value,
        },
      );
    }

    return acc;
  }, []);

/**
 * Transforms an array of field objects by adding a `type` property with a value of `field` to each object.
 * This is used to prepare fields for advanced search functionality, ensuring each field object includes
 * the necessary type information.
 *
 * @param {Array<{value: string, text: string}>} [fields=[]] - An array of objects representing the fields to be
 * prepared.
 * Each object should have at least a `value` and a `text` property.
 * @returns {Array<{value: string, text: string, type: string}>} An array of the same objects provided in the input,
 * but with an additional `type` property set to `field`.
 */
export const prepareAdvancedSearchFields = (fields = []) => (
  fields.map(({ value, text }) => ({ value, text, type: ADVANCED_SEARCH_ITEM_TYPES.field }))
);

/**
 * Prepares advanced search conditions by mapping each condition to an object with its value, type, and text.
 * The type is always set to the `condition` type from `ADVANCED_SEARCH_ITEM_TYPES`.
 *
 * @param {Array} [conditions=[]] - An array of conditions to be prepared. Defaults to an empty array if not provided.
 * @returns {Array} An array of objects where each object represents a prepared condition with properties: value,
 * type, and text.
 */
export const prepareAdvancedSearchConditions = (conditions = []) => (
  conditions.map(condition => ({ value: condition, type: ADVANCED_SEARCH_ITEM_TYPES.condition, text: condition }))
);
