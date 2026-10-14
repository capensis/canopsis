<template>
  <c-modal-wrapper close>
    <template #title>
      <span>{{ config.assignedInstruction.name }}</span>
    </template>
    <template #text>
      <v-fade-transition>
        <c-layout
          v-if="pending"
          justify="center"
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </c-layout>
        <remediation-instruction-simple-execute
          v-else
          :executing="executing"
          :instruction-execution="instructionExecution"
          :jobs="jobs"
          @run:jobs="runJobs"
        />
      </v-fade-transition>
    </template>
  </c-modal-wrapper>
</template>

<script>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { MODALS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import {
  isInstructionExecutionAborted,
  isInstructionExecutionFinished,
  isInstructionExecutionRunning,
} from '@/helpers/entities/remediation/instruction-execution/form';

import { useI18n } from '@/hooks/i18n';
import { useInnerModal } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useSocket } from '@/hooks/socket';
import { useRemediationInstruction } from '@/hooks/store/modules/remediation-instruction';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

import RemediationInstructionSimpleExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-simple-execute.vue';

export default {
  name: MODALS.executeRemediationSimpleInstruction,

  components: {
    RemediationInstructionSimpleExecute,
  },

  props: {
    modal: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const executing = ref(false);
    const instruction = ref(null);
    const instructionExecution = ref(null);

    const { t } = useI18n();
    const popups = usePopups();
    const socket = useSocket();
    const { config, close } = useInnerModal(props);

    const { fetchRemediationInstructionWithoutStore } = useRemediationInstruction();
    const {
      fetchRemediationInstructionExecutionWithoutStore,
      createRemediationInstructionExecution,
    } = useRemediationInstructionExecution();

    const instructionId = computed(() => config.value.assignedInstruction?._id);

    const {
      pending,
      handler: fetchInstruction,
    } = usePendingHandler(async () => {
      try {
        instruction.value = await fetchRemediationInstructionWithoutStore({ id: instructionId.value });
      } catch (err) {
        console.warn(err);
      }
    }, true);

    const instructionExecutionId = computed(() => {
      const { execution } = config.value.assignedInstruction ?? {};

      return execution?._id ?? instructionExecution.value?._id;
    });

    const instructionJobs = computed(() => instruction.value?.jobs?.map(({ job }) => job));
    const jobs = computed(() => instructionExecution.value?.jobs ?? instructionJobs.value);
    const socketRoomName = computed(() => `${SOCKET_ROOMS.execution}/${instructionExecutionId.value}`);

    const closeModal = () => {
      if (config.value.onClose) {
        config.value.onClose();
      }

      close();
    };

    const socketCloseHandler = () => {
      if (!socket.isConnectionOpen) {
        closeModal();
        popups.error({
          text: t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    };

    const socketCloseRoomHandler = () => {
      if (!isInstructionExecutionAborted(instructionExecution.value)) {
        return;
      }

      closeModal();

      popups.error({
        text: t('remediation.instructionExecute.popups.wasAborted', {
          instructionName: instructionExecution.value?.name,
        }),
        autoClose: false,
      });
    };

    const fetchInstructionExecution = async () => {
      try {
        if (instructionExecutionId.value) {
          instructionExecution.value = await fetchRemediationInstructionExecutionWithoutStore({
            id: instructionExecutionId.value,
          });

          const newExecuting = isInstructionExecutionRunning(instructionExecution.value);

          if (newExecuting !== executing.value) {
            executing.value = newExecuting;
          }
        }
      } catch (err) {
        console.error(err);

        popups.error({ text: err.error || t('errors.default') });

        closeModal();
      }
    };

    const setJobs = async (jobsPayload) => {
      await fetchInstructionExecution();

      instructionExecution.value.jobs = jobsPayload;
    };

    const joinToSocketRoom = () => {
      socket
        .on(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
        .join(socketRoomName.value)
        .addListener(setJobs);

      if (!isInstructionExecutionFinished(instructionExecution.value)) {
        socket.on(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler);
      }
    };

    const leaveFromSocketRoom = () => {
      socket
        .off(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
        .off(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
        .leave(socketRoomName.value)
        .removeListener(setJobs);
    };

    const createInstructionExecution = async () => {
      try {
        instructionExecution.value = await createRemediationInstructionExecution({
          data: {
            alarm: config.value.alarmId,
            instruction: instructionId.value,
          },
        });

        executing.value = true;

        if (config.value.onExecute) {
          await config.value.onExecute();
        }
      } catch (err) {
        console.error(err);
        popups.error({ text: err.error || t('errors.default') });

        closeModal();
      }
    };

    const runJobs = async () => {
      await createInstructionExecution();

      joinToSocketRoom();
    };

    onMounted(async () => {
      fetchInstruction();

      if (config.value.assignedInstruction.execution) {
        await fetchInstructionExecution();

        joinToSocketRoom();
      }
    });

    onBeforeUnmount(leaveFromSocketRoom);

    return {
      pending,
      executing,
      instructionExecution,
      jobs,
      config,
      runJobs,
    };
  },
};
</script>
