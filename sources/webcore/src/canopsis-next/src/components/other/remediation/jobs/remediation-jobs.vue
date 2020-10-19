<template lang="pug">
  v-card-text
    remediation-jobs-list(
      :remediation-jobs="remediationJobs",
      :pending="remediationJobsPending",
      :total-items="remediationJobsMeta.total_count",
      :pagination.sync="pagination",
      @remove-selected="showRemoveSelectedRemediationJobsModal",
      @remove="showRemoveRemediationJobModal",
      @edit="showEditRemediationJobModal"
    )
</template>

<script>
import { MODALS } from '@/constants';

import entitiesRemediationJobsMixin from '@/mixins/entities/remediation/jobs';
import localQueryMixin from '@/mixins/query-local/query';

import RemediationJobsList from './remediation-jobs-list.vue';

export default {
  components: { RemediationJobsList },
  inject: ['$validator'],
  mixins: [
    entitiesRemediationJobsMixin,
    localQueryMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      this.fetchRemediationJobsList({ params: this.getQuery() });
    },

    showEditRemediationJobModal(remediationJob) {
      this.$modals.show({
        name: MODALS.createRemediationJob,
        config: {
          remediationJob,
          action: async (job) => {
            await this.updateRemediationJob({ id: remediationJob._id, data: job });

            this.$popups.success({
              text: this.$t('modals.createRemediationInstruction.edit.popups.success', {
                jobName: remediationJob.name,
              }),
            });

            await this.fetchJobsList();
          },
        },
      });
    },

    showRemoveRemediationJobModal(remediationJob) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeRemediationJob({ id: remediationJob._id });
            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedRemediationJobsModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id: id }) => this.removeRemediationJob({ id })));
            await this.fetchList();
          },
        },
      });
    },
  },
};
</script>
