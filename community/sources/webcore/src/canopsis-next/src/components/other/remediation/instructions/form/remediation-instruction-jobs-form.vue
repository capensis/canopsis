<template>
  <v-layout
    class="gap-2"
    column
  >
    <c-draggable-list-field
      v-field="jobs"
      :disabled="disabled"
      :group="draggableGroup"
      handle=".job-drag-handler"
      ghost-class="grey"
      class="flex-column gap-2"
    >
      <v-card
        v-for="(job, index) in jobs"
        :key="job.key"
      >
        <v-card-text>
          <remediation-instruction-job-field
            v-field="jobs[index]"
            :jobs="jobsItems"
            :name="job.key"
            :job-number="index + 1"
            :disabled="disabled"
            :loading="pending"
            @remove="removeJob(index)"
          />
        </v-card-text>
      </v-card>
    </c-draggable-list-field>

    <c-btn-with-error
      :error="hasJobsErrors ? $t('remediation.instruction.errors.jobRequired') : ''"
      :disabled="disabled"
      outlined
      @click="addJob"
    >
      {{ $t('remediation.instruction.addJob') }}
    </c-btn-with-error>
  </v-layout>
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { remediationInstructionJobToForm } from '@/helpers/entities/remediation/instruction/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationJob } from '@/hooks/store/modules/remediation-job';
import { useValidator } from '@/hooks/validator/validator';
import { useValidationAttachMinValueForField } from '@/hooks/validator/validation-attach-min-value';

import RemediationInstructionJobField from './fields/remediation-instruction-job-field.vue';

export default {
  components: {
    RemediationInstructionJobField,
  },
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
  setup(props, { emit }) {
    const draggableGroup = {
      name: 'remediation-instruction-jobs',
    };

    const validator = useValidator();
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);
    const { fetchRemediationJobsListWithoutStore } = useRemediationJob();

    const jobsItems = ref([]);

    const hasJobsErrors = computed(() => validator?.errors?.has(props.name) ?? false);

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
      () => props.jobs.length,
    );

    const addJob = () => addItemIntoArray(remediationInstructionJobToForm());

    watch(() => props.jobs, asyncValidateMinValueRule);

    onMounted(fetchList);

    return {
      draggableGroup,

      jobsItems,
      pending,

      hasJobsErrors,

      addJob,
      removeJob: removeItemFromArray,
    };
  },
};
</script>
