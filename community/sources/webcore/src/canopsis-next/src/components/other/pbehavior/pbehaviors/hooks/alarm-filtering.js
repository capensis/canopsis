import { ref } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useSocket } from '@/hooks/socket';
import { usePbehaviorPatterns } from '@/hooks/store/modules/pbehavior-patterns';

/**
 * Custom hook for managing pbehavior alarm filtering operations.
 *
 * Provides a clean interface for running alarm filtering with automatic
 * socket management, loading states, and error handling.
 *
 * @param {Function} refreshHandler - Function to call when alarm filtering is completed
 *
 * @returns {Object} An object containing:
 * @property {import('vue').Ref<boolean>} pending - Reactive loading state
 * @property {Function} runAlarmFiltering - Function to trigger alarm filtering
 *
 * @example
 * // Usage in a component
 * const { pending, runAlarmFiltering } = useAlarmFiltering();
 *
 * // In template: :loading="pending"
 * // On click: @click="runAlarmFiltering"
 */
export const useAlarmFiltering = (refreshHandler = () => {}) => {
  const { t } = useI18n();
  const popups = usePopups();
  const socket = useSocket();
  const { runAlarmFiltering: runAlarmFilteringAction } = usePbehaviorPatterns();
  const pending = ref(false);

  /**
   * Runs the alarm filtering process with proper socket management.
   *
   * - Connects to 'pbehavior-patterns' websocket room
   * - Sends API request to trigger filtering
   * - Listens for completion message ({ ok: true })
   * - Handles cleanup on success/error
   */
  const runAlarmFiltering = async () => {
    pending.value = true;

    const completionListener = ({ ok } = {}) => {
      if (ok) {
        pending.value = false;

        socket.leave(SOCKET_ROOMS.pbehaviorPatterns).removeListener(completionListener);
        refreshHandler?.();
      }
    };

    try {
      socket.join(SOCKET_ROOMS.pbehaviorPatterns).addListener(completionListener);
      await runAlarmFilteringAction();
    } catch (err) {
      pending.value = false;

      socket.leave(SOCKET_ROOMS.pbehaviorPatterns).removeListener(completionListener);
      console.error('Error running alarm filtering:', err);
      popups.error({ text: t('pbehavior.alarmFilteringError') });
    }
  };

  return {
    pending,
    runAlarmFiltering,
  };
};
