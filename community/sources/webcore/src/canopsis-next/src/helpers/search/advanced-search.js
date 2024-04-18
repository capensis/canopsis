import { orderBy, map } from 'lodash';

import {
  ADVANCED_SEARCH_ITEM_TYPES,
  ADVANCED_SEARCH_NEXT_ITEM_TYPES,
  ADVANCED_SEARCH_UNION_CONDITIONS,
  ADVANCED_SEARCH_UNION_REGEXP_PATTERN,
  ADVANCED_SEARCH_NOT,
  ADVANCED_SEARCH_CONDITIONS,
} from '@/constants';

import { uid } from '@/helpers/uid';

/**
 * @typedef { 'field' | 'condition' | 'value' | 'union' } AdvancedSearchItemType
 */

/**
 * @typedef {Object} AdvancedSearchField
 * @property {string} value
 * @property {string} text
 * @property {string} [selectorText]
 * @property {AdvancedSearchItem[]} [items]
 */

/**
 * @typedef {AdvancedSearchField} AdvancedSearchListItem
 * @property {AdvancedSearchItemType} type
 */

/**
 * @typedef {AdvancedSearchListItem} AdvancedSearchItem
 * @property {string} key
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
 * Converts an advanced search string into an array of search items, each represented as an object.
 * This function parses a search string based on predefined union conditions, item types, and conditions.
 * It supports complex queries with logical operators and is designed to work with a dynamic set of column names.
 *
 * @param {string} search - The search string to be parsed.
 * @param {Array<{text: string}>} fields - An array of objects representing the fields, where each column has a `text`
 * property.
 * @returns {{internalSearch: string, value: Array}} An object containing two properties: `internalSearch` which is
 * a string that could not be parsed, and `value` which is an array of parsed items.
 *
 * @example
 * // Example usage:
 * const search = "name = John AND age > 30";
 * const columns = [{ text: "name" }, { text: "age" }];
 * const result = advancedSearchStringToArray(search, columns);
 * console.log(result);
 * // Output:
 * // {
 * //   internalSearch: ``,
 * //   value: [
 * //     { key: `uid1`, value: `name`, type: `field`, text: `name`, not: false },
 * //     { key: `uid2`, value: `=`, type: `condition`, text: `=` },
 * //     { key: `uid3`, value: `John`, type: `value`, text: `John` },
 * //     { key: `uid4`, value: `AND`, type: `union`, text: `AND` },
 * //     { key: `uid5`, value: `age`, type: `field`, text: `age`, not: false },
 * //     { key: `uid6`, value: `>`, type: `condition`, text: `>` },
 * //     { key: `uid7`, value: `30`, type: `value`, text: `30` }
 * //   ]
 * // }
 */
export const advancedSearchStringToArray = (search = '', fields = []) => {
  const result = {
    internalSearch: '',
    value: [],
  };

  if (!search) {
    return result;
  }

  const searchWithoutDash = search.replace(/^\s*-\s*/, '');

  try {
    const items = searchWithoutDash.split(ADVANCED_SEARCH_UNION_REGEXP_PATTERN);
    const columnsForRegexp = orderBy(map(fields, 'text'), text => text.length, ['desc']).join('|');
    const itemRegexp = new RegExp(`^(?<not>${ADVANCED_SEARCH_NOT})?\\s*(?<field>(${columnsForRegexp})[\\w._]*|[\\w._]+)?\\s*(?<condition>${Object.values(ADVANCED_SEARCH_CONDITIONS).join('|')})?\\s*(?<value>.+)?$`, 'i');

    for (let i = 0; i < items.length; i += 1) {
      const item = items[i];

      const trimmedItem = item.trim();

      if (!trimmedItem) {
        continue;
      }

      const unionItem = ADVANCED_SEARCH_UNION_CONDITIONS[trimmedItem.toLocaleLowerCase()];

      if (unionItem) {
        result.value.push({
          key: uid(),
          value: unionItem,
          type: ADVANCED_SEARCH_ITEM_TYPES.union,
          text: unionItem,
        });

        continue;
      }

      const { groups } = trimmedItem.match(itemRegexp);

      if (!groups) {
        result.internalSearch = items.slice(i).join(' ');

        break;
      }

      const {
        not,
        field,
        condition,
        value,
      } = groups;

      if (i !== items.length - 1 && (!field || !condition || !value)) {
        result.internalSearch = items.slice(i).join(' ');

        break;
      }

      if (!field) {
        result.internalSearch = items.slice(i).join(' ');

        break;
      }

      result.value.push({
        key: uid(),
        value: field,
        type: ADVANCED_SEARCH_ITEM_TYPES.field,
        text: field,
        not: !!not,
      });

      if (!condition) {
        result.internalSearch = items.slice(i).map((slicedItem, index) => (index ? slicedItem : value)).join(' ');

        break;
      }

      result.value.push({
        key: uid(),
        value: condition,
        type: ADVANCED_SEARCH_ITEM_TYPES.condition,
        text: condition,
      });

      if (value) {
        result.value.push({
          key: uid(),
          value,
          type: ADVANCED_SEARCH_ITEM_TYPES.value,
          text: value,
        });
      }
    }

    return result;
  } catch (err) {
    console.error(err);

    return {
      internalSearch: searchWithoutDash,
      value: [],
    };
  }
};

/**
 * Converts an array of search items into a formatted string for advanced search.
 * Each item in the array can optionally have a `not` property to prepend `NOT` to the item text.
 *
 * @param {Object<{ not: boolean, text: string, selectorText: string}>[]} array - The array of search items.
 * Each item should have either `selectorText` or `text`.
 * @returns {string} The formatted search string starting with a dash and spaces between items.
 *
 * @example
 * // returns "- NOT apple banana"
 * advancedSearchArrayToString([
 *   { not: true, text: `apple` },
 *   { text: `banana` }
 * ]);
 */
export const advancedSearchArrayToString = (array = []) => (
  `- ${array.map(item => `${item.not ? `${ADVANCED_SEARCH_NOT} ` : ''}${item.selectorText || item.text}`).join(' ')}`
);

/**
 * Transforms an array of field objects by adding a `type` property with a value of `field` to each object.
 * This is used to prepare fields for advanced search functionality, ensuring each field object includes
 * the necessary type information.
 *
 * @param {AdvancedSearchField[]} [fields=[]] - An array of objects representing the fields to be
 * prepared.
 * Each object should have at least a `value` and a `text` property.
 * @returns {AdvancedSearchListItem[]} An array of the same objects provided in the input,
 * but with an additional `type` property set to `field`.
 */
export const prepareAdvancedSearchFields = (fields = []) => (
  fields.map(field => ({
    ...field,
    type: ADVANCED_SEARCH_ITEM_TYPES.field,
    items: field.items?.length ? prepareAdvancedSearchFields(field.items) : undefined,
  }))
);

/**
 * Prepares advanced search conditions by mapping each condition to an object with its value, type, and text.
 * The type is always set to the `condition` type from `ADVANCED_SEARCH_ITEM_TYPES`.
 *
 * @param {string[]} [conditions=[]] - An array of conditions to be prepared. Defaults to an empty array if not
 * provided.
 * @returns {AdvancedSearchListItem[]} An array of objects where each object represents a prepared condition with
 * properties: value, type, and text.
 */
export const prepareAdvancedSearchConditions = (conditions = []) => (
  conditions.map(condition => ({ value: condition, type: ADVANCED_SEARCH_ITEM_TYPES.condition, text: condition }))
);
