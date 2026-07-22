import { MODALS } from '@/constants'

import { useI18n } from '@/hooks/i18n'
import { useModals } from '@/hooks/modals'
import { useEventsRecord } from '@/hooks/store/modules/events-record'

export const useEventsRecordRecording = (fetchListHandler = () => {}) => {
  const { t } = useI18n()
  const modals = useModals()

  const { startCurrentEventsRecord, stopCurrentEventsRecord } = useEventsRecord()

  const startRecording = () =>
    modals.show({
      name: MODALS.applyEventFilter,
      config: {
        title: t('eventsRecord.launchEventRecording'),
        required: false,
        action: (eventPattern) =>
          startCurrentEventsRecord({ data: { event_pattern: eventPattern } }),
        afterSubmit: fetchListHandler,
      },
    })

  const stopRecording = (recordingId) =>
    modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await stopCurrentEventsRecord({ id: recordingId })

          return fetchListHandler()
        },
      },
    })

  return {
    startRecording,
    stopRecording,
  }
}
