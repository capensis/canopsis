<template>
  <div>
    <p v-if="engine.is_unknown" class="pre-line">
      {{ $t('healthcheck.unknown') }}
    </p>
    <p v-if="hasNodeDownIssue" class="pre-line">
      {{ $t('healthcheck.issues.node_down', { name }) }}
    </p>
    <p v-if="hasNoQuorumIssue" class="pre-line">
      {{ $t('healthcheck.issues.no_quorum', { name }) }}
    </p>
    <p v-if="hasQuorumDegradedIssue" class="pre-line">
      {{ $t('healthcheck.issues.quorum_degraded', { name }) }}
    </p>
    <div v-if="hasTooFewInstancesIssue">
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
    <p v-if="hasQueueOverflowedIssue" class="pre-wrap">
      {{ $t('healthcheck.queueOverflowed', { queueLength, maxQueueLength }) }}
    </p>
    <p v-if="hasInstancesConfigMismatchIssue" class="pre-wrap">
      {{ $t('healthcheck.invalidInstancesConfiguration') }}
    </p>
  </div>
</template>

<script>
import {
  HEALTHCHECK_ENGINE_ISSUES,
  HEALTHCHECK_SERVICE_ISSUES,
  PRO_ENGINES,
} from '@/constants';

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

    issues() {
      return this.engine.issues ?? [];
    },

    hasNodeDownIssue() {
      return this.issues.includes(HEALTHCHECK_SERVICE_ISSUES.nodeDown);
    },

    hasNoQuorumIssue() {
      return this.issues.includes(HEALTHCHECK_SERVICE_ISSUES.noQuorum);
    },

    hasQuorumDegradedIssue() {
      return this.issues.includes(HEALTHCHECK_SERVICE_ISSUES.quorumDegraded);
    },

    hasTooFewInstancesIssue() {
      return this.issues.includes(HEALTHCHECK_ENGINE_ISSUES.tooFewInstances);
    },

    hasQueueOverflowedIssue() {
      return this.issues.includes(HEALTHCHECK_ENGINE_ISSUES.queueOverflowed);
    },

    hasInstancesConfigMismatchIssue() {
      return this.issues.includes(HEALTHCHECK_ENGINE_ISSUES.instancesConfigMismatch);
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
