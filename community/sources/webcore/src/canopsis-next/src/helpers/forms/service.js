import { isNumber, omit } from 'lodash';

import { ENTITIES_STATES, OLD_PATTERNS_FIELDS, PATTERNS_FIELDS } from '@/constants';

import { infosToArray } from './shared/common';
import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

/**
 * @typedef {FilterPatterns} Service
 * @property {string} [_id]
 * @property {string} name
 * @property {string} category
 * @property {boolean} enabled
 * @property {number} impact_level
 * @property {number} sli_avail_state
 * @property {Object|Array} infos
 * @property {Object} entity_patterns
 * @property {string} output_template
 * @property {Object} [coordinates]
 */

/**
 * @typedef {Service} ServiceForm
 * @property {FilterPatternsForm} patterns
 * @property {Object} category
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
  output_template: service.output_template ?? '',
  sli_avail_state: service.sli_avail_state ?? ENTITIES_STATES.ok,
  patterns: filterPatternsToForm(
    service,
    [PATTERNS_FIELDS.entity],
    [OLD_PATTERNS_FIELDS.entity],
  ),
  coordinates: service.coordinates ?? {
    lat: undefined,
    lng: undefined,
  },
});

/**
 * Convert entity service form to entity service object by stack
 *
 * @param {ServiceForm} [form = {}]
 * @returns {Service}
 */
export const formToService = (form = {}) => {
  const service = {
    ...omit(form, ['patterns', 'coordinates']),
    ...formFilterToPatterns(form.patterns, [PATTERNS_FIELDS.entity]),
    category: form.category._id,
  };

  if (isNumber(form.coordinates.lat) && isNumber(form.coordinates.lng)) {
    service.coordinates = form.coordinates;
  }

  return service;
};
