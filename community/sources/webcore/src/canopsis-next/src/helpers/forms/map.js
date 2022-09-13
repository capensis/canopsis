import { MAP_TYPES, MERMAID_THEMES, TREE_OF_DEPENDENCIES_TYPES } from '@/constants';

import uuid from '@/helpers/uuid';
import { addKeyInEntities, mapIds, removeKeyFromEntities } from '@/helpers/entities';

/**
 * @typedef {Object} MapCommonFields
 * @property {string} name
 * @property {string} [_id]
 */

/**
 * @typedef {Object} MapMermaidPoint
 * @property {string} _id
 * @property {MapCommonFields} map
 * @property {Entity} entity
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {MapMermaidPoint} MapMermaidPointForm
 * @property {string} map
 * @property {string} entity
 */

/**
 * @typedef {MapCommonFields} MapMermaidParameters
 * @property {string} theme
 * @property {string} code
 * @property {MapMermaidPoint[]} points
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
 * @property {MapGeoPointForm[]} points
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
 * @typedef {Object} MapTreeOfDependenciesEntity
 * @property {Entity} data
 * @property {Entity[]} pinned
 */

/**
 * @typedef {MapTreeOfDependenciesEntity} MapTreeOfDependenciesEntityForm
 * @property {string} key
 */

/**
 * @typedef { 'treeofdeps' | 'impactchain' } MapTreeOfDependenciesParametersType
 */

/**
 * @typedef {Object} MapTreeOfDependenciesParameters
 * @property {MapTreeOfDependenciesEntity[]} entities
 * @property {MapTreeOfDependenciesParametersType} type
 */

/**
 * @typedef {MapTreeOfDependenciesParameters} MapTreeOfDependenciesParametersForm
 * @property {MapTreeOfDependenciesEntityForm[]} entities
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
 * Convert mermaid point to form object
 *
 * @param {MapMermaidPoint} [point = {}]
 * @returns {MapMermaidPointForm}
 */
export const mermaidPointToForm = (point = {}) => ({
  x: point.x,
  y: point.y,
  entity: point.entity?._id ?? '',
  map: point.map?._id,
  _id: uuid(),
});

/**
 * Convert mermaid point to form object
 *
 * @param {MapMermaidPoint[]} [points = []]
 * @returns {MapMermaidPointForm[]}
 */
export const mermaidPointsToForm = (points = []) => points.map(mermaidPointToForm);

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
  entity: point.entity?._id ?? '',
  map: point.map?._id,
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
  code: parameters.code ?? 'graph TB\n  a-->b',
  points: mermaidPointsToForm(parameters.points),
});

/**
 * Convert map mermaid parameters object to form
 *
 * @param {MapTreeOfDependenciesParameters} [parameters = {}]
 * @returns {MapTreeOfDependenciesParametersForm}
 */
export const mapTreeOfDependenciesParametersToForm = (parameters = {}) => ({
  ...parameters,

  impact: parameters.type === TREE_OF_DEPENDENCIES_TYPES.impactChain,
  entities: addKeyInEntities(parameters.entities ?? []).map(({ entity, pinned_entities: pinned }) => ({
    data: entity,
    pinned,
  })),
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
    [MAP_TYPES.geo]: mapGeoParametersToForm,
    [MAP_TYPES.flowchart]: mapFlowchartParametersToForm,
    [MAP_TYPES.mermaid]: mapMermaidParametersToForm,
    [MAP_TYPES.treeOfDependencies]: mapTreeOfDependenciesParametersToForm,
  }[type];

  return {
    type,
    name: map.name ?? '',
    parameters: prepare(map.parameters),
  };
};

/**
 * Convert form to tree of dependencies map
 *
 * @param {boolean} impact
 * @param {MapTreeOfDependenciesParametersForm} form
 * @returns {MapTreeOfDependenciesParameters}
 */
export const formToMapTreeOfDependenciesParameters = ({ impact, ...form }) => ({
  ...form,

  type: impact ? TREE_OF_DEPENDENCIES_TYPES.impactChain : TREE_OF_DEPENDENCIES_TYPES.treeOfDependencies,
  entities: removeKeyFromEntities(form.entities).map(({ data, pinned }) => ({
    entity: data._id,
    pinned_entities: mapIds(pinned),
  })),
});

/**
 * Convert map form to map
 *
 * @param {MapForm} form
 * @returns {Map}
 */
export const formToMap = (form) => {
  const prepare = {
    [MAP_TYPES.treeOfDependencies]: formToMapTreeOfDependenciesParameters,
  }[form.type];

  return {
    ...form,

    parameters: prepare ? prepare(form.parameters) : form.parameters,
  };
};
