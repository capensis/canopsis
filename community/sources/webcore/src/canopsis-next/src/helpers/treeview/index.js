import { isObject, isArray } from 'lodash';

/**
 * Convert object to tree view object
 *
 * @param parent - Object
 * @param parentKey - Key of object
 * @param [parentPath=''] - Path to parent
 * @param [isParentArray=false] - Is parent converted array into object
 * @return {{children: [string, any], name: *}}
 */
export function convertObjectToTreeview(parent, parentKey, parentPath = '', isParentArray = false) {
  const basePath = !parentPath ? parentKey : parentPath;

  const children = Object.entries(parent).reduce((acc, [key, value]) => {
    const path = isParentArray ? `${basePath}.[${key}]` : `${basePath}.${key}`;

    if (isArray(value)) {
      acc.push({
        ...convertObjectToTreeview({ ...value }, key, path, true),

        isArray: true,
      });
    } else if (isObject(value)) {
      acc.push(convertObjectToTreeview(value, key, path));
    } else {
      acc.push({
        name: key,
        value,
        path,
      });
    }

    return acc;
  }, []);

  return { name: parentKey, children };
}

/**
 * Check if rule is value (not node with children)
 *
 * @param {string} rule
 * @param {Array} operators
 * @returns {number|boolean}
 */
export function isValuePatternRule(rule = '', operators = []) {
  if (!isObject(rule)) {
    return true;
  }

  const items = Object.entries(rule);

  return !!items.length && items.every(([key, value]) => operators.includes(key) && !isObject(value));
}

/**
 * Convert pattern object to treeview items
 *
 * @param {Object} source
 * @param {Array} [operators = []]
 * @param {Array} [prevPath = []]
 */
export function convertPatternToTreeview(source, operators = [], prevPath = []) {
  return Object.entries(source).map(([field, value]) => {
    const path = [...prevPath, field];
    const item = {
      path,

      id: path.join('.'),
      name: field,
    };

    if (isValuePatternRule(value, operators)) {
      item.rule = { field, value };
      item.isSimpleRule = !isObject(value);
    } else {
      item.children = convertPatternToTreeview(value, operators, path);
    }

    return item;
  }, []);
}

/**
 * Convert tree array structure to flat array structure
 *
 * @param {Array} tree
 * @param {string} [itemChildren = 'children']
 * @return {Array}
 */
export const convertTreeArrayToArray = (tree, itemChildren = 'children') => {
  const result = [];

  tree.forEach((item) => {
    const { [itemChildren]: children = [] } = item;

    result.push(item);

    if (children.length) {
      result.push(...convertTreeArrayToArray(children, itemChildren));
    }
  });

  return result;
};
