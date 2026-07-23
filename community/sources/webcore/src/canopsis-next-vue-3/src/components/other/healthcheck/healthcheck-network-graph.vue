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
import { get, isNil, minBy } from 'lodash';

import { COLORS } from '@/config';
import {
  HEALTHCHECK_ENGINES_NAMES,
  HEALTHCHECK_NETWORK_GRAPH_OPTIONS,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_SERVICES_RENDERED_POSITIONS_DIFF_FACTORS,
  MODALS,
} from '@/constants';

import {
  getHealthcheckInfoDiv,
  getHealthcheckNodeRenderedPositionDiff,
} from '@/helpers/charts/healthcheck';
import { convertDateToString } from '@/helpers/date/date';
import { getHealthcheckNodeColor } from '@/helpers/entities/healthcheck/color';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

export default {
  components: { NetworkGraph },

  props: {
    services: {
      type: Array,
      default: () => [],
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
      return this.services.length
        ? [
            ...this.services.map((node) => ({
              classes:
                node.name === HEALTHCHECK_SERVICES_NAMES.rabbit ? [] : ['service-title--bottom'],
              data: {
                id: node.name,
                name: node.name,
                color: getHealthcheckNodeColor(node),
                is_running: node.is_running,
                nodes: node.nodes,
                running_nodes: node.running_nodes,
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
      return this.services.length
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
        ...this.servicesNodes.map((node) => ({ group: 'nodes', ...node })),
        ...this.servicesEdges.map((edge) => ({ group: 'edges', ...edge })),
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
        ({ from, to }) =>
          from === HEALTHCHECK_ENGINES_NAMES.snmp && to === HEALTHCHECK_ENGINES_NAMES.fifo,
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
      return this.enginesGraphEdges.map((edge) => ({
        data: {
          source: edge.from,
          target: edge.to,
          color: COLORS.healthcheck.edgeBlack,
        },
      }));
    },

    enginesElements() {
      return [
        ...this.enginesNodes.map((node) => ({ group: 'nodes', ...node })),
        ...this.enginesEdges.map((edge) => ({ group: 'edges', ...edge })),
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
          tpl: (data) => `<div class="text-subtitle-1">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: '.service-title--bottom',
          valign: 'bottom',
          halign: 'center',
          valignBox: 'bottom',
          halignBox: 'center',
          tpl: (data) => `<div class="text-subtitle-1">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: `node[id="${HEALTHCHECK_SERVICES_NAMES.rabbit}"]`,
          valign: 'bottom',
          halign: 'left',
          valignBox: 'bottom',
          halignBox: 'left',
          tpl: this.getServiceNodeHtmlLabel,
        },
        {
          query: '.service-node--without-node',
          valign: 'center',
          halign: 'center',
          tpl: (data) => `<div class="text-h5">${this.getNodeName(data.id)}</div>`,
        },
        {
          query: `node[id="${HEALTHCHECK_SERVICES_NAMES.events}"]`,
          halign: 'center',
          tpl: (data) => `<div class="text-h5">${this.getNodeName(data.id)}</div>`,
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

            let itemDiffFactors =
              HEALTHCHECK_SERVICES_RENDERED_POSITIONS_DIFF_FACTORS[element.data.id];

            if (!itemDiffFactors && element.data.id === HEALTHCHECK_SERVICES_NAMES.enginesChain) {
              itemDiffFactors = { x: 0, y: hasSNMPNode && hasSNMPEdge ? -1.5 : -0.5 };
            }

            if (!itemDiffFactors) {
              return element;
            }

            return {
              ...element,
              renderedPosition: {
                x:
                  getHealthcheckNodeRenderedPositionDiff(itemDiffFactors.x) * zoom +
                  positionWithMinY.x,
                y:
                  getHealthcheckNodeRenderedPositionDiff(itemDiffFactors.y) * zoom +
                  positionWithMinY.y,
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
        offsetY: HEALTHCHECK_NETWORK_GRAPH_OPTIONS.nodeSize / 2 + 5,
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
    refreshGraph() {
      this.$refs.networkGraph?.remountCytoscape();
    },

    refreshParameters() {
      Object.entries(this.enginesParameters).forEach(([nodeId, nodeParameters]) => {
        const node = this.$refs.networkGraph?.$cy?.getElementById(nodeId);

        if (node?.length) {
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

    getNodeEdgeLabel(name) {
      const nodeEdgeLabelKey = `healthcheck.nodes.${name}.edgeLabel`;

      return this.$te(nodeEdgeLabelKey) ? this.$t(nodeEdgeLabelKey) : name;
    },

    getNodeName(name) {
      const nodeLabelKey = `healthcheck.nodes.${name}.name`;

      return this.$te(nodeLabelKey) ? this.$t(nodeLabelKey) : name;
    },

    getTooltipText(data = {}) {
      const statusKeys = [];
      const itemWithDefaultFlags = {
        is_running: true,
        is_queue_overflown: false,
        is_too_few_instances: false,
        is_diff_instances_config: false,
        ...data,
      };

      if (!itemWithDefaultFlags.is_running) {
        statusKeys.push('healthcheck.notRunning');
      }

      if (itemWithDefaultFlags.is_queue_overflown) {
        statusKeys.push('healthcheck.queueOverflow');
      }

      if (itemWithDefaultFlags.is_too_few_instances) {
        statusKeys.push('healthcheck.lackOfInstances');
      }

      if (itemWithDefaultFlags.is_diff_instances_config) {
        statusKeys.push('healthcheck.diffInstancesConfig');
      }

      if (itemWithDefaultFlags.is_unknown) {
        statusKeys.push('healthcheck.unknown');
      }

      return statusKeys
        .map((stateKey) => this.$t(stateKey, { name: this.getNodeName(data.name) }))
        .join('\n');
    },

    getTooltipContent(data) {
      if (this.getTooltip) {
        return this.getTooltip(data);
      }

      const tooltipTextsMap = {
        [HEALTHCHECK_SERVICES_NAMES.enginesChain]: this.hasInvalidEnginesOrder
          ? this.$t('healthcheck.invalidEnginesOrder')
          : '',
        [HEALTHCHECK_SERVICES_NAMES.timescaleDB]: !data.is_running
          ? this.$t('healthcheck.metricsUnavailable')
          : '',
      };

      const tooltipMessage = tooltipTextsMap[data.id] ?? this.getTooltipText(data);

      return tooltipMessage ? `<div class="pre-wrap">${tooltipMessage}</div>` : '';
    },

    getEnginesChainHtmlLabel({ id }) {
      const nameSpan = `<span class="text-h5">${this.getNodeName(id)}</span>`;
      let wrapperDivClass = '';
      let icon = '';

      if (this.hasInvalidEnginesOrder) {
        wrapperDivClass =
          'healthcheck-network-graph__chain-error text-white px-3 py-1 cursor-pointer';
        icon = '<i class="v-icon material-icons text-white ml-2">warning</i>';
      }

      return `<div class="${wrapperDivClass}">${nameSpan}${icon}</div>`;
    },

    getServiceNodeHtmlLabel(data) {
      const { id, nodes, running_nodes: runningNodes = 0 } = data;
      const elements = [`<div class="text-subtitle-1">${this.getNodeName(id)}</div>`];

      if (this.showDescription && !isNil(nodes)) {
        elements.push(
          getHealthcheckInfoDiv(this.$t('healthcheck.runningNodes', { nodes, runningNodes })),
        );
      }

      return `<div class="ml-2">${elements.join('')}</div>`;
    },

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
      const elements = [`<div class="text-subtitle-1">${this.getNodeName(id)}</div>`];

      if (!isUnknown && this.showDescription) {
        elements.push(
          isNil(time) ? null : getHealthcheckInfoDiv(convertDateToString(time)),
          isNil(queueLength)
            ? null
            : getHealthcheckInfoDiv(
                this.$t('healthcheck.queueLength', {
                  queueLength,
                  maxQueueLength: this.maxQueueLength,
                }),
                isQueueOverflown,
              ),
          isNil(instances) && isNil(minInstances)
            ? null
            : getHealthcheckInfoDiv(
                this.$t('healthcheck.instancesCount', { instances, minInstances }),
                isTooFewInstances,
              ),
        );
      }

      return `<div class="ml-2">${elements.filter(Boolean).join('')}</div>`;
    },

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

    showEngineChainReferenceModal() {
      this.$modals.show({
        name: MODALS.healthcheckEnginesChainReference,
      });
    },

    isWrongEngine(engine = {}) {
      return (
        !engine.is_running ||
        engine.is_queue_overflown ||
        engine.is_too_few_instances ||
        engine.is_diff_instances_config ||
        engine.is_unknown
      );
    },

    nodeTapHandler({ target }) {
      const engine = target.data();

      if (engine.name !== HEALTHCHECK_SERVICES_NAMES.enginesChain) {
        this.showEngineModal(engine);
      } else if (this.hasInvalidEnginesOrder) {
        this.showEngineChainReferenceModal();
      }
    },

    nodeMouseoverHandler({ target }) {
      target.style('border-color', COLORS.healthcheck.edgeGray);
    },

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

:deep(.healthcheck-network-graph__chain-error) {
  display: inline-flex;
  align-items: center;
  color: rgb(var(--v-theme-on-error));
  background: rgb(var(--v-theme-error));
  border-radius: 16px;
}
</style>
