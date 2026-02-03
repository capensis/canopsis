import { cloneDeep } from 'lodash';

import { ENTITY_TYPES } from '@/constants';

import { createNodeElement, createEdgeElement } from '@/helpers/cytoscape/elements';

/**
 * Get icon data by geo map point
 *
 * @param {MapGeoPoint} point
 * @param {number} size
 * @returns {Object}
 */
export const getGeomapMarkerIconOptions = (point, size) => {
  const halfIconSize = size / 2;
  const pixelSize = `${size}px`;

  return {
    style: {
      width: pixelSize,
      height: pixelSize,
      maxWidth: 'unset',
      maxHeight: 'unset',
    },
    size,
    anchor: point.entity
      ? [halfIconSize, size]
      : [halfIconSize, halfIconSize],
  };
};

/**
 * Get text for displaying in map components for entity
 *
 * @param {Entity} entity
 * @returns {string}
 */
export const getMapEntityText = entity => (
  entity.type === ENTITY_TYPES.resource
    ? entity._id
    : entity.name
);

/**
 * @typedef {Object} TreeOfDependenciesMapEntity
 * @property {Entity} entity
 * @property {string[]} [dependencies]
 */

/**
 *
 * @param {Entity[]} [entities = []]
 * @returns {Object<string, TreeOfDependenciesMapEntity>}
 */
export const normalizeTreeOfDependenciesMapEntities = (entities = [], childrenKey = 'dependencies') => (
  entities.reduce((acc, { entity, pinned_entities: pinnedEntities }) => {
    const newEntity = {
      entity: cloneDeep(entity),
      [childrenKey]: [],
    };

    pinnedEntities.forEach((pinnedEntity) => {
      const { _id: id } = pinnedEntity;

      newEntity[childrenKey].push(id);

      if (!acc[id]) {
        acc[id] = {
          entity: cloneDeep(pinnedEntity),
        };
      }
    });

    acc[entity._id] = newEntity;

    return acc;
  }, {})
);

/**
 * Get entity children elements for cytoscape
 *
 * @param {Object} params
 * @param {Entity} params.entity - The entity object
 * @param {string[]} [params.childrenIds=[]] - Array of children IDs
 * @param {string[]} [params.handledChildrenIds=[]] - Array of already handled children IDs to prevent cycles
 * @param {Object} params.entitiesById - Map of entities by ID
 * @param {string} [params.childrenKey='dependencies'] - Key name for children in entitiesById
 * @param {Function} params.getEventsNodeElementByEntity - Function to get events node element
 * @param {Function} params.getShowMoreElements - Function to get show more elements
 * @returns {Array} Array of cytoscape elements
 */
export const getEntityChildrenElements = ({
  entity,
  childrenIds = [],
  handledChildrenIds = [],
  entitiesById,
  withEvents = false,
  childrenKey = 'dependencies',
  getEventsNodeElementByEntity,
  getShowMoreElements,
}) => {
  const getChildElements = (childId) => {
    const childData = entitiesById[childId];
    if (!childData) {
      return [];
    }

    const childChildrenIds = childData[childrenKey] || [];
    const { entity: child } = childData;
    const isCycle = handledChildrenIds.includes(childId);
    const hasChildren = !!childChildrenIds.length;

    const elements = [];

    if (!isCycle) {
      const childChildrenElements = getEntityChildrenElements({
        entity: child,
        childrenIds: childChildrenIds,
        handledChildrenIds: [...handledChildrenIds, childId],
        entitiesById,
        childrenKey,
        getEventsNodeElementByEntity,
        getShowMoreElements,
      });

      elements.push(
        createNodeElement(childId, child, hasChildren),
        ...childChildrenElements,
      );
    }

    elements.push(createEdgeElement(childId, entity._id));

    if (withEvents && !child.state_setting?.title) {
      elements.push(...getEventsNodeElementByEntity(child));
    }

    return elements;
  };

  const childrenElements = childrenIds.flatMap(getChildElements);

  childrenElements.push(...getShowMoreElements(entity));

  return childrenElements;
};
