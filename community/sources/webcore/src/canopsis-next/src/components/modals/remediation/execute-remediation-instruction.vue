<template>
  <v-form @submit.prevent>
    <c-modal-wrapper close>
      <template #title>
        <span>{{ config.assignedInstruction.name }}</span>
      </template>
      <template #text>
        <v-fade-transition>
          <remediation-instruction-execute
            v-if="instructionExecution"
            :instruction-execution="instructionExecution"
            :next-pending="nextPending"
            :previous-pending="previousPending"
            @next-step="nextStep"
            @next-operation="nextOperation"
            @previous-operation="previousOperation"
            @execute-job="executeJob"
            @cancel-job-execution="cancelJobExecution"
          />
          <c-layout
            v-else
            justify="center"
          >
            <v-progress-circular
              color="primary"
              indeterminate
            />
          </c-layout>
        </v-fade-transition>
      </template>
      <template v-if="hasActions" #actions>
        <remediation-instruction-execute-btns
          :instruction-execution="instructionExecution"
          @close="closeModal"
        />
      </template>
    </c-modal-wrapper>
  </v-form>
</template>

<script>
import { pick } from 'lodash';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { MODALS, RESPONSE_STATUSES } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import {
  isInstructionExecutionAborted,
  isInstructionExecutionCompleted,
  isInstructionExecutionFailed,
  isInstructionExecutionPaused,
  isInstructionExecutionRunning,
} from '@/helpers/entities/remediation/instruction-execution/form';
import { getEmptyRemediationJobExecution } from '@/helpers/entities/remediation/job/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useSocket } from '@/hooks/socket';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';
import { useRemediationJobExecution } from '@/hooks/store/modules/remediation-job-execution';

import RemediationInstructionExecuteBtns from '@/components/other/remediation/instruction-execute/partials/remediation-instruction-execute-btns.vue';
import RemediationInstructionExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-execute.vue';

