<template>
  <network-graph
    ref="networkGraph"
    :options="options"
    :node-html-label-options="nodeHtmlLabelsOptions"
    :tooltip-options="tooltipOptions"
    class="healthcheck-network-graph"
    @node:tap="nodeTapHandler"
    @node:mouseover="nodeMouseoverHandler"
    @node:mouseout="nodeMouseoutHandler"
  />
</template>

<script>
import { get, minBy, isNil } from 'lodash';

import { COLORS } from '@/config';
import {
  MODALS,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_ENGINES_NAMES,
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS,
  HEALTHCHECK_SERVICES_RENDERED_POSITIONS_DIFF_FACTORS,
} from '@/constants';

import { getHealthcheckNodeRenderedPositionDiff } from '@/helpers/charts/healthcheck';
import { getHealthcheckNodeColor } from '@/helpers/entities/healthcheck/color';
import { convertDateToString } from '@/helpers/date/date';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

export default {
  components: { NetworkGraph },
  mixins: [healthcheckNodesMixin],
  props: {
    services: {
      type: Array,
      required: false,
    },
    enginesGraph: {
      type: Object,
      default: () => ({}),
    },
    enginesParameters: {
      type: Object,
      default: () => ({}),
    },
    hasInvalidEnginesOrder: {
      type: Boolean,
      default: false,
    },
    showDescription: {
      type: Boolean,
      default: false,
    },
    maxQueueLength: {
      type: Number,
      default: 0,
    },
    getTooltip: {
      type: Function,
      required: false,
    },
  },
  computed: {
    servicesNodes() {
      return this.services
        ? [
          ...this.services.map(node => ({
            classes: node.name === HEALTHCHECK_SERVICES_NAMES.rabbit ? [] : ['service-title--bottom'],
            data: {
              id: node.name,
              name: node.name,
              color: getHealthcheckNodeColor(node),
              is_running: node.is_running,
            },
          })),

          {
            classes: ['service-node--without-node'],
            data: {
              id: HEALTHCHECK_SERVICES_NAMES.events,
              name: HEALTHCHECK_SERVICES_NAMES.events,
            },
          },
          {
            classes: ['service-node--without-node'],
            data: {
              id: HEALTHCHECK_SERVICES_NAMES.healthcheck,
              name: HEALTHCHECK_SERVICES_NAMES.healthcheck,
            },
          },
          {
            data: {
              id: HEALTHCHECK_SERVICES_NAMES.api,
              name: HEALTHCHECK_SERVICES_NAMES.api,
              color: COLORS.primary,
            },
          },
          {
            locked: true,
            classes: ['service-node--without-node'],
            data: {
              id: HEALTHCHECK_SERVICES_NAMES.enginesChain,
              name: HEALTHCHECK_SERVICES_NAMES.enginesChain,
            },
          },
        ]
        : [];
    },

    servicesEdges() {
      return this.services
        ? [
          {
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.mongo,
              target: HEALTHCHECK_SERVICES_NAMES.api,
              color: COLORS.healthcheck.edgeGray,
              label: this.getNodeEdgeLabel(HEALTHCHECK_SERVICES_NAMES.mongo),
            },
          },
          {
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.events,
              target: HEALTHCHECK_SERVICES_NAMES.rabbit,
              color: COLORS.healthcheck.edgeBlack,
            },
          },
          {
            classes: ['service-edge--diagonal-right'],
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.rabbit,
              target: HEALTHCHECK_SERVICES_NAMES.api,
              color: COLORS.healthcheck.edgeGray,
              label: this.getNodeEdgeLabel(HEALTHCHECK_SERVICES_NAMES.rabbit),
            },
          },
          {
            classes: ['service-edge--diagonal-left'],
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.timescaleDB,
              target: HEALTHCHECK_SERVICES_NAMES.api,
              color: COLORS.healthcheck.edgeGray,
              label: this.getNodeEdgeLabel(HEALTHCHECK_SERVICES_NAMES.timescaleDB),
            },
          },
          {
            classes: ['service-edge--multiline'],
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.redis,
              target: HEALTHCHECK_SERVICES_NAMES.api,
              color: COLORS.healthcheck.edgeGray,
              label: this.getNodeEdgeLabel(HEALTHCHECK_SERVICES_NAMES.redis),
            },
          },
          {
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.api,
              target: HEALTHCHECK_SERVICES_NAMES.healthcheck,
              color: COLORS.healthcheck.edgeGray,
            },
          },
          {
            data: {
              source: HEALTHCHECK_SERVICES_NAMES.rabbit,
              target: HEALTHCHECK_ENGINES_NAMES.fifo,
              color: COLORS.healthcheck.edgeBlack,
            },
          },
          {
            classes: ['service-edge--multiline'],
            data: {
              source: HEALTHCHECK_ENGINES_NAMES.fifo,
              target: HEALTHCHECK_SERVICES_NAMES.redis,
              color: COLORS.healthcheck.edgeGray,
              label: this.getNodeEdgeLabel(HEALTHCHECK_ENGINES_NAMES.fifo),
            },
          },
        ]
        : [];
    },

    servicesElements() {
      return [
        ...this.servicesNodes.map(node => ({ group: 'nodes', ...node })),
        ...this.servicesEdges.map(edge => ({ group: 'edges', ...edge })),
      ];
    },

    enginesGraphNodes() {
      return get(this.enginesGraph, 'nodes') || [];
    },

    enginesGraphEdges() {
      return get(this.enginesGraph, 'edges') || [];
    },

    hasSNMPNode() {
      return this.enginesGraphNodes.includes(HEALTHCHECK_ENGINES_NAMES.snmp);
    },

    hasSNMPEdge() {
      return this.enginesGraphEdges.some(
        ({ from, to }) => from === HEALTHCHECK_ENGINES_NAMES.snmp && to === HEALTHCHECK_ENGINES_NAMES.fifo,
      );
    },

    enginesNodes() {
      const hasFifoNode = this.enginesGraphNodes.includes(HEALTHCHECK_ENGINES_NAMES.fifo);
      const nodes = [...this.enginesGraphNodes];
      const parameters = { ...this.enginesParameters };

      if (!hasFifoNode) {
        nodes.push(HEALTHCHECK_ENGINES_NAMES.fifo);

        if (!parameters[HEALTHCHECK_ENGINES_NAMES.fifo]) {
          parameters[HEALTHCHECK_ENGINES_NAMES.fifo] = { is_running: true };
        }
      }

      return nodes.map((name) => {
        const nodeParameters = get(parameters, name, {});
        const data = {
          ...nodeParameters,

          name,
          id: name,
        };

        data.color = getHealthcheckNodeColor(data);

        return {
          classes: ['engine'],
          data,
        };
      });
    },

    enginesEdges() {
      return this.enginesGraphEdges.map(
        edge => ({
          data: {
            source: edge.from,
            target: edge.to,
            color: COLORS.healthcheck.edgeBlack,
          },
        }),
      );
    },

    enginesElements() {
      return [
        ...this.enginesNodes.map(node => ({ group: 'nodes', ...node })),
        ...this.enginesEdges.map(edge => ({ group: 'edges', ...edge })),
      ];
    },

    nodeHtmlLabelsOptions() {
      return [
        {
          query: 'node',
          valign: 'top',
          halign: 'right',
          valignBox: 'top',
          halignBox: 'right',
          tpl: data => `<div class="subtitle-1">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: '.service-title--bottom',
          valign: 'bottom',
          halign: 'center',
          valignBox: 'bottom',
          halignBox: 'center',
          tpl: data => `<div class="subtitle-1">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: `node[id="${HEALTHCHECK_SERVICES_NAMES.rabbit}"]`,
          valign: 'bottom',
          halign: 'left',
          valignBox: 'bottom',
          halignBox: 'left',
          tpl: data => `<div class="subtitle-1">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: '.service-node--without-node',
          valign: 'center',
          halign: 'center',
          tpl: data => `<div class="headline">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: `node[id="${HEALTHCHECK_SERVICES_NAMES.events}"]`,
          halign: 'center',
          tpl: data => `<div class="headline">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: '.engine',
          valign: 'center',
          halign: 'right',
          valignBox: 'center',
          halignBox: 'right',
          tpl: this.getEngineHtmlLabel,
        },
        {
          query: `node[id="${HEALTHCHECK_SERVICES_NAMES.enginesChain}"]`,
          valign: 'center',
          halign: 'center',
          tpl: this.getEnginesChainHtmlLabel,
        },
      ];
    },

    styleOption() {
      return [
        {
          selector: 'core',
          style: {
            'active-bg-size': 0,
          },
        },
        {
          selector: 'node',
          style: {
            width: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSize,
            height: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSize,
            'border-width': 3,
          },
        },
        {
          selector: `node[id="${HEALTHCHECK_SERVICES_NAMES.enginesChain}"]`,
          style: {
            width: 200,
            height: 30,
            'overlay-color': 'white',
          },
        },
        {
          selector: `node[id="${HEALTHCHECK_SERVICES_NAMES.events}"]`,
          style: {
            width: 100,
            height: 30,
          },
        },
        {
          selector: `node[id="${HEALTHCHECK_SERVICES_NAMES.healthcheck}"]`,
          style: {
            width: 140,
            height: 30,
          },
        },
        {
          selector: 'node[color]',
          style: {
            'background-color': 'data(color)',
            'border-color': 'data(color)',
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'target-arrow-shape': 'triangle',
            'target-arrow-color': 'data(color)',
            'line-color': 'data(color)',
          },
        },
        {
          selector: 'edge[label]',
          style: {
            label: 'data(label)',
            color: 'gray',
            'font-size': '14px',
            'text-wrap': 'wrap',
            'text-margin-y': -10,
          },
        },
        {
          selector: '.service-node--without-node',
          style: {
            'background-opacity': 0,
            'border-width': 0,
          },
        },
        {
          selector: '.service-edge--vertical',
          style: {
            'text-rotation': '90deg',
            'text-margin-x': 10,
            'text-margin-y': 0,
          },
        },
        {
          selector: '.service-edge--diagonal-right',
          style: {
            'text-rotation': '45deg',
            'text-margin-x': 10,
            'text-margin-y': 0,
          },
        },
        {
          selector: '.service-edge--diagonal-left',
          style: {
            'text-rotation': '-45deg',
            'text-margin-x': -10,
            'text-margin-y': 0,
          },
        },
        {
          selector: '.service-edge--multiline',
          style: {
            'text-margin-y': 0,
            'line-height': 1.5,
          },
        },
      ];
    },

    options() {
      const { servicesElements, hasSNMPNode, hasSNMPEdge } = this;

      return {
        elements: this.enginesElements,
        style: this.styleOption,
        wheelSensitivity: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.wheelSensitivity,
        minZoom: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.minZoom,
        maxZoom: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.maxZoom,
        layout: {
          name: 'dagre',
          direction: true,
          spacingFactor: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.spacingFactor,
        },
        /**
         * Here we are adding the service elements on the left side of fifo-engine
         */
        ready() {
          const zoom = this.zoom();
          const [fifoNode] = this.getElementById(HEALTHCHECK_ENGINES_NAMES.fifo);
          let positionWithMinY = fifoNode
            ? fifoNode.renderedPosition()
            : minBy(this.nodes().renderedPosition(), 'y');

          if (!positionWithMinY) {
            positionWithMinY = {
              x: 0,
              y: 0,
            };
          }

          const preparedElements = servicesElements.map((element) => {
            if (element.group !== 'nodes') {
              return element;
            }

            let itemDiffFactors = HEALTHCHECK_SERVICES_RENDERED_POSITIONS_DIFF_FACTORS[element.data.id];

            if (!itemDiffFactors && element.data.id === HEALTHCHECK_SERVICES_NAMES.enginesChain) {
              itemDiffFactors = { x: 0, y: hasSNMPNode && hasSNMPEdge ? -1.5 : -0.5 };
            }

            if (!itemDiffFactors) {
              return element;
            }

            return {
              ...element,

              renderedPosition: {
                x: (getHealthcheckNodeRenderedPositionDiff(itemDiffFactors.x) * zoom) + positionWithMinY.x,
                y: (getHealthcheckNodeRenderedPositionDiff(itemDiffFactors.y) * zoom) + positionWithMinY.y,
              },
            };
          });

          this.add(preparedElements);
          this.fit(this.nodes(), HEALTHCHECK_NETWORK_GRAPH_OPTIONS.fitPadding);
        },
      };
    },

    tooltipOptions() {
      return {
        offsetY: (HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSize / 2) + 5,
        getContent: this.getTooltipContent,
      };
    },
  },
  watch: {
    services() {
      this.$nextTick(this.refreshGraph);
    },

    enginesGraph() {
      this.$nextTick(this.refreshGraph);
    },

    enginesParameters() {
      this.$nextTick(this.refreshParameters);
    },
  },
  methods: {
    /**
     * Refresh whole cytoscape graph
     */
    refreshGraph() {
      this.$refs.networkGraph.remountCytoscape();
    },

    /**
     * Refresh data on cytoscape nodes cytoscape
     */
    refreshParameters() {
      Object.entries(this.enginesParameters).forEach(([nodeId, nodeParameters]) => {
        const node = this.$refs.networkGraph.$cy.getElementById(nodeId);

        if (node.length) {
          const data = { ...node.data(), ...nodeParameters };
          const color = getHealthcheckNodeColor(data);

          data.color = color;

          node.data(data);
          node.style({
            'background-color': color,
            'border-color': color,
          });
        }
      });
    },

    /**
     * Get content for the tooltip by node data
     *
     * @param {Object} data
     * @return {string}
     */
    getTooltipContent(data) {
      if (this.getTooltip) {
        return this.getTooltip(data);
      }

      const TOOLTIP_TEXTS_MAP = {
        [HEALTHCHECK_SERVICES_NAMES.enginesChain]: this.hasInvalidEnginesOrder
          ? this.$t('healthcheck.invalidEnginesOrder')
          : '',
        [HEALTHCHECK_SERVICES_NAMES.timescaleDB]: !data.is_running
          ? this.$t('healthcheck.metricsUnavailable')
          : '',
      };

      const tooltipMessage = TOOLTIP_TEXTS_MAP[data.id] ?? this.getTooltipText(data);

      return tooltipMessage ? `<div class="pre-wrap">${tooltipMessage}</div>` : '';
    },

    /**
     * Get html label for engines chain service node
     *
     * @param {string} id
     * @return {string}
     */
    getEnginesChainHtmlLabel({ id }) {
      const nameSpan = `<span class="headline">${this.getNodeName(id)}</span>`;
      let wrapperDivClass = '';
      let icon = '';

      if (this.hasInvalidEnginesOrder) {
        wrapperDivClass = 'v-chip theme--light px-3 py-1 error white--text cursor-pointer';
        icon = '<i class="v-icon material-icons theme--dark ml-2">warning</i>';
      }

      return `<div class="${wrapperDivClass}">${nameSpan}${icon}</div>`;
    },

    /**
     * Get html label for engine node
     *
     * @param {string} id
     * @param {number} time
     * @param {number} instances
     * @param {number} minInstances
     * @param {number} queueLength
     * @param {number} maxQueueLength
     * @param {boolean} isTooFewInstances
     * @param {boolean} isQueueOverflown
     * @param {boolean} isUnknown
     * @return {string}
     */
    getEngineHtmlLabel({
      id,
      time,
      instances,
      min_instances: minInstances,
      queue_length: queueLength,
      is_too_few_instances: isTooFewInstances,
      is_queue_overflown: isQueueOverflown,
      is_unknown: isUnknown,
    }) {
      const elements = [`<div class="subtitle-1">${this.getNodeName(id)}</div>`];

      const getInfoDiv = (message, hasError) => `<div class="body-1 grey--text darken-3 ${hasError ? 'error--text' : ''}">${message}</div>`;

      if (!isUnknown && this.showDescription) {
        elements.push(
          isNil(time)
            ? null
            : getInfoDiv(convertDateToString(time)),
          isNil(queueLength)
            ? null
            : getInfoDiv(
              this.$t('healthcheck.queueLength', { queueLength, maxQueueLength: this.maxQueueLength }),
              isQueueOverflown,
            ),
          isNil(instances) && isNil(minInstances)
            ? null
            : getInfoDiv(
              this.$t('healthcheck.instancesCount', { instances, minInstances }),
              isTooFewInstances,
            ),
        );
      }

      return `<div class="ml-2">${elements.filter(Boolean).join('')}</div>`;
    },

    /**
     * Show modal window for each 'engine' node
     *
     * @param {Object} engine
     */
    showEngineModal(engine) {
      const excludedServices = [
        HEALTHCHECK_SERVICES_NAMES.api,
        HEALTHCHECK_SERVICES_NAMES.healthcheck,
        HEALTHCHECK_SERVICES_NAMES.events,
      ];

      if (excludedServices.includes(engine.name)) {
        return;
      }

      if (this.isWrongEngine(engine)) {
        this.$modals.show({
          name: MODALS.healthcheckEngine,
          config: {
            engine,
            maxQueueLength: this.maxQueueLength,
          },
        });
      }
    },

    /**
     * Show modal window for 'enginesChain' node
     */
    showEngineChainReferenceModal() {
      this.$modals.show({
        name: MODALS.healthcheckEnginesChainReference,
      });
    },

    /**
     * Handler for the 'tap' event on node for cytoscape
     *
     * @param {Object} target
     */
    nodeTapHandler({ target }) {
      const engine = target.data();

      if (engine.name !== HEALTHCHECK_SERVICES_NAMES.enginesChain) {
        this.showEngineModal(engine);
      } else if (this.hasInvalidEnginesOrder) {
        this.showEngineChainReferenceModal(engine);
      }
    },

    /**
     * Handler for the 'mouseover' event on node for cytoscape
     *
     * @param {Object} target
     */
    nodeMouseoverHandler({ target }) {
      target.style('border-color', COLORS.healthcheck.edgeGray);
    },

    /**
     * Handler for the 'mouseout' event on node for cytoscape
     *
     * @param {Object} target
     */
    nodeMouseoutHandler({ target }) {
      target.style('border-color', target.data('color'));
    },
  },
};
</script>

<style lang="scss" scoped>
.healthcheck-network-graph {
  overflow: hidden;
}
</style>
