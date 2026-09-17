<template>
  <c-card-iterator-form
    :value="jobsField"
    :disabled="disabled"
    :draggable-group="draggableGroup"
    :required-error-message="t('remediation.instruction.errors.jobRequired')"
    :add-button-label="t('remediation.instruction.addJob')"
    :name="name"
    item-key="key"
    iterator-class="mb-2"
    required
    @add="addJob"
    @update:value="jobsField = $event"
  >
    <template #item="{ index, item: job }">
      <remediation-instruction-job-field
        :job="jobsField[index]"
        :jobs="jobsItems"
        :name="job.key"
        :job-number="index + 1"
        :disabled="disabled"
        :loading="pending"
        @update:job="updateItemInArray(index, $event)"
        @remove="removeJob(index)"
      />
    </template>
  </c-card-iterator-form>
</template>

<script>
import { ref, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { remediationInstructionJobToForm } from '@/helpers/entities/remediation/instruction/form';

import { useFormArray } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';

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
  setup() {
    const { t } = useI18n();
    const jobsField = useFormArray('jobs');
    const { fetchRemediationJobsListWithoutStore } = useRemediationJob();

    const draggableGroup = {
      name: 'remediation-instruction-jobs',
    };

    const jobsItems = ref([]);

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const { data: jobs } = await fetchRemediationJobsListWithoutStore({
        params: {
          limit: MAX_LIMIT,
        },
      });

      jobsItems.value = jobs;
    });

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

    onMounted(fetchList);

    return {
      t,
      draggableGroup,
      jobsField,
      jobsItems,
      pending,
      updateItemInArray,
      addJob,
      removeJob,
    };
  },
};
</script>
