import { isPlainObject, mapValues } from 'lodash';

/**
 * Wraps a variable name in handlebars template syntax.
 *
 * @param {string} value - The variable name to wrap in template syntax.
 * @returns {string} The variable name wrapped in double curly braces.
 */
export const variableTemplatePreparer = value => `{{ ${value} }}`;

/**
 * Converts an object into an array of variable objects with text and value properties.
 *
 * @param {Object} [obj={}] - The input object to convert.
 * @param {Function} [valuePreparer=variableTemplatePreparer] - Value preparing function.
 * @returns {Array<{text: string, value: string, variables?: Array}>}
 */
export const objectToVariables = (obj = {}, valuePreparer = variableTemplatePreparer) => Object.entries(obj).map(([text, value]) => {
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
 * Processes variable children into a hierarchical structure.
 *
 * @param {Array<{name: string, value: string|Array}>} [children=[]]
 * @returns {Array<{text: string, value?: string, variables?: Array}>}
 */
export const varsChildrenToVariablesProcess = (children = []) => children.map(({ alias, name, value }) => Array.isArray(value)
      ? { alias, text: name, variables: varsChildrenToVariablesProcess(value) }
      : { alias, text: name, value },
  );

/**
 * Converts variables object into a hierarchical variables structure.
 *
 * @param {Object} [templateVars={}] - Object containing variable groups.
 * @returns {Object}
 */
export const varsToVariables = (templateVars = {}) => mapValues(templateVars, varsChildrenToVariablesProcess);

/**
 * Checks if a variables list has at least one variable.
 *
 * @param {Array} [variables=[]] - Array of variable objects to check.
 * @returns {boolean}
 */
export const hasAtLeastOneVariable = (variables = []) => !variables.length
    ? false
    : variables.some((variable = {}) => variable.variables ? hasAtLeastOneVariable(variable.variables) : true,
      );
