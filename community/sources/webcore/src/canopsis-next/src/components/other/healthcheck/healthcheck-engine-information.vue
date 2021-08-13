<template lang="pug">
  div
    p.pre-wrap(v-if="!engine.is_running") {{ $t('healthcheck.engineDown', { name }) }}
    div(v-if="engine.is_too_few_instances")
      div.pre-wrap {{ $t('healthcheck.activeInstances', { instances, minInstances, optimalInstances }) }}
      healthcheck-engine-instance-diagram(
        :instances="instances",
        :min-instances="minInstances",
        :optimal-instances="optimalInstances",
        :is-cat-engine="isCatEngine"
      )
    p.pre-wrap(v-if="engine.is_queue_overflown")
      | {{ $t('healthcheck.queueOverflowed', { queueLength, maxQueueLength }) }}
    p.pre-wrap(v-if="engine.is_diff_instances_config") {{ $t('healthcheck.invalidInstancesConfiguration') }}
</template>

<script>
import { CAT_ENGINES } from '@/constants';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';

import HealthcheckEngineInstanceDiagram from './healthcheck-engine-instance-diagram.vue';

export default {
  components: { HealthcheckEngineInstanceDiagram },
  mixins: [healthcheckNodesMixin],
  props: {
    engine: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    name() {
      return this.getNodeName(this.engine.name);
    },

    isCatEngine() {
      return CAT_ENGINES.includes(this.engine.name);
    },

    queueLength() {
      return this.engine.queue_length;
    },

    maxQueueLength() {
      return this.engine.max_queue_length;
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
  },
};
</script>
