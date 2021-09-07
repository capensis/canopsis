<template lang="pug">
  v-layout(column)
    v-layout(row, align-center)
      v-flex.pr-2(xs1)
        c-draggable-step-number(
          drag-class="job-drag-handler",
          :disabled="disabled",
          :color="hasChildrenError ? 'error' : 'primary'"
        ) {{ jobNumber }}
      v-flex(xs10)
        v-autocomplete(
          v-field="job.job",
          v-validate="'required'",
          :items="jobs",
          :label="$tc('remediationInstructions.job')",
          :error-messages="errors.collect(jobFieldName)",
          :name="jobFieldName",
          :disabled="disabled",
          item-text="name",
          item-value="_id",
          return-object
        )
      v-flex(xs1)
        v-layout(justify-center)
          c-action-btn(v-if="!disabled", type="delete", @click="$emit('remove')")
    v-flex(offset-xs1, xs11)
      c-workflow-field(
        v-field="job.stop_on_fail",
        :disabled="disabled",
        :label="$t('remediationInstructions.jobWorkflow')",
        :continue-label="$t('remediationInstructions.remainingJob')"
      )
</template>

<script>
import { formMixin, validationChildrenMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
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
    jobs: {
      type: Array,
      default: () => [],
    },
    jobNumber: {
      type: [Number, String],
      default: 0,
    },
    name: {
      type: String,
      default: 'job',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    jobFieldName() {
      return `${this.name}.job`;
    },
  },
};
</script>
