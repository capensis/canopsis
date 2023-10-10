<template lang="pug">
  modal-wrapper.fill-height(close)
    template(#title="")
      span {{ title }}
    template(#text="")
      div.pre-wrap {{ $t('healthcheck.chainConfigurationInvalid') }}
      v-fade-transition
        v-layout(v-if="pending", justify-center)
          v-progress-circular(color="primary", indeterminate)
        div.healthcheck-engine-chain-reference(v-else)
          healthcheck-network-graph(
            :engines-graph="enginesGraph",
            :engines-parameters="enginesParameters"
          )
    template(#actions="")
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.ok') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MODALS, HEALTHCHECK_SERVICES_NAMES } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { healthcheckNodesMixin } from '@/mixins/healthcheck/healthcheck-nodes';
import { entitiesInfoMixin } from '@/mixins/entities/info';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';

import ModalWrapper from '../modal-wrapper.vue';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  name: MODALS.healthcheckEnginesChainReference,
  components: { HealthcheckNetworkGraph, ModalWrapper },
  mixins: [modalInnerMixin, healthcheckNodesMixin, entitiesInfoMixin],
  data() {
    return {
      pending: false,
      enginesGraph: {
        nodes: [],
        edges: [],
      },
    };
  },
  computed: {
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
  methods: {
    ...mapActions({
      fetchEnginesOrderWithoutStore: 'fetchEnginesOrderWithoutStore',
    }),

    async fetchList() {
      try {
        this.pending = false;
        this.enginesGraph = await this.fetchEnginesOrderWithoutStore();
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss">
.healthcheck-engine-chain-reference {
  height: 65vh;
}
</style>
