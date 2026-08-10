import { MODALS, SINGLE_EVENT_RESEND_DELAY, TIME_UNITS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useEventsRecord } from '@/hooks/store/modules/events-record';

export const useEventsRecordResending = () => {
  const { t } = useI18n();
  const modals = useModals();
  const popups = usePopups();

  const { playbackEventsRecord, stopEventsRecordPlayback, fetchCurrentEventsRecord } =
    useEventsRecord();

  const resend = async (eventsRecordId, eventIds = [], delay = SINGLE_EVENT_RESEND_DELAY) => {
    await playbackEventsRecord({
      id: eventsRecordId,
      data: { delay, event_ids: eventIds },
    });

    popups.info({ text: t('eventsRecord.selectedEventsAreResent') });

    return fetchCurrentEventsRecord();
  };

  const startResending = (eventsRecordId, eventIds = []) => {
    if (eventIds.length === 1) {
      return resend(eventsRecordId, eventIds);
    }

    return modals.show({
      name: MODALS.duration,
      config: {
        title: t('eventsRecord.resendEvents', 1),
        label: t('eventsRecord.delayBetweenEvents'),
        units: [
          { value: TIME_UNITS.millisecond, text: 'common.times.millisecond' },
          { value: TIME_UNITS.second, text: 'common.times.second' },
        ],
        action: delay => resend(eventsRecordId, eventIds, delay),
      },
    });
  };

  const stopResending = eventsRecordId =>
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await stopEventsRecordPlayback({ id: eventsRecordId });

          return fetchCurrentEventsRecord();
        },
      },
    });

  return {
    startResending,
    stopResending,
  };
};
