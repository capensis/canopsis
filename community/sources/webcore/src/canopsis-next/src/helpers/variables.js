import { isPlainObject } from 'lodash';

/**
 * Converts an object into an array of variable objects with text and value properties
 * Recursively processes nested objects to create a hierarchical structure
 *
 * @param {Object} [obj={}] - The input object to convert
 * @returns {Array<{text: string, value: string, variables?: Array}>} Array of variable objects where:
 *   - text: The property name from the input object
 *   - value: Same as the text property
 *   - variables: Present only for nested objects, contains recursively processed child properties
 *
 * @example
 * const input = {
 *   name: 'John',
 *   address: {
 *     city: 'London',
 *     country: 'UK'
 *   }
 * };
 * objectToVariables(input);
 * // Returns:
 * // [
 * //   { text: 'name', value: 'name' },
 * //   {
 * //     text: 'address',
 * //     value: 'address',
 * //     variables: [
 * //       { text: 'city', value: 'city' },
 * //       { text: 'country', value: 'country' }
 * //     ]
 * //   }
 * // ]
 */
export const objectToVariables = (obj = {}) => Object.entries(obj).map(([text, value]) => (
  isPlainObject(value) ? { text, value: text, variables: objectToVariables(value) } : { text, value: text }
));

/**
 * Checks if a variables list has at least one variable.
 *
 * Recursively checks nested variable structures to determine if any variable exists.
 *
 * @param {Array} [variables=[]] - Array of variable objects to check.
 * @returns {boolean} True if at least one variable exists, false otherwise.
 *
 * @example
 * hasAtLeastOneVariable([]);
 * // Returns: false
 *
 * hasAtLeastOneVariable([{ text: 'name', value: 'John' }]);
 * // Returns: true
 *
 * hasAtLeastOneVariable([{ text: 'user', variables: [] }]);
 * // Returns: false
 *
 * hasAtLeastOneVariable([{ text: 'user', variables: [{ text: 'name', value: 'John' }] }]);
 * // Returns: true
 */
export const hasAtLeastOneVariable = (variables = []) => (
  !variables.length
    ? false
    : variables.some((variable = {}) => (
      variable.variables
        ? hasAtLeastOneVariable(variable.variables)
        : true
    )));
