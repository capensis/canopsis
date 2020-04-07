<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.healthcheck') }}
    v-layout(row, justify-center)
      v-chip(:color="overallInfo.color", text-color="white")
        v-progress-circular(v-if="pending", width="2", size="20", indeterminate)
        span(v-else) {{ overallInfo.text }}
    div.white.mt-2
      v-data-table(
        :headers="headers",
        :items="items",
        :loading="pending",
        hide-actions
      )
        tr(slot="items", slot-scope="{ item }")
          td
            span {{ item.type }}
          td
            v-chip(:color="item.state ? 'primary' : 'error'", small)
          td
            span {{ item.message }}
          td
            span {{ item.timestamp | date('long', false) }}
    .fab
      v-layout(column)
        refresh-btn(@click="fetchList")
</template>

<script>
import { isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { HEALTHCHECK_FIELDS } from '@/constants';

import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';

const { mapActions } = createNamespacedHelpers('healthcheck');

export default {
  components: { RefreshBtn },
  data() {
    return {
      healthcheck: {},
      pending: false,
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.type'), value: 'type' },
        { text: this.$t('tables.healthcheck.columns.state'), value: 'state' },
        { text: this.$t('tables.healthcheck.columns.message'), value: 'message' },
        { text: this.$t('tables.healthcheck.columns.timestamp'), value: 'timestamp' },
      ];
    },
    items() {
      if (isEmpty(this.healthcheck)) {
        return [];
      }

      return HEALTHCHECK_FIELDS.map(field => ({
        type: field,
        state: !this.healthcheck[field],
        timestamp: this.healthcheck.timestamp,
        message: this.healthcheck[field] || this.$t('tables.healthcheck.defaultMessage'),
      }));
    },
    overallInfo() {
      if (this.healthcheck.overall || isEmpty(this.healthcheck)) {
        return {
          text: this.$t('tables.healthcheck.overall.success'),
          color: 'primary',
        };
      }

      return {
        text: this.$t('tables.healthcheck.overall.error'),
        color: 'error',
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchHealthcheckWithoutStore: 'fetchWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      this.healthcheck = await this.fetchHealthcheckWithoutStore();

      this.pending = false;
    },
  },
};
</script>
