import { isPlainObject, mapValues } from 'lodash';

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
 * Processes template variable children into a hierarchical structure
 * Recursively converts array-based template variables into nested variable objects
 *
 * @param {Array<{name: string, value: string|Array}>} [children=[]] - Array of template variable children
 * @returns {Array<{text: string, value?: string, variables?: Array}>} Array of processed variable objects where:
 *   - text: The name property from the input child
 *   - value: Present only for non-array values, contains the original value
 *   - variables: Present only for array values, contains recursively processed child elements
 *
 * @example
 * const children = [
 *   { name: 'user', value: 'John' },
 *   { name: 'settings', value: [
 *     { name: 'theme', value: 'dark' },
 *     { name: 'language', value: 'en' }
 *   ]}
 * ];
 * templateVarsChildrenToVariablesProcess(children);
 * // Returns:
 * // [
 * //   { text: 'user', value: 'John' },
 * //   {
 * //     text: 'settings',
 * //     variables: [
 * //       { text: 'theme', value: 'dark' },
 * //       { text: 'language', value: 'en' }
 * //     ]
 * //   }
 * // ]
 */
export const templateVarsChildrenToVariablesProcess = (children = []) => children.map(({ name, value }) => (
  Array.isArray(value)
    ? { text: name, variables: templateVarsChildrenToVariablesProcess(value) }
    : { text: name, value }
));

/**
 * Converts template variables object into a hierarchical variables structure
 * Maps each template variable group through the children processing function
 *
 * @param {Object} [templateVars={}] - Object containing template variable groups
 * @returns {Object} Object with the same keys as input, where each value is processed
 *                   through templateVarsChildrenToVariablesProcess
 *
 * @example
 * const templateVars = {
 *   user: [
 *     { name: 'name', value: 'John' },
 *     { name: 'email', value: 'john@example.com' }
 *   ],
 *   system: [
 *     { name: 'version', value: '1.0.0' }
 *   ]
 * };
 * templateVarsToVariables(templateVars);
 * // Returns:
 * // {
 * //   user: [
 * //     { text: 'name', value: 'John' },
 * //     { text: 'email', value: 'john@example.com' }
 * //   ],
 * //   system: [
 * //     { text: 'version', value: '1.0.0' }
 * //   ]
 * // }
 */
export const templateVarsToVariables = (templateVars = {}) => (
  mapValues(templateVars, templateVarsChildrenToVariablesProcess)
);
