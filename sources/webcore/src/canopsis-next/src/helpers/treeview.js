import { isObject, isArray } from 'lodash';

export default function convertObjectFieldToTreeBranch(branch, branchName, prevPath = '') {
  const children = Object.keys(branch).reduce((acc, field) => {
    if (isArray(branch[field]) || !isObject(branch[field])) {
      const path = prevPath ? `${prevPath}.${branchName}.${field}` : `${branchName}.${field}`;
      acc.push({
        name: field,
        value: branch[field],
        path,
        isArray: isArray(branch[field]),
      });
    } else {
      const path = prevPath ? `${prevPath}.${branchName}` : `${branchName}`;
      acc.push(convertObjectFieldToTreeBranch(branch[field], field, path));
    }

    return acc;
  }, []);

  return { name: branchName, children };
}

