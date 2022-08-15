import { MAP_TYPES, MERMAID_THEMES } from '@/constants';

import { addKeyInEntities, removeKeyFromEntities } from '@/helpers/entities';

/**
 * @typedef {Object} MapCommonFields
 * @property {string} name
 * @property {string} [_id]
 */

/**
 * @typedef {MapCommonFields} MapMermaidProperties
 * @property {string} theme
 * @property {string} code
 */

/**
 * @typedef {MapMermaidProperties} MapMermaidPropertiesForm
 */

/**
 * @typedef {MapCommonFields} MapMermaid
 * @property {'mermaid'} type
 * @property {MapMermaidProperties} properties
 */

/**
 * @typedef {Object} MapGeoProperties
 */

/**
 * @typedef {MapGeoProperties} MapGeoPropertiesForm
 */

/**
 * @typedef {MapCommonFields} MapGeo
 * @property {'geo'} type
 * @property {MapGeoProperties} properties
 */

/**
 * @typedef {Object} MapFlowchartProperties
 */

/**
 * @typedef {MapFlowchartProperties} MapFlowchartPropertiesForm
 */

/**
 * @typedef {MapCommonFields} MapFlowchart
 * @property {'flowchart'} type
 * @property {MapFlowchartProperties} properties
 */

/**
 * @typedef {Object} MapTreeOfDependenciesEntity
 * @property {Entity} data
 * @property {Entity[]} pinned
 */

/**
 * @typedef {MapTreeOfDependenciesEntity} MapTreeOfDependenciesEntityForm
 * @property {string} key
 */

/**
 * @typedef {Object} MapTreeOfDependenciesProperties
 * @property {MapTreeOfDependenciesEntity[]} entities
 * @property {boolean} impact
 */

/**
 * @typedef {MapTreeOfDependenciesProperties} MapTreeOfDependenciesPropertiesForm
 * @property {MapTreeOfDependenciesEntityForm[]} entities
 */

/**
 * @typedef {MapCommonFields} MapTreeOfDependencies
 * @property {'treeOfDependencies'} type
 * @property {MapTreeOfDependenciesProperties} properties
 */

/**
 * @typedef {MapMermaid} Map
 */

/**
 * @typedef {MapMermaid} MapForm
 */

/**
 * Convert map geo properties object to form
 *
 * @param {MapGeoProperties} [properties = {}]
 * @returns {MapGeoPropertiesForm}
 */
export const mapGeoPropertiesToForm = properties => ({ ...properties });

/**
 * Convert map flowchart properties object to form
 *
 * @param {MapFlowchartProperties} [properties = {}]
 * @returns {MapFlowchartPropertiesForm}
 */
export const mapFlowchartPropertiesToForm = properties => ({ ...properties });

/**
 * Convert map mermaid properties object to form
 *
 * @param {MapMermaidProperties} [properties = {}]
 * @returns {MapMermaidPropertiesForm}
 */
export const mapMermaidPropertiesToForm = (properties = {}) => ({
  theme: properties.theme ?? MERMAID_THEMES.default,
  code: properties.code ?? 'graph TB\na-->b',
});

/**
 * Convert map mermaid properties object to form
 *
 * @param {MapTreeOfDependenciesProperties} [properties = {}]
 * @returns {MapTreeOfDependenciesPropertiesForm}
 */
export const mapTreeOfDependenciesPropertiesToForm = (properties = {}) => ({
  ...properties,

  entities: addKeyInEntities(properties.entities),
});

/**
 * Convert map object to map form
 *
 * @param {Map} [map = {}]
 * @returns {MapForm}
 */
export const mapToForm = (map = {}) => {
  const type = map.type ?? MAP_TYPES.flowchart;

  const prepare = {
    [MAP_TYPES.geo]: mapGeoPropertiesToForm,
    [MAP_TYPES.flowchart]: mapFlowchartPropertiesToForm,
    [MAP_TYPES.mermaid]: mapMermaidPropertiesToForm,
    [MAP_TYPES.treeOfDependencies]: mapTreeOfDependenciesPropertiesToForm,
  }[type];

  return {
    type,
    name: map.name ?? '',
    properties: prepare(map.properties),
  };
};

/**
 * Convert form to tree of dependencies map
 *
 * @param {MapTreeOfDependenciesPropertiesForm} form
 * @returns {MapTreeOfDependenciesProperties}
 */
export const formToMapTreeOfDependenciesProperties = form => ({
  ...form,

  entities: removeKeyFromEntities(form.entities),
});

/**
 * Convert map form to map
 *
 * @param {MapForm} form
 * @returns {Map}
 */
export const formToMap = (form) => {
  const prepare = {
    [MAP_TYPES.treeOfDependencies]: formToMapTreeOfDependenciesProperties,
  }[form.type];

  return {
    ...form,

    properties: prepare ? prepare(form.properties) : form.properties,
  };
};
