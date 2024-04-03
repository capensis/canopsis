<template>
  <v-layout column>
    <v-layout align-center>
      <v-flex
        class="pr-2"
        xs1
      >
        <c-draggable-step-number
          :disabled="disabled"
          :color="hasChildrenError ? 'error' : 'primary'"
          drag-class="job-drag-handler"
        >
          {{ jobNumber }}
        </c-draggable-step-number>
      </v-flex>
      <v-flex xs10>
        <v-autocomplete
          v-field="job.job"
          v-validate="'required'"
          :items="jobs"
          :label="$tc('remediation.instruction.job')"
          :error-messages="errors.collect(jobFieldName)"
          :name="jobFieldName"
          :disabled="disabled"
          item-text="name"
          item-value="_id"
          return-object
        />
      </v-flex>
      <v-flex xs1>
        <v-layout justify-center>
          <c-action-btn
            v-if="!disabled"
            type="delete"
            @click="$emit('remove')"
          />
        </v-layout>
      </v-flex>
    </v-layout>
    <v-flex
      offset-xs1
      xs11
    >
      <c-workflow-field
        v-field="job.stop_on_fail"
        :disabled="disabled"
        :label="$t('remediation.instruction.jobWorkflow')"
        :continue-label="$t('remediation.instruction.remainingJob')"
      />
    </v-flex>
  </v-layout>
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
