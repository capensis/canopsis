<template lang="pug">
  v-layout(column)
    v-layout(row, align-center)
      v-flex(xs1)
        c-draggable-step-number(
          drag-class="job-drag-handler",
          :color="hasChildrenError ? 'error' : 'primary'"
        ) {{ jobNumber }}
      v-flex(xs10)
        remediation-instruction-jobs-select(v-field="job.job")
      v-flex(xs1)
        v-layout(justify-center)
          c-action-btn(type="delete", @click="remove")
    v-flex(offset-xs1, xs11)
      c-workflow-field(
        v-field="job.stop_on_fail",
        :label="$t('remediationInstructions.jobWorkflow')",
        :continue-label="$t('remediationInstructions.remainingJob')"
      )
</template>

<script>
import formMixin from '@/mixins/form';
import validationChildrenMixin from '@/mixins/form/validation-children';

import RemediationInstructionJobsSelect from './remediation-instruction-jobs-select.vue';

export default {
  inject: ['$validator'],
  components: { RemediationInstructionJobsSelect },
  mixins: [
    formMixin,
    validationChildrenMixin,
  ],
  model: {
    prop: 'job',
    event: 'input',
  },
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
    jobNumber: {
      type: [Number, String],
      default: 0,
    },
  },
  data() {
    return {
      expanded: true,
    };
  },
  computed: {
    fieldName() {
      return this.job.key ? this.job.key : '';
    },

    nameFieldName() {
      return `${this.fieldName}.name`;
    },
  },
  methods: {
    remove() {
      this.$emit('remove');
    },
  },
};
</script>

<style lang="scss" scoped>
  .job-expand {
    margin: 24px 2px 0 2px !important;
    width: 20px;
    height: 20px;
  }
</style>
