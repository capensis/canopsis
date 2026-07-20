import { storeToRefs } from 'pinia'

import { useInfosStore } from '@/stores/modules/infos'

/**
 * Infos store facade.
 *
 * @returns {Object}
 */
export const useInfos = () => {
  const infosStore = useInfosStore()
  const { alarmInfos, alarmInfosRules, entityInfos, pending } = storeToRefs(infosStore)

  return {
    alarmInfos,
    alarmInfosRules,
    entityInfos,
    infosPending: pending,
    fetchInfos: infosStore.fetch,
  }
}
