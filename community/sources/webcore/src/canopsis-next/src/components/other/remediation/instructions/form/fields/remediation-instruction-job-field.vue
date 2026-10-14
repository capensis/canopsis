<template>
  <c-card-iterator-item
    :item-number="jobNumber"
    :disabled="disabled"
    :removable="!disabled"
    static-delete
    @remove="$emit('remove')"
  >
    <template #header>
      <v-row>
        <v-col cols="12">
          <c-autocomplete
            :model-value="job.job"
            :loading="loading"
            :items="jobs"
            :label="t('remediation.instruction.job')"
            :name="`${name}.job`"
            :disabled="disabled"
            item-title="name"
            item-value="_id"
            return-object
            required
            @update:model-value="updateField('job', $event)"
          />
        </v-col>
        <v-col cols="12">
          <c-workflow-field
          :model-value="job.stop_on_fail"
          :name="`${name}.stop_on_fail`"
            :disabled="disabled"
            :label="t('remediation.instruction.jobWorkflow')"
            :continue-label="t('remediation.instruction.remainingJob')"
            @update:model-value="updateField('stop_on_fail', $event)"
          />
        </v-col>
      </v-row>
    </template>
  </c-card-iterator-item>
</template>

<script>
import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

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

    return {
      t,
      updateField,
    };
  },
};
</script>
