import { keyBy, omit } from 'lodash';

import { COLORS } from '@/config';
import { MAP_TYPES, MERMAID_THEMES, TREE_OF_DEPENDENCIES_TYPES } from '@/constants';

import uuid from '@/helpers/uuid';
import { shapeToForm } from '@/helpers/flowchart/shapes';
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
 * @property {boolean} is_entity_coordinates
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
 * @typedef {Object} MapFlowchartPoint
 * @property {string} _id
 * @property {string} [shape]
 * @property {MapCommonFields} map
 * @property {Entity} entity
 * @property {number} [x]
 * @property {number} [y]
 */

/**
 * @typedef {Object} MapFlowchartParameters
 * @property {Shape[]} shapes
 * @property {string} background_color
 * @property {MapFlowchartPoint[]} points
 */

/**
 * @typedef {MapFlowchartPoint} MapFlowchartPointForm
 * @property {string} map
 * @property {string} entity
 */

/**
 * @typedef {MapFlowchartParameters} MapFlowchartParametersForm
 * @property {MapFlowchartPointForm[]} points
 */

/**
 * @typedef {MapCommonFields} MapFlowchart
 * @property {'flowchart'} type
 * @property {MapFlowchartParameters} parameters
 */

/**
 * @typedef {Object} MapTreeOfDependenciesEntity
 * @property {Entity} entity
 * @property {Entity[]} pinned_entities
 */

/**
 * @typedef {Object} MapTreeOfDependenciesEntityRequest
 * @property {string} entity
 * @property {string[]} pinned_entities
 */

/**
 * @typedef {MapTreeOfDependenciesEntity} MapTreeOfDependenciesEntityForm
 * @property {string} key
 * @property {Entity} entity
 * @property {Entity[]} pinned
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
  _id: point._id ?? uuid(),
});

/**
 * Convert mermaid point to form object
 *
 * @param {MapFlowchartPoint} [point = {}]
 * @returns {MapFlowchartPointForm}
 */
export const flowchartPointToForm = (point = {}) => ({
  x: point.x,
  y: point.y,
  entity: point.entity?._id ?? '',
  shape: point.shape ?? '',
  map: point.map?._id,
  _id: point._id ?? uuid(),
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
 * @returns {MapGeoPointForm}
 */
export const geomapPointToForm = (point = {}) => ({
  coordinates: point.coordinates ?? {
    lat: 0,
    lng: 0,
  },
  is_entity_coordinates: !!point.entity?.coordinates,
  entity: point.entity?._id ?? '',
  map: point.map?._id,
  _id: point._id ?? uuid(),
});

/**
 * Convert geomap points to form array
 *
 * @param {MapGeoPoint[]} [points = {}]
 * @returns {MapGeoPointForm[]}
 */
export const geomapPointsToForm = (points = []) => points.map(geomapPointToForm);

/**
 * Convert flowchart points to form array
 *
 * @param {MapFlowchartPoint[]} [points = {}]
 * @returns {MapFlowchartPointForm[]}
 */
export const flowchartPointsToForm = (points = []) => points.map(flowchartPointToForm);

/**
 * Convert flowchart shapes to form
 *
 * @param {Shape[]} [shapes = []]
 * @returns {Object.<string, Shape>}
 */
export const flowchartShapesToForm = (shapes = []) => keyBy(shapes.map(shapeToForm), '_id');

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
export const mapFlowchartParametersToForm = (parameters = {}) => ({
  shapes: parameters.shapes ? flowchartShapesToForm(parameters.shapes) : {},
  background_color: parameters.background_color ?? COLORS.flowchart.background[0],
  points: flowchartPointsToForm(parameters.points),
});

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
 * Convert entities array of tree of dependencies to form
 *
 * @param {MapTreeOfDependenciesEntity[] | [{}]} entities
 * @returns {MapTreeOfDependenciesEntityForm[]}
 */
export const mapTreeOfDependenciesParametersEntitiesToForm = (entities = [{}]) => (
  addKeyInEntities(entities.map(({ entity, pinned_entities: pinned = [] }) => ({
    entity,
    pinned,
  })))
);

/**
 * Convert map mermaid parameters object to form
 *
 * @param {MapTreeOfDependenciesParameters} [parameters = {}]
 * @returns {MapTreeOfDependenciesParametersForm}
 */
export const mapTreeOfDependenciesParametersToForm = (parameters = {}) => ({
  ...parameters,

  impact: parameters.type === TREE_OF_DEPENDENCIES_TYPES.impactChain,
  entities: mapTreeOfDependenciesParametersEntitiesToForm(parameters.entities),
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
 * Convert entities array of tree of dependencies to form
 *
 * @param {MapTreeOfDependenciesEntityForm[]} entities
 * @returns {MapTreeOfDependenciesEntityRequest[]}
 */
export const formToMapTreeOfDependenciesParametersEntities = (entities = []) => (
  removeKeyFromEntities(entities).map(({ entity, pinned }) => ({
    entity: entity._id,
    pinned_entities: mapIds(pinned),
  }))
);

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
  entities: formToMapTreeOfDependenciesParametersEntities(form.entities),
});

/**
 * Convert form parameters to flowchart map parameters
 *
 * @param {MapFlowchartPointForm[]} points
 * @returns {MapFlowchartPoint[]}
 */
export const formPointsToMapFlowchartPoints = points => points.map(
  point => omit(
    point,
    point.shape ? ['x', 'y'] : ['shape'],
  ),
);

/**
 * Convert form parameters to flowchart map parameters
 *
 * @param {MapFlowchartParametersForm} form
 * @returns {MapFlowchartParameters}
 */
export const formToMapFlowchartParameters = form => ({
  ...form,

  points: formPointsToMapFlowchartPoints(form.points),
  shapes: Object.values(form.shapes),
});

/**
 * Convert form parameters to geomap parameters
 *
 * @param {MapGeoParametersForm} form
 * @returns {MapGeoParameters}
 */
export const formToMapGeomapParameters = form => ({
  ...form,

  points: form.points.map(point => omit(point, ['is_entity_coordinates'])),
});

/**
 * Convert map form to map
 *
 * @param {MapForm} form
 * @returns {Map}
 */
export const formToMap = (form) => {
  const prepare = {
    [MAP_TYPES.geo]: formToMapGeomapParameters,
    [MAP_TYPES.flowchart]: formToMapFlowchartParameters,
    [MAP_TYPES.treeOfDependencies]: formToMapTreeOfDependenciesParameters,
  }[form.type];

  return {
    ...form,

    parameters: prepare ? prepare(form.parameters) : form.parameters,
  };
};
