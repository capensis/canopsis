import { MODALS, SINGLE_EVENT_RESEND_DELAY, TIME_UNITS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEventsRecord } from '@/hooks/store/modules/events-record';
import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';

/**
 * Hook to manage resending events records.
 * This hook provides methods to start and stop resending events records, as well as handling the delay between events.
 *
 * @returns {Object} An object containing methods for managing resending events records.
 */
export const useEventsRecordResending = () => {
  const { t, tc } = useI18n();
  const modals = useModals();

  /**
   * STORE
   */
  const { playbackEventsRecordEvents, stopPlaybackEventsRecordEvents } = useEventsRecord();
  const { current, fetchEventsRecordCurrent } = useEventsRecordCurrent();

  /**
   * Resend events with a specified delay.
   *
   * @param {string} eventsRecordId - The ID of the events record.
   * @param {string[]} [eventIds = []] - An array of event IDs to resend.
   * @param {Duration} [delay = SINGLE_EVENT_RESEND_DELAY] - The delay between resending the events.
   * @returns {Promise}
   */
  const resend = async (eventsRecordId, eventIds = [], delay = SINGLE_EVENT_RESEND_DELAY) => {
    await playbackEventsRecordEvents({
      id: eventsRecordId,
      data: { delay, event_ids: eventIds },
    });

    return fetchEventsRecordCurrent();
  };

  /**
   * Start resending events with optional delay configuration.
   *
   * @param {string} eventsRecordId - The ID of the events record.
   * @param {string[]} [eventIds = []] - An array of event IDs to resend.
   * @returns {Promise}
   */
  const startResending = (eventsRecordId, eventIds = []) => {
    if (eventIds.length === 1) {
      return resend(eventsRecordId, eventIds);
    }

    return modals.show({
      name: MODALS.duration,
      config: {
        title: tc('eventsRecord.resendEvents', 1),
        label: t('eventsRecord.delayBetweenEvents'),
        units: [
          { value: TIME_UNITS.millisecond, text: 'common.times.millisecond' },
          { value: TIME_UNITS.second, text: 'common.times.second' },
        ],
        action: delay => resend(eventsRecordId, eventIds, delay),
      },
    });
  };

  /**
   * Stop resending events and confirm the action.
   *
   * @returns {Promise}
   */
  const stopResending = () => modals.show({
    name: MODALS.confirmation,
    config: {
      action: async () => {
        await stopPlaybackEventsRecordEvents({ id: current.value._id });

        return fetchEventsRecordCurrent();
      },
    },
  });

  return {
    startResending,
    stopResending,
  };
};
