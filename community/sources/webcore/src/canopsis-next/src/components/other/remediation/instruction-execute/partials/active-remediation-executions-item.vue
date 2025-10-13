<template>
  <v-list-item @click="showExecutionModal">
    <v-list-item-content>
      <v-list-item-title class="primary--text cursor-pointer mb-3">
        {{ execution.name }}
      </v-list-item-title>
      <v-list-item-subtitle>
        <v-layout class="gap-2" column>
          <span>{{ execution.alarm?.display_name }}</span>
          <v-layout v-if="hasJobs" align-center>
            <strong class="mr-2">{{ $tc('remediation.instruction.job', 2) }}: </strong>
            <active-remediation-executions-item-jobs :jobs="execution.jobs" />
          </v-layout>
          <span v-if="execution.current_operation">
            <strong class="mr-2">{{ $t('common.step') }} {{ stepFullNumber }}:</strong>
            <span>{{ execution.current_operation.name }}</span>
          </span>
          <span v-if="timeToComplete" class="mt-1 grey--text">
            <v-icon class="mr-2" color="grey" small>timer</v-icon>
            <span>{{ $t('remediation.instructionExecute.timeToComplete', { duration: timeToComplete }) }}</span>
          </span>
        </v-layout>
      </v-list-item-subtitle>
    </v-list-item-content>
    <v-list-item-action v-if="hasActions">
      <c-action-btn
        :loading="removing"
        class="remove-card-btn"
        type="delete"
        icon="close"
        small
        icon-small
        @click="removeExecution"
      />
    </v-list-item-action>
  </v-list-item>
</template>

<script>
import { computed, toRef } from 'vue';

import { REMEDIATION_INSTRUCTION_TYPES, LONG_DURATION_FORMAT } from '@/constants';

import { convertDurationToString } from '@/helpers/date/duration';
import { getOperationNumber } from '@/helpers/entities/remediation/instruction/form';

import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

import {
  useShowExecutionModal,
} from '@/components/other/remediation/instruction-execute/hooks/active-remediation-executions-item';

import ActiveRemediationExecutionsItemJobs from './active-remediation-executions-item-jobs.vue';

export default {
  components: {
    ActiveRemediationExecutionsItemJobs,
  },
  props: {
    execution: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { readRemediationInstructionExecution } = useRemediationInstructionExecution();

    const currentOperation = computed(() => props.execution?.current_operation);
    const timeToComplete = computed(() => (
      currentOperation.value?.time_to_complete
        ? convertDurationToString(currentOperation.value.time_to_complete, LONG_DURATION_FORMAT)
        : null
    ));

    const stepFullNumber = computed(() => (
      currentOperation.value
        ? getOperationNumber(currentOperation.value.step_index + 1, currentOperation.value.index)
        : null
    ));

    const hasJobs = computed(() => Object.values(props.execution.jobs).some(value => !!value));

    const hasActions = computed(() => (
      props.execution.type === REMEDIATION_INSTRUCTION_TYPES.simpleManual && !props.execution.jobs?.running
    ));

    const {
      pending: removing,
      handler: removeExecutionAction,
    } = usePendingHandler(() => readRemediationInstructionExecution({ id: props.execution._id }));

    const removeExecution = async () => {
      await removeExecutionAction({ id: props.execution._id });

      emit('refresh');
    };

    const { showExecutionModal } = useShowExecutionModal(toRef(props, 'execution'));

    return {
      timeToComplete,
      stepFullNumber,
      hasJobs,
      hasActions,
      showExecutionModal,
      removing,
      removeExecution,
    };
  },
};
</script>
