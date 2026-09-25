import { computed, onBeforeUnmount, ref, unref, watch } from 'vue';

import Socket from '@/plugins/socket/services/socket';

import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { isWebhookExecutionFinished } from '@/helpers/entities/webhook-execution/entity';

import { useFormContext, useFormErrors } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useSocket } from '@/hooks/socket';
import { useValidationFormErrors } from '@/hooks/validator/validation-form-errors';

/**
 * Handles shared test execution flow for webhook-based rule forms.
 *
 * The hook keeps selected alarm, pending state, current execution status,
 * validation errors and socket room lifecycle. Domain-specific components
 * provide only their socket room, API actions and payload mapper.
 *
 * @param {Object} options Hook options.
 * @param {import('vue').Ref<Object>|Object} options.form Reactive rule form used for validation
 * and payload mapping.
 * @param {string} options.socketRoom Base socket room name for the execution.
 * @param {Function} options.createExecution Action that creates a test execution.
 * @param {Function} options.fetchExecution Action that fetches execution status by id after
 * socket error.
 * @param {Function} options.getExecutionData Maps `{ alarm, form }` to API request payload.
 * @param {Function} [options.errorsToValidation] Maps API error paths to form field paths.
 * @returns {{
 *   alarm: import('vue').Ref<string>,
 *   pending: import('vue').Ref<boolean>,
 *   executionStatus: import('vue').Ref<Object|undefined|null>,
 *   hasErrors: import('vue').ComputedRef<boolean>,
 *   alarmsPatternsParams: import('vue').ComputedRef<Object<string, string>>,
 *   clearWebhookStatus: Function,
 *   runTestExecution: Function,
 * }}
 */
export const useAlarmWebhookTestQuery = ({
  form,
  socketRoom,
  createExecution,
  fetchExecution,
  getExecutionData,
  errorsToValidation = errors => errors,
}) => {
  const { t } = useI18n();
  const popups = usePopups();
  const modals = useModals();
  const socket = useSocket();
  const formContext = useFormContext();
  const formErrors = useFormErrors();
  const { setFormErrors } = useValidationFormErrors(unref(form));

  const alarm = ref('');
  const pending = ref(false);
  const executionStatus = ref(undefined);

  const hasErrors = computed(() => Object.keys(formErrors.value).length > 0);
  const alarmsPatternsParams = computed(() =>
    Object.entries(formFilterToPatterns(unref(form).patterns)).reduce((acc, [key, value]) => {
      acc[key] = JSON.stringify(value);

      return acc;
    }, {}),
  );

  const getSocketRoomName = id => `${socketRoom}/${id}`;

  const setExecutionStatus = (nextExecutionStatus) => {
    executionStatus.value = nextExecutionStatus;
  };

  const clearWebhookStatus = () => {
    executionStatus.value = null;
  };

  const socketCloseHandler = () => {
    if (!socket?.isConnectionOpen) {
      modals.hide();
      popups.error({
        text: t('remediation.instructionExecute.popups.connectionError'),
        autoClose: false,
      });
    }
  };

  const socketCloseRoomHandler = () => {
    modals.hide();
  };

  const handleSocketError = async () => {
    executionStatus.value = await fetchExecution({
      id: executionStatus.value._id,
    });
  };

  const joinToSocketRoom = () => {
    socket
      ?.on(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
      .on(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
      .on(Socket.EVENTS_TYPES.error, handleSocketError)
      .join(getSocketRoomName(executionStatus.value._id))
      .addListener(setExecutionStatus);
  };

  const leaveFromSocketRoom = () => {
    socket
      ?.off(Socket.EVENTS_TYPES.customClose, socketCloseHandler)
      .off(Socket.EVENTS_TYPES.closeRoom, socketCloseRoomHandler)
      .off(Socket.EVENTS_TYPES.error, handleSocketError)
      .leave(getSocketRoomName(executionStatus.value._id))
      .removeListener(setExecutionStatus);
  };

  const runTestExecution = async () => {
    const { valid: isFormValid } = await formContext.validate();

    if (!isFormValid) {
      return;
    }

    pending.value = true;
    clearWebhookStatus();

    try {
      executionStatus.value = await createExecution({
        data: getExecutionData({
          alarm: alarm.value,
          form: unref(form),
        }),
      });

      joinToSocketRoom();
    } catch (err) {
      if (err.error) {
        popups.error({ text: err.error });
      } else {
        setFormErrors(errorsToValidation(err));
      }

      executionStatus.value = undefined;
    } finally {
      pending.value = false;
    }
  };

  watch(executionStatus, (nextExecutionStatus) => {
    if (isWebhookExecutionFinished(nextExecutionStatus)) {
      leaveFromSocketRoom();
    }
  });

  onBeforeUnmount(() => {
    if (executionStatus.value) {
      leaveFromSocketRoom();
    }
  });

  return {
    alarm,
    pending,
    executionStatus,
    hasErrors,
    alarmsPatternsParams,
    clearWebhookStatus,
    runTestExecution,
  };
};
