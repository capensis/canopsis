<template>
  <c-layout column>
    <c-layout align="center">
      <v-flex class="pr-2" xs1>
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
          v-model="jobField"
          v-validate="'required'"
          :loading="loading"
          :items="jobs"
          :label="$t('remediation.instruction.job')"
          :error-messages="errors.collect(jobFieldName)"
          :name="jobFieldName"
          :disabled="disabled"
          item-title="name"
          item-value="_id"
          return-object
        />
      </v-flex>
      <v-flex xs1>
        <c-layout justify="center">
          <c-action-btn v-if="!disabled" type="delete" @click="$emit('remove')" />
        </c-layout>
      </v-flex>
    </c-layout>
    <v-flex offset-xs1 xs11>
      <c-workflow-field
        v-model="stopOnFailField"
        :disabled="disabled"
        :label="$t('remediation.instruction.jobWorkflow')"
        :continue-label="$t('remediation.instruction.remainingJob')"
      />
    </v-flex>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
  inject: ['$validator'],
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
    loading: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:job', 'remove'],
  setup(props, { emit }) {
    const { hasChildrenError } = useValidationChildren();

    const jobFieldName = computed(() => `${props.name}.job`);
    // TODO: fix that
    const createFieldModel = field =>
      computed({
        get: () => props.job[field],
        set: value =>
          emit('update:job', {
            ...props.job,
            [field]: value,
          }),
      });
    const jobField = createFieldModel('job');
    const stopOnFailField = createFieldModel('stop_on_fail');

    return {
      hasChildrenError,
      jobField,
      stopOnFailField,
      jobFieldName,
    };
  },
};
</script>
