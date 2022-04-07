<template lang="pug">
  healthcheck-network-graph(
    :engines-graph="enginesGraph",
    :engines-parameters="enginesParameters",
    :get-tooltip="getEngineTooltip"
  )
</template>

<script>
import { sortBy } from 'lodash';

import { ENGINES_NAMES_TO_QUEUE_NAMES } from '@/constants';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

export default {
  components: { HealthcheckNetworkGraph },
  props: {
    engines: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    sortedEngines() {
      return sortBy(this.engines, ['name']);
    },

    enginesNames() {
      return this.sortedEngines.map(({ name }) => name);
    },

    enginesEdges() {
      return this.sortedEngines.reduce((acc, engine) => {
        const from = ENGINES_NAMES_TO_QUEUE_NAMES[engine.consume_queue];
        const to = ENGINES_NAMES_TO_QUEUE_NAMES[engine.publish_queue];

        if (from && to) {
          acc.push({ from, to });
        }

        return acc;
      }, []);
    },

    enginesGraph() {
      return {
        nodes: this.enginesNames,
        edges: this.enginesEdges,
      };
    },

    enginesParameters() {
      return this.enginesNames.reduce((acc, name) => {
        acc[name] = { name, is_running: true };

        return acc;
      }, {});
    },
  },
  methods: {
    getEngineTooltip(data) {
      const message = [
        this.$t(`healthcheck.nodes.${data.name}.name`),
        this.$t(`healthcheck.nodes.${data.name}.description`),
      ].join('\n');

      return `<div class="pre-wrap">${message}</div>`;
    },
  },
};
</script>
