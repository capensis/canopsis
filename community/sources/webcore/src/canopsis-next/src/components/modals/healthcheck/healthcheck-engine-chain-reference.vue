<template lang="pug">
  modal-wrapper.fill-height(close)
    template(slot="title")
      span {{ title }}
    template(slot="text")
      div.pre-wrap {{ $t('healthcheck.chainConfigurationInvalid') }}
      div.healthcheck-engine-chain-reference
        healthcheck-network-graph(:engines="engines")
    template(slot="actions")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
</template>

<script>
import {
  MODALS,
  HEALTHCHECK_ENGINES_NAMES,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_ENGINES_REFERENCE_EDGES,
  HEALTHCHECK_ENGINES_CAT_REFERENCE_EDGES,
  CAT_ENGINES,
} from '@/constants';

import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.healthcheckEngineChainReference,
  components: { HealthcheckNetworkGraph, ModalWrapper },
  mixins: [healthcheckNodesMixin, entitiesInfoMixin],
  computed: {
    engines() {
      const nodeNames = Object.values(HEALTHCHECK_ENGINES_NAMES)
        .filter(name => this.isCatVersion || !CAT_ENGINES.includes(name));

      return {
        nodes: nodeNames.map(name => ({
          name,
          is_running: true,
        })),
        edges: this.isCatVersion
          ? HEALTHCHECK_ENGINES_CAT_REFERENCE_EDGES
          : HEALTHCHECK_ENGINES_REFERENCE_EDGES,
      };
    },

    title() {
      return this.getNodeName(HEALTHCHECK_SERVICES_NAMES.enginesChain);
    },
  },
};
</script>

<style lang="scss">
.healthcheck-engine-chain-reference {
  height: 65vh;
}
</style>
