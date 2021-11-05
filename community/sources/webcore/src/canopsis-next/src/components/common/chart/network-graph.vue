<template lang="pug">
  div.network-graph
    div.network-graph__canvas-wrapper(ref="canvasWrapper")
    div.v-tooltip__content.menuable__content__active.network-graph__tooltip(ref="tooltip")
</template>

<script>
import { isEqual, debounce } from 'lodash';

import { HEALTHCHECK_NETWORK_GRAPH_OPTIONS } from '@/constants';

import cytoscape from '@/services/cytoscape';

export default {
  props: {
    options: {
      type: Object,
      default: () => ({}),
    },
    nodeHtmlLabelOptions: {
      type: Array,
      default: () => [],
    },
    /**
     * @typedef {Object} NetworkGraphTooltipOptions
     * @property {number} offsetX
     * @property {number} offsetY
     * @property {function} getContent
     */
    tooltipOptions: {
      type: Object,
      default: () => ({}),
    },
  },
  watch: {
    nodeHtmlLabelOptions(newOptions, oldOptions) {
      if (!isEqual(newOptions, oldOptions)) {
        this.$cy.nodeHtmlLabel(newOptions);
      }
    },
  },
  created() {
    this.debouncedResizeHandler = debounce(this.resizeHandler, 50);
  },
  mounted() {
    this.mountCytoscape();

    window.addEventListener('resize', this.debouncedResizeHandler);
  },
  beforeDestroy() {
    this.destroyCytoscape();

    window.removeEventListener('resize', this.debouncedResizeHandler);
  },
  methods: {
    resizeHandler() {
      this.$cy.invalidateSize();
      this.$cy.fit(this.$cy.nodes(), HEALTHCHECK_NETWORK_GRAPH_OPTIONS.fitPadding);
    },

    mountCytoscape(options = this.options) {
      this.$cy = cytoscape({
        container: this.$refs.canvasWrapper,

        ...options,
      });

      this.$cy.nodeHtmlLabel(this.nodeHtmlLabelOptions);
      this.$cy.bind('mouseover', 'node', this.nodeMouseoverHandler);
      this.$cy.bind('mouseout', 'node', this.nodeMouseoutHandler);
      this.$cy.bind('tap', 'node', this.nodeTapHandler);
    },

    destroyCytoscape() {
      if (this.$cy) {
        this.$cy.destroy();
        this.$cy = null;
      }
    },

    remountCytoscape(options) {
      this.destroyCytoscape();
      this.mountCytoscape(options);
    },

    /**
     * Show tooltip with special position and content
     *
     * @param {{ x: number, y: number }} renderedPosition
     * @param {Object} data
     * @param {number} [zoom = 1]
     */
    showTooltip({ renderedPosition, data, zoom = 1 }) {
      const { getContent, offsetX = 0, offsetY = 0 } = this.tooltipOptions;

      if (!getContent) {
        return;
      }

      const tooltipInnerHTML = getContent(data);

      if (tooltipInnerHTML) {
        const { tooltip: tooltipEl } = this.$refs;
        const x = renderedPosition.x - (offsetX * zoom);
        const y = renderedPosition.y - (offsetY * zoom);

        tooltipEl.innerHTML = tooltipInnerHTML;
        tooltipEl.style.left = `${Math.round(x)}px`;
        tooltipEl.style.top = `${Math.round(y)}px`;
        tooltipEl.style.opacity = 1;
      }
    },

    /**
     * Hide tooltip
     */
    hideTooltip() {
      this.$refs.tooltip.style.opacity = 0;
    },

    /**
     * Handler for the 'tap' event on node for cytoscape
     *
     * @param {Object} event
     */
    nodeTapHandler(event) {
      this.$emit('node:tap', event);
    },

    /**
     * Handler for the 'mouseover' event on node for cytoscape
     *
     * @param {Object} event
     */
    nodeMouseoverHandler(event) {
      const { target, cy } = event;

      this.showTooltip({
        renderedPosition: target.renderedPosition(),
        data: target.data(),
        zoom: cy.zoom(),
      });

      this.$refs.canvasWrapper.classList.add('network-graph__canvas-wrapper--pointer');

      this.$emit('node:mouseover', event);
    },

    /**
     * Handler for the 'mouseout' event on node for cytoscape
     *
     * @param {Object} event
     */
    nodeMouseoutHandler(event) {
      this.hideTooltip();
      this.$refs.canvasWrapper.classList.remove('network-graph__canvas-wrapper--pointer');

      this.$emit('node:mouseout', event);
    },
  },
};
</script>

<style lang="scss" scoped>
.network-graph {
  position: relative;
  width: 100%;
  height: 100%;

  &__tooltip {
    transform: translate(-50%, -100%) !important;
  }

  &__canvas-wrapper {
    position: relative;
    width: inherit;
    height: inherit;
    background: white;
    cursor: grabbing;

    &--pointer {
      cursor: pointer;
    }
  }
}

.v-tooltip__content {
  position: absolute;
  top: 0;
  opacity: 0;
  transition: all .2s linear;
  transform: translate(-50%, 50%);
  pointer-events: none;
}
</style>