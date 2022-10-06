<template lang="pug">
  v-layout.mt-2(column)
    v-layout(row)
      v-flex(v-if="!jobs.length", xs12)
        v-alert(:value="true", type="info") {{ $t('remediationInstructions.emptyJobs') }}
    h3.subheading.font-weight-bold {{ $t('remediationInstructions.listJobs') }}
    draggable(v-field="jobs", :options="draggableOptions")
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
      ) {{ $t('remediationInstructions.addJob') }}
      span.error--text(v-show="hasJobsErrors") {{ $t('remediationInstructions.errors.jobRequired') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { MAX_LIMIT } from '@/constants';
import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { remediationInstructionJobToForm } from '@/helpers/forms/remediation-instruction';

import { formArrayMixin } from '@/mixins/form';
import { entitiesRemediationJobMixin } from '@/mixins/entities/remediation/job';

import RemediationInstructionJobField from './fields/remediation-instruction-job-field.vue';

export default {
  inject: ['$validator'],
  components: {
    Draggable,
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

    draggableOptions() {
      return {
        disabled: this.disabled,
        animation: VUETIFY_ANIMATION_DELAY,
        handle: '.job-drag-handler',
        ghostClass: 'white',
        group: {
          name: 'remediation-instruction-jobs',
          pull: false,
          put: false,
        },
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
