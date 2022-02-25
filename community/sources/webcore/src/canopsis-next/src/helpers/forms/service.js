import { cloneDeep } from 'lodash';

import { ENTITIES_STATES } from '@/constants';

import { infosToArray } from './shared/common';

/**
 * @typedef {Object} ServiceForm
 * @property {string} name
 * @property {string} category
 * @property {boolean} enabled
 * @property {number} impact_level
 * @property {number} sli_avail_state
 * @property {Object|Array} infos
 * @property {Object} entity_patterns
 * @property {string} output_template
 */

/**
 * @typedef {ServiceForm} Service
 * @property {string} _id
 */

/**
 * @typedef {Service} ServiceForm
 */

/**
 * Convert entity service object to entity service form
 *
 * @param {Service} [service = {}]
 * @returns {ServiceForm}
 */
export const serviceToForm = (service = {}) => ({
  impact_level: service.impact_level,
  name: service.name ?? '',
  category: service.category ?? '',
  enabled: service.enabled ?? true,
  infos: infosToArray(service.infos),
  entity_patterns: service.entity_patterns ? cloneDeep(service.entity_patterns) : [],
  output_template: service.output_template ?? '',
  sli_avail_state: service.sli_avail_state ?? ENTITIES_STATES.ok,
});

/**
 * Convert entity service form to entity service object by stack
 *
 * @param {ServiceForm} [form = {}]
 * @returns {Service}
 */
export const formToService = (form = {}) => ({
  ...form,
  category: form.category._id,
});
