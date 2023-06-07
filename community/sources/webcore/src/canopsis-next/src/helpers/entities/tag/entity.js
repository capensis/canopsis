import { TAG_TYPES } from '@/constants/tag';

/**
 * @typedef { 'imported' | 'created' } TagTypes
 */

/**
 * @typedef {Object & FilterPatterns} AlarmTag
 * @property {string} name
 * @property {string} color
 * @property {TagTypes} [type]
 */

/**
 * Check tag is imported
 *
 * @param {AlarmTag} tag
 * @returns {boolean}
 */
export const isImportedTag = tag => tag?.type === TAG_TYPES.imported;
