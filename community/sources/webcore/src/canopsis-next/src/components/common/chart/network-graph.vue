<template>
  <div class="network-graph">
    <div
      ref="canvasWrapper"
      class="network-graph__canvas-wrapper"
    />
    <div
      ref="tooltip"
      class="v-tooltip__content menuable__content__active network-graph__tooltip"
    />
  </div>
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
    ctrlWheelZoom: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    nodeHtmlLabelOptions(newOptions, oldOptions) {
      if (!isEqual(newOptions, oldOptions)) {
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
  beforeDestroy() {
    this.destroyCytoscape();

    window.removeEventListener('resize', this.debouncedResizeHandler);
  },
  methods: {
    /**
     * Custom handler for ctrl + wheel zooming
     *
     * @param {WheelEvent} e
     */
    wheelHandler(e) {
      const r = this.$cy.renderer();

      const inBoxSelection = () => r.selection[4] !== 0;

      if (r.scrollingPage) {
        return;
      } // while scrolling, ignore wheel-to-zoom

      if (!e.ctrlKey) {
        return;
      }

      const { cy } = r;
      const zoom = cy.zoom();
      const pan = cy.pan();
      const pos = r.projectIntoViewport(e.clientX, e.clientY);
      const rpos = [pos[0] * zoom + pan.x, pos[1] * zoom + pan.y];

      if (r.hoverData.draggingEles || r.hoverData.dragging || r.hoverData.cxtStarted || inBoxSelection()) {
        // if pan dragging or cxt dragging, wheel movements make no zoom
        e.preventDefault();
        return;
      }

      if (cy.panningEnabled() && cy.userPanningEnabled() && cy.zoomingEnabled()) {
        e.preventDefault();
        r.data.wheelZooming = true;

        clearTimeout(r.data.wheelTimeout);

        r.data.wheelTimeout = setTimeout(() => {
          r.data.wheelZooming = false;
          r.redrawHint('eles', true);
          r.redraw();
        }, 150);

        let diff;

        if (e.deltaY != null) {
          diff = e.deltaY / -250;
        } else if (e.wheelDeltaY != null) {
          diff = e.wheelDeltaY / 1000;
        } else {
          diff = e.wheelDelta / 1000;
        }

        diff *= r.wheelSensitivity;
        const needsWheelFix = e.deltaMode === 1;

        if (needsWheelFix) {
          // fixes slow wheel events on ff/linux and ff/windows
          diff *= 33;
        }

        let newZoom = cy.zoom() * (10 ** diff);

        if (e.type === 'gesturechange') {
          newZoom = r.gestureStartZoom * e.scale;
        }

        cy.zoom({
          level: newZoom,
          renderedPosition: {
            x: rpos[0],
            y: rpos[1],
          },
        });

        cy.emit(e.type === 'gesturechange' ? 'pinchzoom' : 'scrollzoom');
      }
    },

    resizeHandler() {
      this.$cy.invalidateSize();
      this.$cy.fit(this.$cy.nodes(), HEALTHCHECK_NETWORK_GRAPH_OPTIONS.fitPadding);
    },

    mountCytoscape(options = this.options) {
      this.$cy = cytoscape({
        container: this.$refs.canvasWrapper,
        userZoomingEnabled: !this.ctrlWheelZoom,

        ...options,
      });

      this.$cy.nodeHtmlLabel(this.nodeHtmlLabelOptions);
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

      if (this.ctrlWheelZoom) {
        this.$el.removeEventListener('wheel', this.wheelHandler);
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

        tooltipEl.innerHTML = tooltipInnerHTML;

        const { height } = tooltipEl.getBoundingClientRect();

        const x = renderedPosition.x - (offsetX * zoom);
        const y = renderedPosition.y - (offsetY * zoom);

        tooltipEl.style.left = `${Math.round(x)}px`;
        tooltipEl.style.top = `${Math.round(y < height ? height : y)}px`;
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
