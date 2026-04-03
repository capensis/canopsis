import { isPlainObject, mapValues } from 'lodash';

/**
 * Wraps a variable name in handlebars template syntax
 *
 * @param {string} value - The variable name to wrap in template syntax
 * @returns {string} The variable name wrapped in double curly braces (handlebars format)
 *
 * @example
 * variableTemplatePreparer('user.name');
 * // Returns: '{{ user.name }}'
 */
export const variableTemplatePreparer = value => `{{ ${value} }}`;

/**
 * Converts an object into an array of variable objects with text and value properties
 * Recursively processes nested objects to create a hierarchical structure
 *
 * @param {Object} [obj={}] - The input object to convert
 * @param {Function} [valuePreparer=variableTemplatePreparer] - A function to prepare the value before
 *                                                              adding to the variable
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
export const objectToVariables = (
  obj = {},
  valuePreparer = variableTemplatePreparer,
) => Object.entries(obj).map(([text, value]) => {
  const preparedValue = valuePreparer ? valuePreparer(text) : text;

  return isPlainObject(value)
    ? {
      text,
      value: preparedValue,
      variables: objectToVariables(value, childValue => valuePreparer(`${text}.${childValue}`)),
    }
    : {
      text,
      value: preparedValue,
    };
});

/**
 * Processes variable children into a hierarchical structure
 * Recursively converts array-based variables into nested variable objects
 *
 * @param {Array<{name: string, value: string|Array}>} [children=[]] - Array of variable children
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
 * varsChildrenToVariablesProcess(children);
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
export const varsChildrenToVariablesProcess = (children = []) => children.map(({ alias, name, value }) => (
  Array.isArray(value)
    ? { alias, text: name, variables: varsChildrenToVariablesProcess(value) }
    : { alias, text: name, value }
));

/**
 * Converts variables object into a hierarchical variables structure
 * Maps each template variable group through the children processing function
 *
 * @param {Object} [templateVars={}] - Object containing variable groups
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
 * varsToVariables(templateVars);
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
export const varsToVariables = (templateVars = {}) => (
  mapValues(templateVars, varsChildrenToVariablesProcess)
);

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
