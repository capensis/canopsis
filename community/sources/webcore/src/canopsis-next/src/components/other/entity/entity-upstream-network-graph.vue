<template>
  <entity-network-graph
    ref="entityNetworkGraphElement"
    :options="options"
    :ready="ready"
    :loading="isLoading"
    :color-indicator="colorIndicator"
    :meta-by-entity-id="metaByEntityId"
    :has-children="hasChildren"
    :on-badge-click="toggleChildren"
    :on-show-more-click="showMore"
  />
</template>

<script>
import { omit, isEmpty } from 'lodash';
import { ref, computed, watch, onMounted } from 'vue';

import { PAGINATION_LIMIT, VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  ROOT_CAUSE_DIAGRAM_OPTIONS,
  ENTITY_UPSTREAM_GRAPH_LAYOUT_OPTIONS,
  COLOR_INDICATOR_TYPES_WITH_STATUS,
  ENTITY_FIELDS,
  SORT_ORDERS,
} from '@/constants';

import { useEntityNetworkGraph } from '@/hooks/charts/entity-network-graph';
import { useService } from '@/hooks/store/modules/service';

import EntityNetworkGraph from '@/components/common/chart/entity-network-graph.vue';

export default {
  components: { EntityNetworkGraph },
  props: {
    entity: {
      type: Object,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const colorIndicator = COLOR_INDICATOR_TYPES_WITH_STATUS.status;

    const entityNetworkGraphElement = ref(null);
    const ready = ref(false);
    const pendingEntities = ref(true);

    const { fetchEntityUpstreamWithoutStore, fetchEntityDownstreamsWithoutStore } = useService();

    const fetchParent = async (id) => {
      const parent = await fetchEntityUpstreamWithoutStore({ id });

      return {
        data: isEmpty(parent) ? [] : [{ ...parent, isParent: true }],
      };
    };

    const fetchChildren = async (id, page) => fetchEntityDownstreamsWithoutStore({
      id,
      params: {
        page,
        limit: PAGINATION_LIMIT,
        with_flags: true,
        define_state: true,
        sort_by: ENTITY_FIELDS.impactState,
        sort: SORT_ORDERS.desc.toLowerCase(),
      },
    });

    const {
      metaByEntityId,
      toggleChildren,
      initRelations,
      showMore,
      resetEntities,

      entitiesElements: rawEntitiesElements,
    } = useEntityNetworkGraph({
      entity: props.entity,
      fetchHandler: async (id, page, entity) => {
        let parentPromise;
        let childrenPromise;

        if (id === props.entity._id) {
          if (page === 1) {
            parentPromise = fetchParent(id);
          }

          childrenPromise = fetchChildren(id, page);
        } else if (entity.isParent) {
          parentPromise = fetchParent(id);
        } else {
          childrenPromise = fetchChildren(id, page);
        }

        const [
          { data: parentData } = { data: [] },
          { data: childrenData, meta: childrenMeta } = { data: [], meta: {} },
        ] = await Promise.all([parentPromise, childrenPromise]);

        return {
          data: [...parentData, ...childrenData],
          meta: childrenMeta,
        };
      },
    });

    const entitiesElements = computed(() => {
      const elements = rawEntitiesElements.value;

      return elements.map((element) => {
        if (element.group === 'edges' && element.data) {
          const sourceEntity = elements.find(e => e.group === 'nodes' && e.data?.id === element.data.source);
          const targetEntity = elements.find(e => e.group === 'nodes' && e.data?.id === element.data.target);
          const { isParent: sourceIsParent } = sourceEntity?.data?.entity ?? {};
          const { isParent: targetIsParent } = targetEntity?.data?.entity ?? {};

          if ((sourceIsParent && !targetIsParent) || (sourceIsParent && targetIsParent)) {
            return {
              ...element,
              data: {
                ...element.data,
                source: element.data.target,
                target: element.data.source,
              },
            };
          }
        }

        return element;
      });
    });

    const isLoading = computed(() => props.pending || pendingEntities.value);

    const styleOption = computed(() => [
      {
        selector: 'node',
        style: {
          width: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
          height: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
        },
      },
      {
        selector: 'node[showMore]',
        style: {
          'background-opacity': 0,
          'border-width': 0,
          width: 128,
          height: 34,
        },
      },
      {
        selector: 'edge',
        style: {
          width: 2,
          'curve-style': 'bezier',
          'line-color': 'silver',
          'target-arrow-shape': 'vee',
          'target-arrow-color': 'silver',
          'arrow-scale': 1.5,
          'target-distance-from-node': 30,
        },
      },
    ]);

    const options = computed(() => {
      const opts = {
        ...omit(ROOT_CAUSE_DIAGRAM_OPTIONS, ['nodeSize']),

        style: styleOption.value,
        elements: entitiesElements.value,
      };

      if (entitiesElements.value.length) {
        opts.layout = {
          ...ENTITY_UPSTREAM_GRAPH_LAYOUT_OPTIONS,
        };
      }

      return opts;
    });

    /**
     * Check if entity has children based on upstream/downstream count
     *
     * @param {object} entity - Entity object
     * @returns {boolean} True if entity has children
     */
    const hasChildren = (entity = {}) => (entity.isParent && entity.upstream) || entity.downstream_count > 0;

    watch(() => props.entity, entity => resetEntities(entity));

    onMounted(async () => {
      pendingEntities.value = true;

      await initRelations(props.entity._id);

      /**
       * @desc: We are waiting modal showing animation
       */
      setTimeout(() => {
        entityNetworkGraphElement.value.fit(20);
        ready.value = true;
      }, VUETIFY_ANIMATION_DELAY);

      pendingEntities.value = false;
    });

    return {
      colorIndicator,
      entityNetworkGraphElement,
      metaByEntityId,
      ready,
      isLoading,
      options,
      toggleChildren,
      showMore,
      hasChildren,
    };
  },
};
</script>
