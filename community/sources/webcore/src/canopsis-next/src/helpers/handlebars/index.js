import promisedHandlebars from 'promised-handlebars';
import HandlebarsLib from 'handlebars';

import store from '@/store';

import * as helpers from './helpers';

const Handlebars = promisedHandlebars(HandlebarsLib);

/**
 * Compile template
 *
 * @param {string} template
 * @param {Object} [context = {}]
 * @returns {Promise}
 */
export async function compile(template, context = {}) {
  const handleBarFunction = Handlebars.compile(template ?? '');
  const preparedContext = {
    env: store.getters['templateVars/items'] ?? {},

    ...context,
  };

  const result = await handleBarFunction(preparedContext);

  const element = document.createElement('div');

  element.innerHTML = result;

  return element.innerHTML;
}

/**
 * Register handlebars helper
 *
 * @param {string} name
 * @param {Function} helper
 * @returns {*}
 */
export function registerHelper(name, helper) {
  if (Handlebars.helpers[name]) {
    return;
  }

  Handlebars.registerHelper(name, helper);
}

/**
 * Unregister handlebars helper
 *
 * @param {string} name
 * @returns {*}
 */
export function unregisterHelper(name) {
  Handlebars.unregisterHelper(name);
}

/**
 * Get all node variables
 *
 * @param {
 *   hbs.AST.BlockStatement |
 *   hbs.AST.MustacheStatement |
 *   hbs.AST.Program |
 *   hbs.AST.PathExpression |
 *   hbs.AST.SubExpression |
 *   hbs.AST.HashPair |
 *   hbs.AST.ContentStatement |
 *   hbs.AST.Statement
 * } node
 * @returns {string[]}
 */
const getVariablesFromNode = (node) => {
  switch (node?.type) {
    case 'MustacheStatement':
    case 'BlockStatement': {
      const variables = [];

      if (node.hash?.pairs) {
        node.hash.pairs.forEach((item) => {
          variables.push(...getVariablesFromNode(item));
        });
      }

      if (node.program) {
        variables.push(...getVariablesFromNode(node.program));
      }

      if (node.params) {
        node.params.forEach((item) => {
          variables.push(...getVariablesFromNode(item));
        });
      }

      if (node.path) {
        variables.push(...getVariablesFromNode(node.path));
      }

      return variables;
    }
    case 'Program':
      return node.body?.reduce((acc, bodyNode) => {
        acc.push(...getVariablesFromNode(bodyNode));

        return acc;
      }, []) ?? [];
    case 'PathExpression':
      return node?.original ? [node.original] : [];
    case 'SubExpression':
      return node.params.reduce((acc, item) => {
        acc.push(...getVariablesFromNode(item));

        return acc;
      }, []);
    case 'HashPair':
      return getVariablesFromNode(node.value);
    default:
      return [];
  }
};

/**
 * Get all using variables in the template
 *
 * @param {string} template
 * @returns {string[]}
 */
export const getTemplateVariables = template => getVariablesFromNode(
  Handlebars.parseWithoutProcessing(template),
)
  .filter(variable => !Handlebars.helpers[variable]);

/**
 * Register global helpers
 */
registerHelper('duration', helpers.durationHelper);
registerHelper('state', helpers.alarmStateHelper);
registerHelper('tags', helpers.alarmTagsHelper);
registerHelper('request', helpers.requestHelper);
registerHelper('timestamp', helpers.timestampHelper);
registerHelper('internal-link', helpers.internalLinkHelper);
registerHelper('compare', helpers.compareHelper);
registerHelper('concat', helpers.concatHelper);
registerHelper('sum', helpers.sumHelper);
registerHelper('minus', helpers.minusHelper);
registerHelper('mul', helpers.mulHelper);
registerHelper('divide', helpers.divideHelper);
registerHelper('capitalize', helpers.capitalizeHelper);
registerHelper('capitalize-all', helpers.capitalizeAllHelper);
registerHelper('lowercase', helpers.lowercaseHelper);
registerHelper('uppercase', helpers.uppercaseHelper);
registerHelper('replace', helpers.replaceHelper);
registerHelper('copy', helpers.copyHelper);
registerHelper('json', helpers.jsonHelper);
