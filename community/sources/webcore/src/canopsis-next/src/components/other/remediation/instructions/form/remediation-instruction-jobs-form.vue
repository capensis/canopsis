<template lang="pug">
  v-layout.mt-2(column)
    v-layout(row)
      v-flex(v-if="!jobs.length", xs12)
        v-alert(:value="true", type="info") {{ $t('remediation.instruction.emptyJobs') }}
    h3.subheading.font-weight-bold {{ $t('remediation.instruction.listJobs') }}
    c-draggable-list-field(
      v-field="jobs",
      :disabled="disabled",
      :group="draggableGroup",
      handle=".job-drag-handler",
      ghost-class="white"
    )
      v-card.my-2(v-for="(job, index) in jobs", :key="job.key")
        v-card-text
          remediation-instruction-job-field.py-1(
            v-field="jobs[index]",
            :jobs="jobsItems",
            :name="job.key",
            :job-number="index + 1",
            :disabled="disabled",
            @remove="removeItemFromArray(index)"
          )
    v-layout(row, align-center)
      v-btn.ml-0(
        :color="hasJobsErrors ? 'error' : 'primary'",
        :disabled="disabled",
        outline,
        @click="addJob"
      ) {{ $t('remediation.instruction.addJob') }}
      span.error--text(v-show="hasJobsErrors") {{ $t('remediation.instruction.errors.jobRequired') }}
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { remediationInstructionJobToForm } from '@/helpers/forms/remediation-instruction';

import { formArrayMixin } from '@/mixins/form';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';

import RemediationInstructionJobField from './fields/remediation-instruction-job-field.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionJobField,
  },
  mixins: [formArrayMixin, entitiesRemediationJobMixin],
  model: {
    prop: 'jobs',
    event: 'input',
  },
  props: {
    jobs: {
      type: Array,
      default: () => ([]),
    },
    name: {
      type: String,
      default: 'jobs',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isDragging: false,
      jobsItems: [],
    };
  },
  computed: {
    hasJobsErrors() {
      return this.errors.has(this.name);
    },

    draggableGroup() {
      return {
        name: 'remediation-instruction-jobs',
        pull: false,
        put: false,
      };
    },
  },
  watch: {
    jobs() {
      this.$validator.validate(this.name);
    },
  },
  mounted() {
    this.fetchList();
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: 'min_value:1',
      getter: () => this.jobs.length,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    addJob() {
      this.addItemIntoArray(remediationInstructionJobToForm());
    },

    async fetchList() {
      const { data: jobs } = await this.fetchRemediationJobsListWithoutStore({
        params: {
          limit: MAX_LIMIT,
        },
      });

      this.jobsItems = jobs;
    },
  },
};
</script>
