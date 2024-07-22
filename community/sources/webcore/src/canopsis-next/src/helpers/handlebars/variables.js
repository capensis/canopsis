import { Handlebars } from './handlebars';

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

      if (node.inverse) {
        variables.push(...getVariablesFromNode(node.inverse));
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
).filter(variable => !Handlebars.helpers[variable]);
