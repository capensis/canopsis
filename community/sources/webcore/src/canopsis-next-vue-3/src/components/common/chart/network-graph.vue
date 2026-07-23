<template>
  <div class="network-graph">
    <div ref="canvasWrapper" class="network-graph__canvas-wrapper" />
    <div ref="tooltip" class="v-tooltip__content network-graph__tooltip" />
  </div>
</template>

<script>
import { debounce, isEqual } from 'lodash';

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
    tooltipOptions: {
      type: Object,
      default: () => ({}),
    },
    ctrlWheelZoom: {
      type: Boolean,
      default: false,
    },
  },

  watch: {
    nodeHtmlLabelOptions(newOptions, oldOptions) {
      if (!isEqual(newOptions, oldOptions) && this.$cy) {
        this.$cy.nodeHtmlLabel(newOptions, {
          enablePointerEvents: true,
        });
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

  beforeUnmount() {
    this.destroyCytoscape();

    window.removeEventListener('resize', this.debouncedResizeHandler);
  },

  methods: {
    wheelHandler(event) {
      const renderer = this.$cy.renderer();
      const inBoxSelection = () => renderer.selection[4] !== 0;

      if (renderer.scrollingPage || !event.ctrlKey) {
        return;
      }

      const { cy } = renderer;
      const zoom = cy.zoom();
      const pan = cy.pan();
      const position = renderer.projectIntoViewport(event.clientX, event.clientY);
      const renderedPosition = [position[0] * zoom + pan.x, position[1] * zoom + pan.y];

      if (
        renderer.hoverData.draggingEles ||
        renderer.hoverData.dragging ||
        renderer.hoverData.cxtStarted ||
        inBoxSelection()
      ) {
        event.preventDefault();
        return;
      }

      if (cy.panningEnabled() && cy.userPanningEnabled() && cy.zoomingEnabled()) {
        event.preventDefault();
        renderer.data.wheelZooming = true;

        clearTimeout(renderer.data.wheelTimeout);

        renderer.data.wheelTimeout = setTimeout(() => {
          renderer.data.wheelZooming = false;
          renderer.redrawHint('eles', true);
          renderer.redraw();
        }, 150);

        let diff;

        if (event.deltaY != null) {
          diff = event.deltaY / -250;
        } else if (event.wheelDeltaY != null) {
          diff = event.wheelDeltaY / 1000;
        } else {
          diff = event.wheelDelta / 1000;
        }

        diff *= renderer.wheelSensitivity;

        if (event.deltaMode === 1) {
          diff *= 33;
        }

        let newZoom = cy.zoom() * 10 ** diff;

        if (event.type === 'gesturechange') {
          newZoom = renderer.gestureStartZoom * event.scale;
        }

        cy.zoom({
          level: newZoom,
          renderedPosition: {
            x: renderedPosition[0],
            y: renderedPosition[1],
          },
        });

        cy.emit(event.type === 'gesturechange' ? 'pinchzoom' : 'scrollzoom');
      }
    },

    resizeHandler() {
      if (!this.$cy) {
        return;
      }

      this.$cy.invalidateSize();
      this.$cy.fit(this.$cy.nodes(), HEALTHCHECK_NETWORK_GRAPH_OPTIONS.fitPadding);
    },

    mountCytoscape(options = this.options) {
      this.$cy = cytoscape({
        container: this.$refs.canvasWrapper,
        userZoomingEnabled: !this.ctrlWheelZoom,
        ...options,
      });

      this.$cy.nodeHtmlLabel(this.nodeHtmlLabelOptions, {
        enablePointerEvents: true,
      });
      this.$cy.bind('mouseover', 'node', this.nodeMouseoverHandler);
      this.$cy.bind('mouseout', 'node', this.nodeMouseoutHandler);
      this.$cy.bind('tap', 'node', this.nodeTapHandler);

      if (this.ctrlWheelZoom) {
        this.$el.addEventListener('wheel', this.wheelHandler);
      }
    },

    destroyCytoscape() {
      if (this.$cy) {
        this.$cy.destroy();
        this.$cy = null;
      }

      if (this.ctrlWheelZoom && this.$el) {
        this.$el.removeEventListener('wheel', this.wheelHandler);
      }
    },

    remountCytoscape(options) {
      this.destroyCytoscape();
      this.mountCytoscape(options);
    },

    showTooltip({ renderedPosition, data, zoom = 1 }) {
      const { getContent, offsetX = 0, offsetY = 0 } = this.tooltipOptions;

      if (!getContent) {
        return;
      }

      const tooltipInnerHTML = getContent(data);

      if (tooltipInnerHTML) {
        const { tooltip: tooltipElement } = this.$refs;

        tooltipElement.innerHTML = tooltipInnerHTML;

        const { height } = tooltipElement.getBoundingClientRect();
        const x = renderedPosition.x - offsetX * zoom;
        const y = renderedPosition.y - offsetY * zoom;

        tooltipElement.style.left = `${Math.round(x)}px`;
        tooltipElement.style.top = `${Math.round(y < height ? height : y)}px`;
        tooltipElement.style.opacity = 1;
      }
    },

    hideTooltip() {
      this.$refs.tooltip.style.opacity = 0;
    },

    nodeTapHandler(event) {
      this.$emit('node:tap', event);
    },

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
    cursor: grabbing;

    &--pointer {
      cursor: pointer;
    }
  }
}

.v-tooltip__content {
  position: absolute;
  top: 0;
  max-width: 420px;
  padding: 8px 12px;
  color: rgb(var(--v-theme-surface));
  background: rgba(var(--v-theme-on-surface), 0.9);
  border-radius: 4px;
  opacity: 0;
  transition: all 0.2s linear;
  transform: translate(-50%, 50%);
  pointer-events: none;
  z-index: 10;

  :deep(*) {
    color: inherit !important;
  }
}
</style>
