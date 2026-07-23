<template>
  <div>
    <p v-if="!engine.is_running" class="pre-line">
      {{ systemDownMessage }}
    </p>
    <div v-if="engine.is_too_few_instances">
      <div class="pre-wrap">
        {{ $t('healthcheck.activeInstances', { instances, minInstances, optimalInstances }) }}
      </div>
      <healthcheck-engine-instance-diagram
        :instances="instances"
        :min-instances="minInstances"
        :optimal-instances="optimalInstances"
        :is-pro-engine="isProEngine"
      />
    </div>
    <p v-if="engine.is_queue_overflown" class="pre-wrap">
      {{ $t('healthcheck.queueOverflowed', { queueLength, maxQueueLength }) }}
    </p>
    <p v-if="engine.is_diff_instances_config" class="pre-wrap">
      {{ $t('healthcheck.invalidInstancesConfiguration') }}
    </p>
  </div>
</template>

<script>
import { HEALTHCHECK_ENGINES_NAMES, HEALTHCHECK_SERVICES_NAMES, PRO_ENGINES } from '@/constants';

import HealthcheckEngineInstanceDiagram from './partials/healthcheck-engine-instance-diagram.vue';

export default {
  components: { HealthcheckEngineInstanceDiagram },

  props: {
    engine: {
      type: Object,
      default: () => ({}),
    },
    maxQueueLength: {
      type: Number,
      default: 0,
    },
  },

  computed: {
    name() {
      return this.getNodeName(this.engine.name);
    },

    isProEngine() {
      return PRO_ENGINES.includes(this.engine.name);
    },

    queueLength() {
      return this.engine.queue_length;
    },

    instances() {
      return this.engine.instances;
    },

    minInstances() {
      return this.engine.min_instances;
    },

    optimalInstances() {
      return this.engine.optimal_instances;
    },

    systemDownMessage() {
      const messageKey =
        {
          [HEALTHCHECK_ENGINES_NAMES.fifo]: 'healthcheck.engineDown',
          [HEALTHCHECK_SERVICES_NAMES.timescaleDB]: 'healthcheck.timescaleDown',
        }[this.engine.name] || 'healthcheck.engineDownOrSlow';

      return this.$t(messageKey, { name: this.name });
    },
  },

  methods: {
    getNodeName(name) {
      const nodeLabelKey = `healthcheck.nodes.${name}.name`;

      return this.$te(nodeLabelKey) ? this.$t(nodeLabelKey) : name;
    },
  },
};
</script>
