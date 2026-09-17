<template>
  <c-layout column>
    <c-layout gap="2" align="center">
      <c-draggable-step-number
        :disabled="disabled"
        :color="hasChildrenError ? 'error' : 'primary'"
        drag-class="job-drag-handler"
      >
        {{ jobNumber }}
      </c-draggable-step-number>
      <div class="flex-grow-1">
        <c-autocomplete
          :model-value="job.job"
          :loading="loading"
          :items="jobs"
          :label="t('remediation.instruction.job')"
          :name="jobFieldName"
          :disabled="disabled"
          item-title="name"
          item-value="_id"
          return-object
          required
          @update:model-value="updateField('job', $event)"
        />
      </div>
      <c-layout justify="center">
        <c-action-btn v-if="!disabled" type="delete" @click="$emit('remove')" />
      </c-layout>
    </c-layout>
    <c-workflow-field
      :model-value="job.stop_on_fail"
      :disabled="disabled"
      :label="t('remediation.instruction.jobWorkflow')"
      :continue-label="t('remediation.instruction.remainingJob')"
      @update:model-value="updateField('stop_on_fail', $event)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useValidationChildren } from '@/hooks/validator/validation-children';

export default {
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
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const { hasChildrenError } = useValidationChildren();

    const jobFieldName = computed(() => `${props.name}.job`);

    return {
      t,
      hasChildrenError,
      jobFieldName,
      updateField,
    };
  },
};
</script>
