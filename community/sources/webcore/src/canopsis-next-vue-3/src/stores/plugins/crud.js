import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import request from '@/services/request'

export const DEFAULT_ENTITY_STORE_TYPES = {
  fetchList: 'fetchList',
  fetchListWithPreviousParams: 'fetchListWithPreviousParams',
  fetchListWithoutStore: 'fetchListWithoutStore',
  create: 'create',
  update: 'update',
  remove: 'remove',
}

export const createBasicCRUDActions = ({
  route,
  withWithoutStore = false,
} = {}) => {
  const actions = {
    create({ data } = {}) {
      return request.post(route, data)
    },

    update({ id, data } = {}) {
      return request.put(`${route}/${encodeURIComponent(id)}`, data)
    },

    remove({ id } = {}) {
      return request.delete(`${route}/${encodeURIComponent(id)}`)
    },
  }

  if (withWithoutStore) {
    actions.fetchListWithoutStore = ({ params } = {}) => request.get(route, { params })
  }

  return actions
}

export const createCRUDStore = ({
  id,
  route,
  dataPreparer = response => response?.data,
  metaPreparer = response => response?.meta,
  withFetchingParams = false,
  withMeta = false,
  withWithoutStore = false,
  bulkRoute,
  actions: additionalActions = {},
} = {}) => defineStore(id, () => {
  const items = ref([])
  const meta = ref({})
  const pending = ref(false)
  const fetchingParams = ref()

  const hasItems = computed(() => items.value.length > 0)

  const fetchList = async ({ params } = {}) => {
    try {
      pending.value = true

      if (withFetchingParams) {
        fetchingParams.value = params
      }

      const response = await request.get(route, { params })

      items.value = dataPreparer(response) ?? []

      if (withMeta) {
        meta.value = metaPreparer(response) ?? {}
      }

      return response
    } catch (err) {
      console.error(err)

      throw err
    } finally {
      pending.value = false
    }
  }

  const fetchListWithPreviousParams = () => fetchList({ params: fetchingParams.value })

  const fetchListWithoutStore = options => request.get(route, options)

  const bulkRemove = ({ data } = {}) => request.delete(bulkRoute, { data })

  return {
    items,
    meta,
    pending,
    fetchingParams,
    hasItems,

    fetchList,
    ...(withFetchingParams ? { fetchListWithPreviousParams } : {}),
    ...(withWithoutStore ? { fetchListWithoutStore } : {}),
    ...(bulkRoute ? { bulkRemove } : {}),

    ...createBasicCRUDActions({ route }),
    ...additionalActions,
  }
})
