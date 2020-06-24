import { get, isUndefined, omit, sortBy } from 'lodash';
import flatten from 'flat';

import {
  NOT_COMPLETED_USER_RIGHTS,
  USERS_RIGHTS,
  USERS_RIGHTS_MASKS,
  USERS_RIGHTS_TYPES_TO_MASKS,
  USER_RIGHTS_PREFIXES,
} from '@/constants';

/**
 * Check user access for a right
 *
 * @param {Object} user
 * @param {string} rightId
 * @param {number} rightMask
 * @returns {boolean}
 */
export function checkUserAccess(user, rightId, rightMask) {
  const checksum = get(user, ['rights', rightId, 'checksum'], 0);

  return (checksum & rightMask) === rightMask;
}

/**
 * Get masks for right by right.type
 *
 * @param {Object} right
 * @returns {*}
 */
export function getRightMasks(right) {
  return USERS_RIGHTS_TYPES_TO_MASKS[right.type]
    ? USERS_RIGHTS_TYPES_TO_MASKS[right.type]
    : [USERS_RIGHTS_MASKS.default];
}

/**
 * Get right checkbox value for the form
 *
 * @param {Object} right
 * @param {Object} role
 * @param {Object} changedRole
 * @param {number} rightMask
 * @returns {boolean}
 */
export function getCheckboxValue(right, role, changedRole, rightMask = USERS_RIGHTS_MASKS.default) {
  const checkSum = get(role, ['rights', right._id, 'checksum'], 0);

  const changedCheckSum = get(changedRole, [right._id]);
  const currentCheckSum = isUndefined(changedCheckSum) ? checkSum : changedCheckSum;
  const rightType = currentCheckSum & rightMask;

  return rightType === rightMask;
}

/**
 * Get prepared grouped rights for the rights page
 *
 * @param {Array<Object>} rights
 * @param {Array<Object>} [views = []]
 * @param {Array<Object>} [playlists = []]
 * @returns {*}
 */
export function getGroupedRights(rights, views = [], playlists = []) {
  const allBusinessRightsIds = flatten(USERS_RIGHTS.business);

  const { exploitation: exploitationTechnicalRights, ...adminTechnicalRights } = USERS_RIGHTS.technical;
  const adminTechnicalRightsValues = Object.values(adminTechnicalRights);
  const exploitationTechnicalRightsValues = Object.values(exploitationTechnicalRights);

  const groupedRights = rights.reduce((acc, right) => {
    const rightId = String(right._id, '\'');
    const view = views.find(({ _id }) => _id === rightId);
    const playlist = playlists.find(({ _id }) => _id === rightId);

    if (view) {
      acc.view.push({
        ...right,

        desc: right.desc.replace(view._id, view.name),
      });
    } else if (playlist) {
      acc.playlist.push({
        ...right,

        desc: right.desc.replace(playlist._id, playlist.name),
      });
    } else if (adminTechnicalRightsValues.indexOf(rightId) !== -1) {
      acc.technical.admin.push(right);
    } else if (exploitationTechnicalRightsValues.indexOf(rightId) !== -1) {
      acc.technical.exploitation.push(right);
    } else if (
      Object.values(allBusinessRightsIds).indexOf(rightId) !== -1 ||
      NOT_COMPLETED_USER_RIGHTS.some(id => rightId.startsWith(id))
    ) {
      const [parentKey] = right._id.split('_');

      if (acc.business[parentKey]) {
        acc.business[parentKey].push(right);
      }
    }

    return acc;
  }, {
    business: {
      [USER_RIGHTS_PREFIXES.business.common]: [],
      [USER_RIGHTS_PREFIXES.business.alarmsList]: [],
      [USER_RIGHTS_PREFIXES.business.context]: [],
      [USER_RIGHTS_PREFIXES.business.weather]: [],
      [USER_RIGHTS_PREFIXES.business.counter]: [],
    },
    view: [],
    playlist: [],
    technical: {
      admin: [],
      exploitation: [],
    },
  });

  /**
   * We are using order which one we've defined on the reduce accumulator initial value.
   * For not `number`/`number string` object keys ordering is staying like we define
   */
  groupedRights.business = Object.entries(groupedRights.business)
    .map(([key, groupRights]) => ({ key: `rights.business.${key}`, rights: sortBy(groupRights, ['desc']) }));

  /**
   * Ordering behavior have the same behavior as for `business`
   */
  groupedRights.technical = Object.entries(groupedRights.technical)
    .map(([key, groupRights]) => ({ key: `rights.technical.${key}`, rights: sortBy(groupRights, ['desc']) }));

  groupedRights.view.push(...groupedRights.playlist);

  return omit(groupedRights, ['playlist']);
}
