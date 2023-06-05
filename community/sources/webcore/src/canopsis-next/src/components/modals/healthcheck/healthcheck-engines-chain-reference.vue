<template lang="pug">
  modal-wrapper.fill-height(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      div.pre-wrap {{ $t('healthcheck.chainConfigurationInvalid') }}
      div.healthcheck-engine-chain-reference
        healthcheck-network-graph(
          :engines-graph="enginesGraph",
          :engines-parameters="enginesParameters"
        )
    template(#actions="")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.ok') }}
</template>

<script>
import {
  MODALS,
  HEALTHCHECK_ENGINES_NAMES,
  HEALTHCHECK_SERVICES_NAMES,
  HEALTHCHECK_ENGINES_REFERENCE_EDGES,
  HEALTHCHECK_ENGINES_PRO_REFERENCE_EDGES,
  PRO_ENGINES,
} from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.healthcheckEnginesChainReference,
  components: { HealthcheckNetworkGraph, ModalWrapper },
  mixins: [modalInnerMixin, healthcheckNodesMixin, entitiesInfoMixin],
  computed: {
    enginesGraph() {
      return {
        nodes: Object.values(HEALTHCHECK_ENGINES_NAMES)
          .filter(name => this.isProVersion || !PRO_ENGINES.includes(name)),
        edges: this.isProVersion
          ? HEALTHCHECK_ENGINES_PRO_REFERENCE_EDGES
          : HEALTHCHECK_ENGINES_REFERENCE_EDGES,
      };
    },

    enginesParameters() {
      return this.enginesGraph.nodes.reduce((acc, name) => {
        acc[name] = { name, is_running: true };

        return acc;
      }, {});
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
