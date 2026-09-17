<template>
  <c-layout gap="2" column>
    <c-draggable-list-field
      :model-value="jobsField"
      :disabled="disabled"
      :group="draggableGroup"
      handle=".job-drag-handler"
      ghost-class="grey"
      class="flex-column gap-2"
      item-key="key"
      @update:model-value="jobsField = $event"
    >
      <template #item="{ element: job, index }">
        <v-card>
          <v-card-text>
            <remediation-instruction-job-field
              :job="job"
              :jobs="jobsItems"
              :name="job.key"
              :job-number="index + 1"
              :disabled="disabled"
              :loading="pending"
              @update:job="updateItemInArray(index, $event)"
              @remove="removeJob(index)"
            />
          </v-card-text>
        </v-card>
      </template>
    </c-draggable-list-field>

    <c-btn-with-error
      :error="hasJobsErrors ? t('remediation.instruction.errors.jobRequired') : ''"
      :disabled="disabled"
      outlined
      @click="addJob"
    >
      {{ t('remediation.instruction.addJob') }}
    </c-btn-with-error>
  </c-layout>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { remediationInstructionJobToForm } from '@/helpers/entities/remediation/instruction/form';

import { useFormArray, useFormErrorsMap } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';
import { useValidationAttachMinValueForField } from '@/hooks/validator/validation-attach-min-value';

import RemediationInstructionJobField from './fields/remediation-instruction-job-field.vue';

export default {
  components: {
    RemediationInstructionJobField,
  },
  props: {
    name: {
      type: String,
      default: 'jobs',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const jobsField = useFormArray('jobs');
    const { hasError } = useFormErrorsMap();
    const { fetchRemediationJobsListWithoutStore } = useRemediationJob();

    const draggableGroup = {
      name: 'remediation-instruction-jobs',
    };

    const jobsItems = ref([]);
    const hasJobsErrors = computed(() => hasError(props.name));

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const { data: jobs } = await fetchRemediationJobsListWithoutStore({
        params: {
          limit: MAX_LIMIT,
        },
      });

      jobsItems.value = jobs;
    });

    const { asyncValidateMinValueRule } = useValidationAttachMinValueForField(
      props.name,
      () => jobsField.value.length,
    );

    const addJob = () => {
      jobsField.value = [...jobsField.value, remediationInstructionJobToForm()];
    };
    const updateItemInArray = (index, value) => {
      const jobs = [...jobsField.value];

      jobs[index] = value;
      jobsField.value = jobs;
    };
    const removeJob = (index) => {
      jobsField.value = jobsField.value.filter((job, jobIndex) => jobIndex !== index);
    };

    watch(jobsField, asyncValidateMinValueRule);

    onMounted(fetchList);

    return {
      t,
      draggableGroup,
      jobsField,
      jobsItems,
      pending,
      hasJobsErrors,
      updateItemInArray,
      addJob,
      removeJob,
    };
  },
};
</script>
