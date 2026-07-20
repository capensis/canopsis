import { defineStore } from 'pinia'

import { API_ROUTES } from '@/config'

import request from '@/services/request'

export const useExternalDataTableStore = defineStore('externalDataTable', () => {
  const convertObjectToFormData = (data = {}) =>
    Object.entries(data).reduce((formData, [key, value]) => {
      formData.append(key, value)

      return formData
    }, new FormData())

  const fetchListWithoutStore = ({ params } = {}) =>
    request.get(API_ROUTES.externalDataTables, { params })

  const fetchWithoutStore = ({ id } = {}) => request.get(`${API_ROUTES.externalDataTables}/${id}`)

  const fetchSchema = ({ id } = {}) => request.get(`${API_ROUTES.externalDataTables}/${id}/schema`)

  const createImport = ({ id, data } = {}) =>
    request.post(`${API_ROUTES.externalDataImport}/${id}`, convertObjectToFormData(data), {
      headers: { 'Content-Type': 'multipart/form-data' },
    })

  const fetchImportData = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.externalDataImport}/${id}/data`, {
      params,
    })

  const fetchImportStatus = ({ id } = {}) =>
    request.get(`${API_ROUTES.externalDataImport}/${id}/status`)
  const previewImport = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.externalDataImport}/${id}/preview`, data)
  const completeImport = ({ id, data } = {}) =>
    request.put(`${API_ROUTES.externalDataImport}/${id}/complete`, data)

  const fetchRecordsListWithoutStore = ({ id, params } = {}) =>
    request.get(`${API_ROUTES.externalDataTables}/${id}/data`, { params })

  const createRecord = ({ table, data } = {}) =>
    request.post(`${API_ROUTES.externalDataTables}/${table}/data`, data)
  const updateRecord = ({ table, id, data } = {}) =>
    request.put(`${API_ROUTES.externalDataTables}/${table}/data/${id}`, data)
  const removeRecord = ({ table, id } = {}) =>
    request.delete(`${API_ROUTES.externalDataTables}/${table}/data/${id}`)
  const bulkRemoveRecords = ({ table, data } = {}) =>
    request.delete(`${API_ROUTES.bulkExternalDataTables}/${table}/data`, { data })

  return {
    fetchListWithoutStore,
    fetchWithoutStore,
    fetchSchema,
    createImport,
    fetchImportData,
    fetchImportStatus,
    previewImport,
    completeImport,
    fetchRecordsListWithoutStore,
    createRecord,
    updateRecord,
    removeRecord,
    bulkRemoveRecords,
  }
})
