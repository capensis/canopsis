<template lang="pug">
  div
    c-the-page-header {{ $t('common.instructionRating') }}
    remediation-instruction-stats-list(
      :remediation-instruction-stats="remediationInstructionStats",
      :pending="remediationInstructionStatsPending",
      :pagination.sync="pagination",
      :total-items="remediationInstructionStatsMeta.total_count"
    )
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionStatsMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;
      // TODO: Should be removed
      params.from = 1;
      params.to = Date.now();

      this.fetchRemediationInstructionStatsList({ params });
    },
  },
};
</script>
