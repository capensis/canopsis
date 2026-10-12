import { isArray, isObject } from 'lodash';

/**
 * Build unique treeview path. Empty object keys must stay in the path,
 * otherwise parent and child share the same id (e.g. alarm.links + "").
 *
 * @param {string} basePath
 * @param {string} key
 * @param {boolean} isParentArray
 * @returns {string}
 */
const getObjectTreeviewPath = (basePath, key, isParentArray) => {
  const keySegment = isParentArray ? `[${key}]` : key;

  if (!basePath) {
    return keySegment;
  }

  if (!isParentArray && key === '') {
    return `${basePath}['']`;
  }

  return `${basePath}.${keySegment}`;
};

/**
 * Convert object to tree view object
 *
 * @param {Object} parent - Object
 * @param {string} parentKey - Key of object
 * @param {string} [parentPath=''] - Path to parent
 * @param {boolean} [isParentArray=false] - Is parent converted array into object
 * @returns {{ children: Array, name: string, path: string }}
 */
export const convertObjectToTreeview = (
  parent,
  parentKey,
  parentPath = '',
  isParentArray = false,
) => {
  const basePath = !parentPath ? parentKey : parentPath;

  const children = Object.entries(parent).reduce((acc, [key, value]) => {
    const path = getObjectTreeviewPath(basePath, key, isParentArray);

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

  return { name: parentKey, path: basePath, children };
};

/**
 * Convert tree array structure to flat array structure
 *
 * @param {Array} tree
 * @param {string} [itemChildren='children']
 * @returns {Array}
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

/**
 * Convert treeview nodes to object
 *
 * @param {Object} node
 * @returns {Object | undefined}
 */
export const convertTreeviewToObject = (node) => {
  if (node.isArray) {
    return node.children?.map(convertTreeviewToObject);
  }

  if (node.children) {
    return node.children?.reduce((acc, item) => {
      acc[item.name] = convertTreeviewToObject(item);

      return acc;
    }, {});
  }

  return node.value;
};
