import { TAG_TYPES } from '@/constants';

/**
 * @typedef { 'imported' | 'created' } TagTypes
 */

/**
 * @typedef {Object & FilterPatterns} AlarmTag
 * @property {string} value
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

/**
 * Check tag is created
 *
 * @param {AlarmTag} tag
 * @returns {boolean}
 */
export const isCreatedTag = tag => tag?.type === TAG_TYPES.created;
