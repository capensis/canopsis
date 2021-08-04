<template lang="pug">
  div.position-relative
    c-page-header
    div#cy(ref="cover")
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import cytoscape from 'cytoscape';
import cytoscapeNodeHtmlLabel from 'cytoscape-node-html-label';
import cytoscapeDagre from 'cytoscape-dagre';

// import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

cytoscapeNodeHtmlLabel(cytoscape);
cytoscape.use(cytoscapeDagre);

export default {
  // components: { HealthcheckNetworkGraph },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
      pending: true,
      response: {
        services: [
          {
            name: 'MongoDB',
            status: 1,
          },
          {
            name: 'Redis',
            status: 0,
          },
          {
            name: 'RabbitMQ',
            status: 0,
          },
        ],
        engines: {
          nodes: [
            {
              name: 'engine-fifo',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-action',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-axe',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-che',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-correlation',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-dynamic-infos',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-pbehavior',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
            {
              name: 'engine-service',
              instances: 1,
              min_instances: 2,
              max_instances: 4,
              queue_length: 0,
              time: 1627885018,
              status: 1,
            },
          ],
          edges: [
            {
              from: 'engine-fifo',
              to: 'engine-che',
            },
            {
              from: 'engine-axe',
              to: 'engine-correlation',
            },
            {
              from: 'engine-che',
              to: 'engine-pbehavior',
            },
            {
              from: 'engine-correlation',
              to: 'engine-service',
            },
            {
              from: 'engine-correlation',
              to: 'engine-che',
            },
            {
              from: 'engine-pbehavior',
              to: 'engine-axe',
            },
            {
              from: 'engine-action',
              to: 'engine-che',
            },
          ],
        },
      },
    };
  },
  computed: {
    services() {
      return this.response.services;
    },

    engines() {
      return this.response.engines;
    },
  },
  mounted() {
    const SPACING_FACTOR = 2;
    const NODES_SPACE = 85;
    const SERVICES_NAMES = {
      mongo: 'MongoDB',
      redis: 'Redis',
      rabbit: 'RabbitMQ',
      events: 'Events',
      api: 'API',
      healthcheck: 'Healthcheck',
      enginesChain: 'Engines chain',
    };

    const ENGINES_NAMES = {
      fifo: 'engine-fifo',
    };

    const getDiff = (factor = 0) => NODES_SPACE * SPACING_FACTOR * factor;

    const SERVICES_RENDERED_POSITIONS_DIFF = {
      [SERVICES_NAMES.events]: { x: getDiff(-3), y: getDiff(1) },
      [SERVICES_NAMES.mongo]: { x: getDiff(-3), y: 0 },
      [SERVICES_NAMES.api]: { x: getDiff(-2), y: 0 },
      [SERVICES_NAMES.rabbit]: { x: getDiff(-2), y: getDiff(1) },
      [SERVICES_NAMES.healthcheck]: { x: getDiff(-2), y: getDiff(-0.5) },
      [SERVICES_NAMES.redis]: { x: getDiff(-1), y: 0 },
      [SERVICES_NAMES.enginesChain]: { x: 0, y: getDiff(-0.5) },
    };

    const servicesNodes = [
      { data: { id: SERVICES_NAMES.events }, classes: ['without-node'] },
      { data: { id: SERVICES_NAMES.mongo, color: '#2fab63' }, classes: ['service-title-bottom'] },
      { data: { id: SERVICES_NAMES.healthcheck }, classes: ['without-node'] },
      { data: { id: SERVICES_NAMES.api, color: '#2fab63' } },
      { data: { id: SERVICES_NAMES.rabbit, color: '#2fab63' }, classes: ['service-title-bottom'] },
      { data: { id: SERVICES_NAMES.redis, color: '#2fab63' }, classes: ['service-title-bottom'] },
      {
        data: { id: SERVICES_NAMES.enginesChain }, classes: ['without-node', 'without-node__error'], locked: true, selectable: false,
      },
    ];

    const servicesEdges = [
      {
        data: {
          source: SERVICES_NAMES.mongo,
          target: SERVICES_NAMES.api,
          label: 'Status check',
          color: 'gray',
        },
      },
      {
        data: {
          source: SERVICES_NAMES.events,
          target: SERVICES_NAMES.rabbit,
          color: 'black',
        },
      },
      {
        classes: ['service-edge__vertical'],
        data: {
          source: SERVICES_NAMES.rabbit,
          target: SERVICES_NAMES.api,
          label: 'Status check',
          color: 'gray',
        },
      },
      {
        classes: ['service-edge__horizontal-multiline'],
        data: {
          source: SERVICES_NAMES.redis,
          target: SERVICES_NAMES.api,
          label: 'FIFO data\nRedis check',
          color: 'gray',
        },
      },
      {
        data: {
          source: SERVICES_NAMES.api,
          target: SERVICES_NAMES.healthcheck,
          color: 'gray',
        },
      },
      {
        data: {
          source: SERVICES_NAMES.rabbit,
          target: ENGINES_NAMES.fifo,
          color: 'black',
        },
      },
      {
        data: {
          source: ENGINES_NAMES.fifo,
          target: SERVICES_NAMES.redis,
          label: 'RabbitMQ status',
          color: 'gray',
        },
      },
    ];

    const elements = {
      nodes: this.response.engines.nodes.map(node => ({ data: { id: node.name, ...node, color: '#2fab63' }, classes: ['engine'] })),
      edges: this.response.engines.edges.map(edge => ({ data: { source: edge.from, target: edge.to, color: 'black' } })),
    };

    this.$cy = cytoscape({
      container: this.$refs.cover,
      style: [
        {
          selector: 'core',
          style: {
            'active-bg-size': 0,
          },
        },
        {
          selector: 'node',
          style: {
            width: 60,
            height: 60,
            'border-width': 3,
          },
        },
        {
          selector: `node[id="${SERVICES_NAMES.enginesChain}"]`,
          style: {
            events: 'no',
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
          selector: '.without-node',
          style: {
            'background-opacity': 0,
            'border-width': 0,
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
          selector: '.service-edge__vertical',
          style: {
            'text-rotation': '90deg',
            'text-margin-x': 10,
            'text-margin-y': 0,
          },
        },
        {
          selector: '.service-edge__horizontal-multiline',
          style: {
            'text-margin-y': 0,
            'line-height': 1.5,
          },
        },
      ],

      elements,
      wheelSensitivity: 0.5,
      layout: {
        name: 'dagre',
        direction: true,
        spacingFactor: SPACING_FACTOR,
      },

      ready() {
        const [fifoNode] = this.nodes(`[id="${ENGINES_NAMES.fifo}"]`);
        const fifoRenderedPosition = fifoNode.renderedPosition();

        this.add([
          ...servicesNodes.map(node => ({
            ...node,

            group: 'nodes',
            renderedPosition: {
              x: SERVICES_RENDERED_POSITIONS_DIFF[node.data.id].x + fifoRenderedPosition.x,
              y: SERVICES_RENDERED_POSITIONS_DIFF[node.data.id].y + fifoRenderedPosition.y,
            },
          })),

          ...servicesEdges.map(edge => ({ group: 'edges', ...edge })),
        ]);

        this.fit(this.nodes(), 15);
      },
    });

    this.$cy.nodeHtmlLabel([
      {
        query: 'node',
        valign: 'top',
        halign: 'right',
        valignBox: 'top',
        halignBox: 'right',
        tpl: data => `<span>${data.id}</span>`,
      },
      {
        query: '.engine',
        valign: 'center',
        halign: 'right',
        valignBox: 'center',
        halignBox: 'right',
        tpl: data =>
          `<div class="ml-1"><div class="subheading">${data.id}</div><div class="body-1 grey--text darken-3">Queue length 312/500</div><div class="body-1 grey--text darken-3">Instances 5/3</div><div></div></div>`,
      },
      {
        query: '.without-node',
        valign: 'center',
        halign: 'center',
        tpl: data => `<span class="headline">${data.id}</span>`,
      },
      {
        query: '.without-node__error',
        valign: 'center',
        halign: 'center',
        tpl: data => `<span class="v-chip theme--light red white--text cursor-pointer"><span class="v-chip__content cursor-pointer headline">${data.id}</span></span>`,
      },
      {
        query: '.service-title-bottom',
        valign: 'bottom',
        halign: 'center',
        valignBox: 'bottom',
        halignBox: 'center',
        tpl: data => `<span class="subheading">${data.id}</span>`,
      },
    ]);

    this.$cy.unbind('mouseover');
    this.$cy.bind('mouseover', 'node', (event) => {
      const { target } = event;

      target.style('border-color', 'gray');
      this.$refs.cover.style.cursor = 'pointer';
    });

    this.$cy.unbind('mouseout');
    this.$cy.bind('mouseout', 'node', (event) => {
      const { target } = event;

      target.style('border-color', target.data('color'));
      this.$refs.cover.style.cursor = 'grabbing';
    });

    this.$cy.unbind('tap');
    this.$cy.bind('tap', 'node', (event) => {
      event.preventDefault();
      event.stopPropagation();
    });

    // this.fetchList();
  },
  beforeDestroy() {
    this.$cy.destroy();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      this.engines = await this.fetchEnginesListWithoutStore();

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
#cy {
  background: white;
  position: relative;
  width: 100%;
  height: 100vh;
  cursor: grabbing;
}
</style>
