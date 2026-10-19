/**
 * Get show-more node ID for an entity
 *
 * @param {string} entityId - The entity ID
 * @returns {string} The show-more node ID
 */
export const getShowMoreNodeId = entityId => `show-more-${entityId}`;

/**
 * Get show-more edge ID for an entity
 *
 * @param {string} entityId - The entity ID
 * @returns {string} The show-more edge ID
 */
export const getShowMoreEdgeId = entityId => `show-more-edge-${entityId}`;

/**
 * Get show-all node ID for an entity
 *
 * @param {string} entityId - The entity ID
 * @returns {string} The show-all node ID
 */
export const getShowAllNodeId = entityId => `show-all-${entityId}`;

/**
 * Get show-all edge ID for an entity
 *
 * @param {string} entityId - The entity ID
 * @returns {string} The show-all edge ID
 */
export const getShowAllEdgeId = entityId => `show-all-edge-${entityId}`;

/**
 * Create cytoscape node element
 *
 * @param {string} id - The node ID
 * @param {Object} entity - The entity object
 * @param {boolean} [opened=false] - Whether the node is opened
 * @returns {Object} Cytoscape node element
 */
export const createNodeElement = (id, entity, opened = false) => ({
  group: 'nodes',
  data: {
    id,
    entity,
    opened,
  },
});

/**
 * Create cytoscape edge element
 *
 * @param {string} sourceId - The source node ID
 * @param {string} targetId - The target node ID
 * @returns {Object} Cytoscape edge element
 */
export const createEdgeElement = (sourceId, targetId) => ({
  group: 'edges',
  data: {
    source: sourceId,
    target: targetId,
  },
});
