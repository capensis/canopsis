import { storeToRefs } from 'pinia'

import { useAlarmStore } from '@/stores/modules/alarm'

/**
 * Alarm store facade.
 *
 * @returns {Object}
 */
export const useAlarm = () => {
  const alarmStore = useAlarmStore()
  const refs = storeToRefs(alarmStore)

  return {
    ...refs,
    ...alarmStore,
  }
}
