<template>
  <v-card-text>
    <remediation-jobs-list
      :remediation-jobs="remediationJobs"
      :pending="remediationJobsPending"
      :total-items="remediationJobsMeta.total_count"
      :options.sync="options"
      :removable="hasDeleteAnyRemediationJobAccess"
      :updatable="hasUpdateAnyRemediationJobAccess"
      :duplicable="hasCreateAnyRemediationJobAccess"
      @remove-selected="showRemoveSelectedRemediationJobsModal"
      @remove="showRemoveRemediationJobModal"
      @duplicate="showDuplicateRemediationJobModal"
      @edit="showEditRemediationJobModal"
    />
  </v-card-text>
</template>

<script>
import { omit } from 'lodash';

import { MODALS } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';
import { permissionsTechnicalRemediationJobMixin } from '@/mixins/permissions/technical/remediation-job';

import RemediationJobsList from './remediation-jobs-list.vue';

export default {
  components: { RemediationJobsList },
  mixins: [
    localQueryMixin,
    entitiesRemediationJobMixin,
    permissionsTechnicalRemediationJobMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      this.fetchRemediationJobsList({ params });
    },

    showEditRemediationJobModal(remediationJob) {
      this.$modals.show({
        name: MODALS.createRemediationJob,
        config: {
          remediationJob,
          title: this.$t('modals.createRemediationJob.edit.title'),
          action: async (job) => {
            await this.updateRemediationJob({ id: remediationJob._id, data: job });

            this.$popups.success({
              text: this.$t('modals.createRemediationJob.edit.popups.success', {
                jobName: job.name,
              }),
            });

            await this.fetchList();
          },
        },
      });
    },

    showDuplicateRemediationJobModal(remediationJob) {
      this.$modals.show({
        name: MODALS.createRemediationJob,
        config: {
          remediationJob: omit(remediationJob, ['_id']),
          title: this.$t('modals.createRemediationJob.duplicate.title'),
          action: async (job) => {
            await this.createRemediationJob({ data: job });

            this.$popups.success({
              text: this.$t('modals.createRemediationJob.duplicate.popups.success', {
                jobName: remediationJob.name,
              }),
            });

            await this.fetchList();
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
