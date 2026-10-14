<template>
  <modal-wrapper close>
    <template #title="">
      <span>{{ config.assignedInstruction.name }}</span>
    </template>
    <template #text="">
      <v-fade-transition>
        <v-layout
          v-if="pending"
          justify-center
        >
          <v-progress-circular
            color="primary"
            indeterminate
          />
        </v-layout>
        <remediation-instruction-simple-execute
          v-else
          :executing="executing"
          :instruction-execution="instructionExecution"
          :jobs="jobs"
          @run:jobs="runJobs"
        />
      </v-fade-transition>
    </template>
  </modal-wrapper>
</template>

<script>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { SOCKET_ROOMS } from '@/config';
import { MODALS } from '@/constants';

import Socket from '@/plugins/socket/services/socket';

import {
  isInstructionExecutionRunning,
  isInstructionExecutionAborted,
  isInstructionExecutionFinished,
} from '@/helpers/entities/remediation/instruction-execution/form';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useSocket } from '@/hooks/socket';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';
import { useStoreModuleHooks } from '@/hooks/store';

import RemediationInstructionSimpleExecute from '@/components/other/remediation/instruction-execute/remediation-instruction-simple-execute.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.executeRemediationSimpleInstruction,
  components: {
    RemediationInstructionSimpleExecute,
    ModalWrapper,
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
    const modals = useModals();
    const popups = usePopups();
    const socket = useSocket();

    const config = computed(() => props.modal?.config ?? {});

    const { useActions: useInstructionActions } = useStoreModuleHooks('remediationInstruction');
    const { fetchRemediationInstructionWithoutStore } = useInstructionActions({
      fetchRemediationInstructionWithoutStore: 'fetchItemWithoutStore',
    });

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

    const {
      fetchRemediationInstructionExecutionWithoutStore,
      createRemediationInstructionExecution,
    } = useRemediationInstructionExecution();

    const instructionExecutionId = computed(() => {
      const { execution } = config.value.assignedInstruction ?? {};

      return execution?._id ?? instructionExecution.value?._id;
    });

    const instructionJobs = computed(() => instruction.value?.jobs?.map(({ job }) => job));

    const jobs = computed(() => instructionExecution.value?.jobs ?? instructionJobs.value);

    const socketRoomName = computed(() => `${SOCKET_ROOMS.execution}/${instructionExecutionId.value}`);

    /**
     * Invokes the modal `onClose` callback when present and hides this modal in the modals store.
     */
    const closeModal = () => {
      if (config.value.onClose) {
        config.value.onClose();
      }

      modals.hide({ id: props.modal.id });
    };

    /**
     * Handles socket `customClose`: when the WebSocket is not open, closes the modal and shows a connection error.
     */
    const socketCloseHandler = () => {
      if (!socket.isConnectionOpen) {
        closeModal();
        popups.error({
          text: t('remediation.instructionExecute.popups.connectionError'),
          autoClose: false,
        });
      }
    };

    /**
     * Handles socket `closeRoom`: when the instruction execution was aborted, closes the modal and notifies the user.
     */
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

    /**
     * Loads the remediation instruction execution by id, updates local state and `executing`, or closes the modal on
     * error.
     */
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

    /**
     * Refreshes execution from the API then applies job updates received from the execution socket channel.
     *
     * @param {Array} jobsPayload - Jobs list pushed by the socket listener.
     */
    const setJobs = async (jobsPayload) => {
      await fetchInstructionExecution();

      instructionExecution.value.jobs = jobsPayload;
    };

    /**
     * Joins the execution socket room for the current execution id, registers connection handlers, and attaches the
     * jobs listener on the room.
     */
    const joinToSocketRoom = () => {
      socket
        .on(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
        .join(socketRoomName.value)
        .addListener(setJobs);

      if (!isInstructionExecutionFinished(instructionExecution.value)) {
        socket.on(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler);
      }
    };

    /**
     * Leaves the execution socket room, removes the jobs listener, and unregisters connection-level handlers.
     */
    const leaveFromSocketRoom = () => {
      socket
        .off(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
        .off(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
        .leave(socketRoomName.value)
        .removeListener(setJobs);
    };

    /**
     * Creates a new instruction execution for the configured alarm and instruction, then runs optional `onExecute`.
     * On failure, shows an error popup and closes the modal.
     */
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

    /**
     * Starts execution via the API then joins the execution socket room for live job updates.
     */
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
