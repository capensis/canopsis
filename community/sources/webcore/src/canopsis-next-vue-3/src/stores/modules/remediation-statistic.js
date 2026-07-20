import { ref } from 'vue'
import { defineStore } from 'pinia'

import { API_ROUTES } from '@/config'

import request from '@/services/request'

export const useRemediationStatisticStore = defineStore('remediationStatistic', () => {
  const metrics = ref([])
  const meta = ref({})
  const pending = ref(false)
  const fetchingParams = ref({})

  const fetchList = async ({ params } = {}) => {
    fetchingParams.value = params
    pending.value = true

    const { data = [], meta: responseMeta = {} } = await request.get(
      API_ROUTES.metrics.remediation,
      { params },
    )

    metrics.value = data
    meta.value = responseMeta
    pending.value = false
  }

  const fetchListWithPreviousParams = () => fetchList({ params: fetchingParams.value })

  return {
    metrics,
    meta,
    pending,
    fetchingParams,
    fetchList,
    fetchListWithPreviousParams,
  }
})
