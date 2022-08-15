import { MAP_TYPES, MERMAID_THEMES } from '@/constants';

import uuid from '@/helpers/uuid';

/**
 * @typedef {Object} MapCommonFields
 * @property {string} name
 * @property {string} [_id]
 */

/**
 * @typedef {MapCommonFields} MapMermaidParameters
 * @property {string} theme
 * @property {string} code
 */

/**
 * @typedef {MapMermaidParameters} MapMermaidParametersForm
 */

/**
 * @typedef {MapCommonFields} MapMermaid
 * @property {'mermaid'} type
 * @property {MapMermaidParameters} parameters
 */

/**
 * @typedef {Object} MapGeoPoint
 * @property {string} _id
 * @property {MapCommonFields} map
 * @property {Entity} entity
 * @property {Object} coordinates
 */

/**
 * @typedef {MapGeoPoint} MapGeoPointForm
 * @property {string} map
 * @property {string} entity
 */

/**
 * @typedef {Object} MapGeoParameters
 * @property {MapGeoPoint[]} points
 */

/**
 * @typedef {MapGeoParameters} MapGeoParametersForm
 */

/**
 * @typedef {MapCommonFields} MapGeo
 * @property {'geo'} type
 * @property {MapGeoParameters} parameters
 */

/**
 * @typedef {Object} MapFlowchartParameters
 */

/**
 * @typedef {MapFlowchartParameters} MapFlowchartParametersForm
 */

/**
 * @typedef {MapCommonFields} MapFlowchart
 * @property {'flowchart'} type
 * @property {MapFlowchartParameters} parameters
 */

/**
 * @typedef {Object} MapTreeOfDependenciesParameters
 */

/**
 * @typedef {MapTreeOfDependenciesParameters} MapTreeOfDependenciesParametersForm
 */

/**
 * @typedef {MapCommonFields} MapTreeOfDependencies
 * @property {'treeOfDependencies'} type
 * @property {MapTreeOfDependenciesParameters} parameters
 */

/**
 * @typedef {MapMermaid} Map
 */

/**
 * @typedef {MapMermaid} MapForm
 */

/**
 * Convert geomap point to form object
 *
 * @param {MapGeoPoint} [point = {}]
 * @returns {MapGeoPoint}
 */
export const geomapPointToForm = (point = {}) => ({
  coordinates: point.coordinates ?? {
    lat: 0,
    lng: 0,
  },
  entity: point.entity ?? '',
  map: point.map,
  _id: uuid(),
});

/**
 * Convert geomap points to form array
 *
 * @param {MapGeoPoint[]} [points = {}]
 * @returns {MapGeoPointForm[]}
 */
export const geomapPointsToForm = (points = []) => points.map(geomapPointToForm);

/**
 * Convert map geo parameters object to form
 *
 * @param {MapGeoParameters} [parameters = {}]
 * @returns {MapGeoParametersForm}
 */
export const mapGeoParametersToForm = (parameters = {}) => ({
  points: geomapPointsToForm(parameters.points),
});

/**
 * Convert map flowchart parameters object to form
 *
 * @param {MapFlowchartParameters} [parameters = {}]
 * @returns {MapFlowchartParametersForm}
 */
export const mapFlowchartParametersToForm = parameters => ({ ...parameters });

/**
 * Convert map mermaid parameters object to form
 *
 * @param {MapMermaidParameters} [parameters = {}]
 * @returns {MapMermaidParametersForm}
 */
export const mapMermaidParametersToForm = (parameters = {}) => ({
  theme: parameters.theme ?? MERMAID_THEMES.default,
  code: parameters.code ?? 'graph TB\na-->b',
});

/**
 * Convert map mermaid parameters object to form
 *
 * @param {MapTreeOfDependenciesParameters} [parameters = {}]
 * @returns {MapTreeOfDependenciesParametersForm}
 */
export const mapTreeOfDependenciesParametersToForm = parameters => ({ ...parameters });

/**
 * Convert map object to map form
 *
 * @param {Map} [map = {}]
 * @returns {MapForm}
 */
export const mapToForm = (map = {}) => {
  const type = map.type ?? MAP_TYPES.flowchart;

  const prepare = {
    [MAP_TYPES.geo]: mapGeoParametersToForm,
    [MAP_TYPES.flowchart]: mapFlowchartParametersToForm,
    [MAP_TYPES.mermaid]: mapMermaidParametersToForm,
    [MAP_TYPES.treeOfDependencies]: mapTreeOfDependenciesParametersToForm,
  }[type];

  return {
    name: map.name ?? '',
    type,
    parameters: prepare(map.parameters),
  };
};

/**
 * Convert map form to map
 *
 * @param {MapForm} form
 * @returns {Map}
 */
export const formToMap = form => ({ ...form });