export default {
  name: MODALS.executeRemediationInstruction,

  components: {
    RemediationInstructionExecute,
    RemediationInstructionExecuteBtns,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const popups = usePopups();
    const socket = useSocket();
    const { config, close } = useInnerModal(props);

    const nextPending = ref(false);
    const previousPending = ref(false);
    const instructionExecution = ref(null);

    const {
      fetchRemediationInstructionExecutionWithoutStore,
      createRemediationInstructionExecution,
      nextRemediationInstructionExecutionOperation,
      nextRemediationInstructionExecutionStep,
      previousRemediationInstructionExecutionOperation,
      resumeRemediationInstructionExecution,
    } = useRemediationInstructionExecution();
    const {
      createRemediationJobExecution,
      cancelRemediationJobExecution,
    } = useRemediationJobExecution();

    const instruction = computed(() => config.value.assignedInstruction);
    const instructionId = computed(() => instruction.value?._id);
    const instructionExecutionId = computed(() => {
      const { execution } = instruction.value ?? {};

      return execution?._id ?? instructionExecution.value?._id;
    });
    const socketRoomName = computed(() => `${SOCKET_ROOMS.execution}/${instructionExecutionId.value}`);
    const hasActions = computed(
      () => instructionExecution.value && isInstructionExecutionRunning(instructionExecution.value),
    );

    const closeModal = () => {
      if (config.value.onClose) {
        config.value.onClose();
      }

      close();
    };

    const leaveFromSocketRoom = () => {
      if (instructionExecutionId.value && socket) {
        socket
          .off(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
          .off(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
          .leave(socketRoomName.value)
          .removeListener(setOperation);
      }
    };

    const socketCloseHandler = () => {
      if (!socket?.isConnectionOpen) {
        closeModal();
        popups.error({
          text: t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    };

    const socketCloseRoomHandler = () => {
      closeModal();
      popups.error({
        text: t('remediation.instructionExecute.popups.wasAborted', {
          instructionName: instructionExecution.value?.name,
        }),
        autoClose: false,
      });
    };

    const setOperation = (operation) => {
      if (!operation) {
        close();
        leaveFromSocketRoom();

        return;
      }

      for (const step of instructionExecution.value.steps) {
        const operationIndex = step.operations
          .findIndex(({ operation_id: operationId }) => operationId === operation.operation_id);

        if (operationIndex !== -1) {
          step.operations[operationIndex] = operation;

          return;
        }
      }
    };

    const joinToSocketRoom = () => {
      if (
        !instructionExecutionId.value
        || isInstructionExecutionFailed(instructionExecution.value)
        || isInstructionExecutionAborted(instructionExecution.value)
        || isInstructionExecutionCompleted(instructionExecution.value)
      ) {
        return;
      }

      socket
        ?.on(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
        .on(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
        .join(socketRoomName.value)
        .addListener(setOperation);
    };

    const tryToCallWithHandleNotFound = async (func) => {
      try {
        await func();
      } catch (err) {
        const message = err?.status === RESPONSE_STATUSES.notFound
          ? t('remediation.instructionExecute.popups.wasRemovedOrDisabled', {
            instructionName: instructionExecution.value?.name,
          })
          : err?.data?.error;

        console.error(err);

        popups.error({ text: message || t('errors.default') });

        closeModal();
      }
    };

    const callFunctionWithPreviousPending = async (func) => {
      previousPending.value = true;

      await tryToCallWithHandleNotFound(func);

      previousPending.value = false;
    };

    const callFunctionWithNextPending = async (func) => {
      nextPending.value = true;

      await tryToCallWithHandleNotFound(func);

      nextPending.value = false;
    };

    const nextStep = (success = false) => callFunctionWithNextPending(async () => {
      instructionExecution.value = await nextRemediationInstructionExecutionStep({
        id: instructionExecutionId.value,
        data: { failed: !success },
      });
    });

    const nextOperation = () => callFunctionWithNextPending(async () => {
      instructionExecution.value = await nextRemediationInstructionExecutionOperation({
        id: instructionExecutionId.value,
      });
    });

    const previousOperation = () => callFunctionWithPreviousPending(async () => {
      instructionExecution.value = await previousRemediationInstructionExecutionOperation({
        id: instructionExecutionId.value,
      });
    });

    const createExecution = () => tryToCallWithHandleNotFound(async () => {
      instructionExecution.value = await createRemediationInstructionExecution({
        data: {
          alarm: config.value.alarmId,
          instruction: instructionId.value,
        },
      });
    });

    const fetchExecution = () => tryToCallWithHandleNotFound(async () => {
      instructionExecution.value = await fetchRemediationInstructionExecutionWithoutStore({
        id: instructionExecutionId.value,
      });
    });

    const resumeExecution = () => tryToCallWithHandleNotFound(async () => {
      instructionExecution.value = await resumeRemediationInstructionExecution({
        id: instructionExecutionId.value,
      });
    });

    const setJob = (job, operation) => {
      setOperation({
        ...operation,

        jobs: operation.jobs.map(operationJob => (
          operationJob.job_id === job.job_id
            ? job
            : operationJob
        )),
      });
    };

    const executeJob = async ({ job, operation }) => {
      try {
        const updatedJob = await createRemediationJobExecution({
          data: {
            execution: instructionExecutionId.value,
            job: job.job_id,
            operation: operation.operation_id,
          },
        });

        setJob(updatedJob, operation);
      } catch (err) {
        console.error(err);

        popups.error({ text: err.error || t('errors.default') });
      }
    };

    const cancelJobExecution = async ({ job, operation }) => {
      try {
        await cancelRemediationJobExecution({ id: job._id });

        const updatedJob = {
          ...getEmptyRemediationJobExecution(),
          ...pick(job, ['_id', 'job_id', 'name', 'payload', 'query']),
        };

        setJob(updatedJob, operation);
      } catch (err) {
        console.error(err);

        popups.error({ text: err.error || t('errors.default') });
      }
    };

    const fetchInstructionExecution = async () => {
      const { execution } = config.value.assignedInstruction;

      if (!execution) {
        await createExecution();
      } else if (isInstructionExecutionPaused(execution)) {
        await resumeExecution();
      } else {
        await fetchExecution();
      }

      if (config.value.onExecute) {
        await config.value.onExecute();
      }
    };

    watch(instructionExecution, async (execution) => {
      if (!execution) {
        return;
      }

      if (!isInstructionExecutionRunning(execution)) {
        const isFailedExecution = isInstructionExecutionFailed(execution)
          || isInstructionExecutionAborted(execution);

        if (isFailedExecution) {
          popups.error({
            text: t('remediation.instructionExecute.popups.failed', {
              instructionName: execution.name,
            }),
          });
        }

        if (config.value.onComplete) {
          await config.value.onComplete(execution);
        }

        await closeModal();
      }
    });

    onMounted(async () => {
      await fetchInstructionExecution();

      joinToSocketRoom();
    });

    onBeforeUnmount(leaveFromSocketRoom);

    return {
      config,
      instructionExecution,
      nextPending,
      previousPending,
      hasActions,
      nextStep,
      nextOperation,
      previousOperation,
      executeJob,
      cancelJobExecution,
      closeModal,
    };
  },
};
</script>
