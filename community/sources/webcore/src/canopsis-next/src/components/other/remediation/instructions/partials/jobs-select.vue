<template>
  <choose-expansion-panel
    :entities="jobs"
    :label="$tc('remediation.instruction.job', 2)"
    content-key="name"
    clearable
    @remove="removeJob"
    @clear="clear"
  >
    <choose-jobs-lists
      :jobs="jobs"
      @select="updateJobs"
    />
  </choose-expansion-panel>
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import ChooseExpansionPanel from '@/components/common/choose-expansion-panel/choose-expansion-panel.vue';
import ChooseJobsLists from '@/components/other/remediation/jobs/partials/choose-jobs-lists.vue';

export default {
  components: { ChooseJobsLists, ChooseExpansionPanel },
  mixins: [formBaseMixin],
  model: {
    prop: 'jobs',
    event: 'input',
  },
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    clear() {
      this.updateModel([]);
    },

    updateJobs(jobs) {
      this.updateModel([...this.jobs, ...jobs]);
    },

    removeJob(removingJob) {
      const updatedEntities = this.jobs.filter(job => job._id !== removingJob._id);

      this.updateModel(updatedEntities);
    },
  },
};
</script>
