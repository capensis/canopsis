import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useEventsRecordCurrent } from '@/hooks/store/modules/events-record-current';

/**
 * Hook to manage event recording functionality.
 *
 * @param {Function} [fetchListHandler = () => {}] - A function to fetch the list of events records.
 * @returns {Object} An object containing methods to start and stop event recording.
 */
export const useEventsRecordRecording = (fetchListHandler = () => {}) => {
  const { t } = useI18n();
  const modals = useModals();

  const { startEventsRecordCurrent, stopEventsRecordCurrent } = useEventsRecordCurrent();

  /**
   * Start recording events based on a specified event pattern.
   */
  const startRecording = () => modals.show({
    name: MODALS.applyEventFilter,
    config: {
      title: t('eventsRecord.launchEventRecording'),
      action: eventPattern => startEventsRecordCurrent({ data: { event_pattern: eventPattern } }),
      afterSubmit: fetchListHandler,
    },
  });

  /**
   * Stop the event recording by id.
   *
   * @param {string} [recordingId] - ID of the recording to stop. Required when multiple recordings exist.
   */
  const stopRecording = recordingId => modals.show({
    name: MODALS.confirmation,
    config: {
      action: async () => {
        await stopEventsRecordCurrent({ id: recordingId });

        return fetchListHandler();
      },
    },
  });

  return {
    startRecording,
    stopRecording,
  };
};
