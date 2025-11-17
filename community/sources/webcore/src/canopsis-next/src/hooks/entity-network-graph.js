import { isEmpty } from 'lodash';
import { ref, computed, unref } from 'vue';

import { normalizeTreeOfDependenciesMapEntities, getEntityChildrenElements } from '@/helpers/entities/map/list';

/**
 * Composable for managing entity network graph data and operations
 *
 * @param {object} config - Configuration object
 * @param {object|Ref<object>} config.entity - Root entity object or ref
 * @param {string|Ref<string>} [config.childrenKey='dependencies'] - Key for children property
 * @param {boolean|Ref<boolean>} [config.isEvent=false] - Whether the entity is an event
 * @param {boolean|Ref<boolean>} [config.withEvents=false] - Whether to include events nodes
 * @param {Function} config.fetchHandler - Handler function to fetch entity data
 * @returns {object} Object containing reactive state and methods for managing the network graph
 */
export const useEntityNetworkGraph = ({
  entity: rootEntity,
  childrenKey = 'dependencies',
  isEvent = false,
  withEvents = false,
  fetchHandler,
}) => {
  if (!fetchHandler) {
    throw new Error('fetchHandler is required');
  }
  const metaByEntityId = ref({});
  const entitiesById = ref(
    normalizeTreeOfDependenciesMapEntities([{ entity: unref(rootEntity), pinned_entities: [] }]),
  );

  /**
   * Generate events node element and its edge for an entity
   *
   * @param {object} entity - Entity object
   * @returns {Array} Array containing node and edge elements for events
   */
  const getEventsNodeElementByEntity = (entity) => {
    const eventsNodeId = `${entity._id}_events-node`;

    return [
      {
        group: 'nodes',
        data: {
          entity,
          id: eventsNodeId,
          isEvents: true,
        },
      },
      {
        group: 'edges',
        data: {
          source: eventsNodeId,
          target: entity._id,
        },
      },
    ];
  };

  /**
   * Generate "show more" node element and its edge for an entity
   *
   * @param {object} entity - Entity object
   * @returns {Array} Array containing show more node and edge elements, or empty array if no more pages
   */
  const getShowMoreElements = (entity) => {
    const meta = metaByEntityId.value[entity._id];

    if (isEmpty(meta) || meta.page >= meta.page_count) {
      return [];
    }

    const showMoreId = `show-all-${entity._id}`;

    return [
      {
        group: 'nodes',
        data: {
          id: showMoreId,
          entity,
          showMore: true,
        },
      },
      {
        group: 'edges',
        data: {
          id: `show-all-edge-${entity._id}`,
          source: showMoreId,
          target: entity._id,
        },
      },
    ];
  };

  /**
   * Update entity data in the entities map
   *
   * @param {string} entityId - Entity ID
   * @param {object} entity - Entity object to update or add
   */
  const updateEntityData = (entityId, entity) => {
    const existingEntity = entitiesById.value[entityId];

    entitiesById.value = {
      ...entitiesById.value,
      [entityId]: existingEntity
        ? {
          ...existingEntity,
          entity: {
            ...existingEntity.entity,
            ...entity,
          },
        }
        : { entity },
    };
  };

  /**
   * Set pending state on a target node
   *
   * @param {object} target - Target node object
   * @param {boolean} pending - Pending state
   */
  const setTargetPending = (target, pending) => target.data?.({ pending });

  /**
   * Computed property that builds the graph elements structure from entities
   *
   * @returns {Array} Array of graph elements (nodes and edges)
   */
  const entitiesElements = computed(() => {
    const unwrappedRootEntity = unref(rootEntity);
    const unwrappedChildrenKey = unref(childrenKey);
    const unwrappedIsEvent = unref(isEvent);
    const unwrappedWithEvents = unref(withEvents);
    const rootElement = entitiesById.value[unwrappedRootEntity._id];
    const { entity, [unwrappedChildrenKey]: children = [] } = rootElement;

    const elements = [
      {
        group: 'nodes',
        data: {
          id: entity._id,
          entity,
          root: true,
          opened: true,
        },
      },
    ];

    if (unwrappedIsEvent) {
      elements.push(...getEventsNodeElementByEntity(entity));

      return elements;
    }

    elements.push(...getEntityChildrenElements({
      entity,
      childrenIds: children,
      handledChildrenIds: [entity._id],
      entitiesById: entitiesById.value,
      childrenKey: unwrappedChildrenKey,
      withEvents: unwrappedWithEvents,
      getEventsNodeElementByEntity,
      getShowMoreElements,
    }));

    return elements;
  });

  /**
   * Show or load more children for a target entity node
   *
   * @param {object|string} target - Target node object or entity ID
   */
  const showChildren = async (target) => {
    const { id } = target.data?.() ?? { id: target };
    const currentPage = metaByEntityId.value[id]?.page ?? 0;
    const nextPage = currentPage + 1;

    setTargetPending(target, true);

    try {
      const { data, meta } = await fetchHandler(id, nextPage, entitiesById.value[id]?.entity);

      metaByEntityId.value = {
        ...metaByEntityId.value,
        [id]: meta,
      };

      if (!data.length) {
        return;
      }

      const childrenIds = data.map((entity) => {
        updateEntityData(entity._id, entity);

        return entity._id;
      });

      const unwrappedChildrenKey = unref(childrenKey);
      const currentChildren = entitiesById.value[id]?.[unwrappedChildrenKey] ?? [];

      entitiesById.value = {
        ...entitiesById.value,
        [id]: {
          ...entitiesById.value[id],
          [unwrappedChildrenKey]: [...currentChildren, ...childrenIds],
        },
      };
    } finally {
      setTargetPending(target, false);
    }
  };

  /**
   * Hide all children of a target entity node
   *
   * @param {object} target - Target node object
   */
  const hideChildren = (target) => {
    const { entity } = target.data();
    const unwrappedChildrenKey = unref(childrenKey);

    entitiesById.value = {
      ...entitiesById.value,
      [entity._id]: {
        ...entitiesById.value[entity._id],
        [unwrappedChildrenKey]: [],
      },
    };

    const { [entity._id]: removed, ...rest } = metaByEntityId.value;
    metaByEntityId.value = rest;
  };

  /**
   * Toggle children visibility for a target entity node
   *
   * @param {object} target - Target node object
   */
  const toggleChildren = async (target) => {
    const { opened, root } = target.data();

    if (!root && opened) {
      return hideChildren(target);
    }

    return showChildren(target);
  };

  /**
   * Initialize relations for an entity by loading its children
   *
   * @param {string} entityId - Entity ID to initialize relations for
   */
  const initRelations = entityId => showChildren(entityId);

  /**
   * Show more items for a target node
   *
   * @param {object} target - Target node object
   */
  const showMore = target => showChildren(target);

  /**
   * Reset the entities state with a new root entity
   *
   * @param {object} newEntity - New root entity object
   */
  const resetEntities = newEntity => (
    entitiesById.value = normalizeTreeOfDependenciesMapEntities([{ entity: newEntity, pinned_entities: [] }])
  );

  return {
    metaByEntityId,
    entitiesById,
    entitiesElements,
    toggleChildren,
    initRelations,
    showMore,
    resetEntities,
  };
};
