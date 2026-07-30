<template>
  <choose-expansion-panel
    :entities="jobs"
    :label="$t('remediation.instruction.job', 2)"
    content-key="name"
    clearable
    @remove="removeJob"
    @clear="clear"
  >
    <choose-jobs-lists :jobs="jobs" @select="updateJobs" />
  </choose-expansion-panel>
</template>

<script>
import ChooseExpansionPanel from '@/components/common/choose-expansion-panel/choose-expansion-panel.vue';
import ChooseJobsLists from '@/components/other/remediation/jobs/partials/choose-jobs-lists.vue';

export default {
  components: { ChooseJobsLists, ChooseExpansionPanel },
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['update:jobs'],
  methods: {
    clear() {
      this.$emit('update:jobs', []);
    },

    updateJobs(jobs) {
      this.$emit('update:jobs', [...this.jobs, ...jobs]);
    },

    removeJob(removingJob) {
      const updatedEntities = this.jobs.filter(job => job._id !== removingJob._id);

      this.$emit('update:jobs', updatedEntities);
    },
  },
};
</script>
