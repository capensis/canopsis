<template lang="pug">
  div.jobs-assigned
    v-layout(row)
      span.subheading {{ $t('remediationInstructionExecute.jobs.title') }}
    v-layout(column)
      v-data-table(:items="jobs", hide-actions)
        template(slot="headers", slot-scope="props")
          tr
            th
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.startedAt') }}
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.launchedAt') }}
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.completedAt') }}
        template(slot="items", slot-scope="props")
          remediation-instruction-assigned-job(
            :job="props.item",
            :key="`job_${props.item.job_id}`",
            @execute-job="executeJob"
          )
          remediation-instruction-assigned-job-alert-row(
            :job="props.item",
            :key="`alert_${props.item.job_id}`",
            @skip="cancelExecuteJob"
          )
</template>

<script>
import entitiesRemediationJobsExecutionsMixin from '@/mixins/entities/remediation/jobs-executions';
import entitiesRemediationInstructionExecutionMixin from '@/mixins/entities/remediation/executions';

import RemediationInstructionAssignedJob from './partials/remediation-instruction-assigned-job.vue';
import RemediationInstructionAssignedJobAlertRow from './partials/remediation-instruction-assigned-job-alert-row.vue';

export default {
  components: {
    RemediationInstructionAssignedJobAlertRow,
    RemediationInstructionAssignedJob,
  },
  mixins: [
    entitiesRemediationJobsExecutionsMixin,
    entitiesRemediationInstructionExecutionMixin,
  ],
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
    executionId: {
      type: String,
      required: true,
    },
    operationId: {
      type: [Number, String],
      required: true,
    },
  },
  methods: {
    async executeJob(job) {
      try {
        await this.createRemediationJobExecution({
          data: {
            execution: this.executionId,
            job: job.job_id,
            operation: this.operationId,
          },
        });
        await this.pingRemediationInstructionExecution({ id: this.executionId });
      } catch (err) {
        this.$popups.error({ text: err.error || this.$t('errors.default') });
      }
    },

    async cancelExecuteJob(job) {
      await this.cancelRemediationJobExecution({ id: job._id });
      await this.pingRemediationInstructionExecution({ id: this.executionId });
    },
  },
};
</script>

<style lang="scss" scoped>
  .jobs-assigned {
    tr {
      border-bottom: none !important;
    }

    tbody tr:hover {
      background: unset !important;
    }
  }
</style>
