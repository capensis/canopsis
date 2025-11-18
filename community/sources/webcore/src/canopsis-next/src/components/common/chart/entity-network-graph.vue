<template>
  <div class="entity-network-graph">
    <c-zoom-overlay>
      <c-progress-overlay :pending="!ready || loading" />
      <network-graph
        ref="networkGraphElement"
        :options="options"
        :tooltip-options="tooltipOptions"
        :node-html-label-options="nodeHtmlLabelsOptions"
        :class="{ 'entity-network-graph__canvas--ready': ready }"
        class="entity-network-graph__canvas fill-height black--text"
        ctrl-wheel-zoom
      />
    </c-zoom-overlay>
  </div>
</template>

<script>
import {
  computed,
  ref,
  nextTick,
  watch,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { COLOR_INDICATOR_TYPES_WITH_STATUS } from '@/constants';

import {
  getBadgeElement,
  getButtonHTML,
  getEntityNodeElement,
  getIconElement,
  getProgressElement,
} from '@/helpers/entities/entity/cytoscape';
import { getEntityColorClass } from '@/helpers/entities/entity/color';

import { useI18n } from '@/hooks/i18n';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

export default {
  components: { NetworkGraph },
  props: {
    options: {
      type: Object,
      default: () => ({}),
    },
    tooltipOptions: {
      type: Object,
      default: () => ({}),
    },
    ready: {
      type: Boolean,
      default: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    colorIndicator: {
      type: String,
      default: COLOR_INDICATOR_TYPES_WITH_STATUS.state,
    },
    metaByEntityId: {
      type: Object,
      default: () => ({}),
    },
    hasChildren: {
      type: Function,
      required: false,
    },
    onBadgeClick: {
      type: Function,
      required: false,
    },
    onShowMoreClick: {
      type: Function,
      required: false,
    },
  },
  setup(props) {
    const { t } = useI18n();

    const networkGraphElement = ref(null);

    /**
     * Generate HTML content for a network graph node
     *
     * @param {Object} node - Node object containing entity information
     * @param {Object} node.entity - Entity data
     * @param {boolean} node.pending - Whether the node is in pending state
     * @param {boolean} node.opened - Whether the node is expanded
     * @param {boolean} node.root - Whether the node is a root node
     * @returns {string} HTML string representing the node content
     */
    const getNodeContent = (node) => {
      const { entity, pending, opened, root } = node;

      const element = getEntityNodeElement(node);

      element.classList.add(getEntityColorClass(entity, props.colorIndicator));

      if (pending || (!root && props.hasChildren?.(entity))) {
        const badge = getBadgeElement();
        badge.dataset.id = entity._id;

        badge.appendChild(
          pending ? getProgressElement() : getIconElement(opened ? 'remove' : 'add', 'white'),
        );

        element.appendChild(badge);
      }

      return element.outerHTML;
    };

    /**
     * Generate HTML content for "Show More" button node
     *
     * @param {Object} node - Node object containing entity information
     * @param {Object} node.entity - Entity data with _id property
     * @returns {string} HTML string representing the "Show More" button
     */
    const getShowMoreButtonContent = (node) => {
      const { entity } = node;
      const meta = props.metaByEntityId[entity._id] ?? {};

      const fetchedEntities = meta.page * meta.per_page;

      return getButtonHTML(
        t('common.showMore', { current: fetchedEntities, total: meta.total_count }),
      );
    };

    const nodeHtmlLabelsOptions = computed(() => [
      {
        query: 'node',
        valign: 'center',
        halign: 'center',
        tpl: getNodeContent,
      },
      {
        query: 'node[showMore]',
        valign: 'center',
        halign: 'center',
        tpl: getShowMoreButtonContent,
      },
    ]);

    /**
     * Run layout algorithm on the network graph
     *
     * Executes the layout configuration on cytoscape instance without animation
     */
    const runLayout = async () => {
      if (networkGraphElement.value.$cy.nodes().empty()) {
        return;
      }

      try {
        await nextTick();

        networkGraphElement.value.$cy.layout({
          ...props.options.layout,
          animate: false,
        }).run();
      } catch (err) {
        console.warn(err);
      }
    };

    /**
     * Reset network graph layout
     *
     * Removes all existing elements from the graph and adds new elements from options, then runs the layout
     */
    const resetLayout = () => {
      if (!networkGraphElement.value?.$cy) {
        return;
      }

      networkGraphElement.value.$cy.elements().remove();
      networkGraphElement.value.$cy.add(props.options.elements);
      runLayout();
    };

    /**
     * Handle tap/click events on network graph nodes
     *
     * Processes clicks on badge elements or "show more" nodes and triggers appropriate callbacks
     *
     * @param {Object} event - Cytoscape tap event
     * @param {Object} event.target - The clicked cytoscape element
     * @param {Event} event.originalEvent - The original DOM event
     */
    const tapHandler = ({ target, originalEvent }) => {
      const { entity, showMore, pending, cycle } = target.data();

      if (cycle || pending) {
        return;
      }

      if (originalEvent.target.classList.contains('v-badge__badge')) {
        const { id } = originalEvent.target.dataset;

        if (id) {
          props.onBadgeClick?.(networkGraphElement.value.$cy.getElementById(id));

          return;
        }
      }

      if (!showMore || !entity) {
        return;
      }

      props.onShowMoreClick?.(networkGraphElement.value.$cy.getElementById(entity._id));
    };

    const fit = (padding = 0) => networkGraphElement.value.$cy.fit(padding);

    watch(() => props.options.elements, () => resetLayout());

    onMounted(() => networkGraphElement.value.$cy.on('tap', tapHandler));
    onBeforeUnmount(() => networkGraphElement.value.$cy.off('tap', tapHandler));

    return {
      nodeHtmlLabelsOptions,
      networkGraphElement,

      /**
       * Expose cytoscape functionality for external use
       */
      resetLayout,
      runLayout,
      fit,
    };
  },
};
</script>

<style lang="scss">
.entity-network-graph {
  position: relative;
  height: 650px;
  width: 100%;
  border-radius: 5px;
  background: white;

  &__canvas {
    opacity: 0;

    &--ready {
      opacity: 1;
    }
  }

  canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  .v-badge__badge {
    top: -7px;
    right: -7px;

    * {
      pointer-events: none;
    }
  }
}
</style>
