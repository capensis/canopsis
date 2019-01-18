import get from 'lodash/get';

export function checkUserAccess(user, rightId, rightMask) {
  const checksum = get(user, ['rights', rightId, 'checksum'], 0);

  return (checksum & rightMask) === rightMask;
}

export default {
  checkUserAccess,
};
