<template>
  <entity-network-graph
    ref="entityNetworkGraphElement"
    :options="options"
    :tooltip-options="tooltipOptions"
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
import { ref, computed, watch, onMounted } from 'vue';
import { omit } from 'lodash';

import { PAGINATION_LIMIT, VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  ROOT_CAUSE_DIAGRAM_OPTIONS,
  ROOT_CAUSE_DIAGRAM_TOOLTIP_OFFSET,
  ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
  COLOR_INDICATOR_TYPES,
  ENTITY_FIELDS,
  SORT_ORDERS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
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
    stateSetting: {
      type: Object,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
  },
  setup(props) {
    const { tc } = useI18n();

    const entityNetworkGraphElement = ref(null);
    const ready = ref(false);
    const pendingEntities = ref(true);

    const { fetchServiceDependenciesWithoutStore } = useService();

    const fetchHandler = async (id, page) => {
      const result = await fetchServiceDependenciesWithoutStore({
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

      return result;
    };

    const isEvent = computed(() => props.stateSetting && !props.stateSetting?.title);

    const {
      metaByEntityId,
      entitiesElements,
      toggleChildren,
      initRelations,
      showMore,
      resetEntities,
    } = useEntityNetworkGraph({
      entity: props.entity,
      withEvents: true,
      isEvent,
      fetchHandler,
    });

    const tooltipOptions = computed(() => ({
      offsetY: (ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize / 2) + ROOT_CAUSE_DIAGRAM_TOOLTIP_OFFSET,
      getContent: ({ isEvents, entity, root }) => {
        if (isEvents || root) {
          return '';
        }

        return entity.state_setting?.title || tc('common.event', 2);
      },
    }));

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
        selector: 'node[isEvents]',
        style: {
          width: 30,
          height: 30,
        },
      },
      {
        selector: 'edge',
        style: {
          width: 2,
          'curve-style': 'bezier',
          'line-color': 'silver',
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
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        };
      }

      return opts;
    });

    /**
     * Check if entity has children based on state settings and dependencies count
     *
     * @param {object} entity - Entity object
     * @returns {boolean} True if entity has children
     */
    const hasChildren = (entity = {}) => entity.state_setting?.title && entity.state_depends_count > 0;

    watch(() => props.entity, () => resetEntities(props.entity));

    onMounted(async () => {
      pendingEntities.value = true;

      if (!isEvent.value) {
        await initRelations(props.entity._id);
      } else {
        entityNetworkGraphElement.value.resetLayout();
      }

      /**
       * @desc: We are waiting modal showing animation
       */
      setTimeout(() => {
        entityNetworkGraphElement.value.fit();
        ready.value = true;
      }, VUETIFY_ANIMATION_DELAY);

      pendingEntities.value = false;
    });

    return {
      entityNetworkGraphElement,
      metaByEntityId,
      ready,
      isLoading,
      options,
      tooltipOptions,
      toggleChildren,
      showMore,
      hasChildren,
    };
  },
};
</script>
