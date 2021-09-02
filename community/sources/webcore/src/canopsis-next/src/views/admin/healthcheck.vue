<template lang="pug">
  div.healthcheck
    c-progress-overlay(:pending="pending")
    c-page-header
    div.mb-3.text-xs-center(v-if="hasAnyError")
      v-chip(color="error", dark)
        span.subheading {{ $t('healthcheck.systemStatusChipError') }}
    v-tabs(v-model="activeTab", centered, slider-color="primary")
      v-tab {{ $t('common.systemStatus') }}
      v-tab {{ $tc('common.graph', 2) }}
      v-tab {{ $tc('common.parameter', 2) }}
    v-tabs-items.white.healthcheck__tabs(v-model="activeTab")
      v-tab-item.healthcheck__graph-tab
        healthcheck-network-graph(
          v-if="!pending && !hasServerError",
          :services="response.services",
          :engines="response.engines",
          :has-invalid-engines-order="response.has_invalid_engines_order",
          :max-queue-length="response.max_queue_length",
          show-description
        )
        h2.my-4.headline.text-xs-center(v-else-if="hasServerError") {{ $t('healthcheck.systemStatusServerError') }}
      v-tab-item(lazy)
        healthcheck-graphs(:max-queue-length="response.max_queue_length")
      v-tab-item(lazy)
        healthcheck-parameters
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import { isEqual } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import HealthcheckNetworkGraph from '@/components/other/healthcheck/exploitation/healthcheck-network-graph.vue';
import HealthcheckGraphs from '@/components/other/healthcheck/exploitation/healthcheck-graphs.vue';
import HealthcheckParameters from '@/components/other/healthcheck/healthcheck-parameters.vue';

import entitiesEngineRunInfoMixin from '@/mixins/entities/engine-run-info';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: { HealthcheckParameters, HealthcheckNetworkGraph, HealthcheckGraphs },
  mixins: [entitiesEngineRunInfoMixin],
  data() {
    return {
      activeTab: 0,
      pending: true,
      response: {},
      hasServerError: false,
    };
  },
  computed: {
    hasAnyError() {
      return this.hasServerError || this.response.has_invalid_engines_order;
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchHealthcheckEnginesWithoutStore: 'fetchEnginesWithoutStore',
    }),

    setData(data) {
      if (!isEqual(data, this.response)) {
        this.response = data;
      }
    },

    async fetchList() {
      try {
        this.hasServerError = false;
        this.pending = true;

        this.response = await this.fetchHealthcheckEnginesWithoutStore();
      } catch (err) {
        this.hasServerError = true;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.healthcheck {
  display: flex;
  flex-direction: column;
  height: 100%;

  &__tabs {
    flex: 1;

    /deep/ .v-window__container {
      height: 100%;
    }
  }

  &__graph-tab {
    height: 100%;
  }
}
</style>
