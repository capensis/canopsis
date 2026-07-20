import { storeToRefs } from 'pinia'

import { useServiceEntityStore } from '@/stores/modules/service-entity'

/**
 * Service entity store facade.
 *
 * @returns {Object}
 */
export const useServiceEntity = () => {
  const serviceEntityStore = useServiceEntityStore()
  const { services, getListByServiceId, getPendingByServiceId, getMetaByServiceId } =
    storeToRefs(serviceEntityStore)

  return {
    serviceEntitiesByService: services,
    getServiceEntitiesListByServiceId: getListByServiceId,
    getServiceEntitiesPendingByServiceId: getPendingByServiceId,
    getServiceEntitiesMetaByServiceId: getMetaByServiceId,
    fetchServiceEntitiesList: serviceEntityStore.fetchList,
  }
}
