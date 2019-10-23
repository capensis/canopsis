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
