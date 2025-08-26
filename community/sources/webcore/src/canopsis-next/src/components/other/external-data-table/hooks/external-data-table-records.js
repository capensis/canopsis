import { omit } from 'lodash';
import { computed, unref, inject, onMounted } from 'vue';

import { MODALS } from '@/constants';

import { pickIds } from '@/helpers/array';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useExternalDataTableRecord } from '@/hooks/store/modules/external-data-table-record';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

/**
 * Provides modal management functions for external data table records.
 *
 * @param {Object} options - Options for configuring the modals.
 * @param {Object} options.externalDataTable - The external data table instance.
 * @param {Function} options.fetchList - Function to fetch the list of records.
 * @returns {Object} An object containing functions to show various modals.
 * @property {Function} showImportExternalDataTablesModal - Shows the modal for importing records.
 * @property {Function} showCreateExternalDataTableRecordModal - Shows the modal for creating a new record.
 * @property {Function} showEditExternalDataTableRecordModal - Shows the modal for editing an existing record.
 * @property {Function} showDuplicateExternalDataTableRecordModal - Shows the modal for duplicating a record.
 * @property {Function} showRemoveExternalDataTableRecordModal - Shows the modal for removing a record.
 * @property {Function} showRemoveSelectedExternalDataTableRecordsModal - Shows the modal for removing selected records.
 */
export const useExternalDataTableRecordsModals = ({
  externalDataTable,
  fetchList,
}) => {
  const { t } = useI18n();
  const modals = useModals();

  const refreshObserver = inject('$refresh', null);

  const {
    createExternalDataTableRecord,
    updateExternalDataTableRecord,
    removeExternalDataTableRecord,
    bulkRemoveExternalDataTableRecord,
  } = useExternalDataTableRecord();

  const { callActionWithPopup } = useCallActionWithPopup();

  /**
   * Shows the modal for importing external data table records.
   */
  const showImportExternalDataTablesModal = () => modals.show({
    name: MODALS.importExternalDataTableRecords,
    config: {
      externalDataTable: unref(externalDataTable),
      afterSubmit: () => callActionWithPopup(() => (refreshObserver?.notify ? refreshObserver.notify() : fetchList())),
    },
  });

  /**
   * Shows the modal for creating a new external data table record.
   */
  const showCreateExternalDataTableRecordModal = () => modals.show({
    name: MODALS.createExternalDataTableRecord,
    config: {
      externalDataTable: unref(externalDataTable),
      action: newExternalDataTableRecord => callActionWithPopup(
        () => createExternalDataTableRecord({
          table: unref(externalDataTable)?._id,
          data: newExternalDataTableRecord,
        }),
        fetchList,
      ),
    },
  });

  /**
   * Shows the modal for editing an existing external data table record.
   *
   * @param {Object} externalDataTableRecord - The record to be edited.
   */
  const showEditExternalDataTableRecordModal = externalDataTableRecord => modals.show({
    name: MODALS.createExternalDataTableRecord,
    config: {
      externalDataTableRecord,
      externalDataTable: unref(externalDataTable),
      title: t('modals.createExternalDataTableRecord.edit.title'),
      action: newExternalDataTableRecord => callActionWithPopup(
        () => updateExternalDataTableRecord({
          id: externalDataTableRecord._id,
          table: unref(externalDataTable)?._id,
          data: newExternalDataTableRecord,
        }),
        fetchList,
      ),
    },
  });

  /**
   * Shows the modal for duplicating an external data table record.
   *
   * @param {Object} externalDataTableRecord - The record to be duplicated.
   */
  const showDuplicateExternalDataTableRecordModal = externalDataTableRecord => modals.show({
    name: MODALS.createExternalDataTableRecord,
    config: {
      externalDataTableRecord: omit(externalDataTableRecord, ['_id']),
      externalDataTable: unref(externalDataTable),
      title: t('modals.createExternalDataTableRecord.duplicate.title'),
      action: newExternalDataTableRecord => callActionWithPopup(
        () => createExternalDataTableRecord({
          table: unref(externalDataTable)?._id,
          data: newExternalDataTableRecord,
        }),
        fetchList,
      ),
    },
  });

  /**
   * Shows the modal for removing an external data table record.
   *
   * @param {string} id - The ID of the record to be removed.
   */
  const showRemoveExternalDataTableRecordModal = id => modals.show({
    name: MODALS.confirmation,
    config: {
      action: () => callActionWithPopup(
        () => removeExternalDataTableRecord({
          id,
          table: unref(externalDataTable)?._id,
        }),
        fetchList,
      ),
    },
  });

  /**
   * Shows the modal for removing selected external data table records.
   *
   * @param {Object[]} externalDataTableRecords - The external data table records to be removed.
   */
  const showRemoveSelectedExternalDataTableRecordsModal = (externalDataTableRecords = []) => modals.show({
    name: MODALS.confirmation,
    config: {
      action: () => callActionWithPopup(
        () => bulkRemoveExternalDataTableRecord({
          data: pickIds(externalDataTableRecords),
          table: unref(externalDataTable)?._id,
        }),
        fetchList,
      ),
    },
  });

  return {
    showImportExternalDataTablesModal,
    showCreateExternalDataTableRecordModal,
    showEditExternalDataTableRecordModal,
    showDuplicateExternalDataTableRecordModal,
    showRemoveExternalDataTableRecordModal,
    showRemoveSelectedExternalDataTableRecordsModal,
  };
};

/**
 * Hook that provides functionality for managing external data table records
 *
 * @param {Object} options - Configuration options
 * @param {Object|Ref<Object>} [options.externalDataTable] - The external data table object
 * @param {Function} [options.fetchExternalDataTable] - The fetch external data table function
 * @param {Object|Ref<Object>} [options.initialQuery] - The initial qurty parameter
 * @returns {Object} Object containing records state and methods
 * @property {Ref<Array>} records - Reactive reference to the fetched records
 * @property {Ref<Object>} meta - Metadata associated with the fetched records
 * @property {Ref<boolean>} pending - Whether records are currently being fetched
 * @property {Ref<Object>} query - Current query state for fetching records
 * @property {ComputedRef<Object>} options - Computed reference to pagination and sorting options
 * @property {Function} updateOptions - Function to update pagination and sorting options
 * @property {ComputedRef<boolean>} hasStructure - Whether the table has a defined column structure
 * @property {Function} fetchList - Function to fetch the records list with current query
 * @property {Function} showImportExternalDataTablesModal - Shows the modal for importing records
 * @property {Function} showCreateExternalDataTableRecordModal - Shows the modal for creating a new record
 * @property {Function} showEditExternalDataTableRecordModal - Shows the modal for editing an existing record
 * @property {Function} showDuplicateExternalDataTableRecordModal - Shows the modal for duplicating a record
 * @property {Function} showRemoveExternalDataTableRecordModal - Shows the modal for removing a record
 * @property {Function} showRemoveSelectedExternalDataTableRecordsModal - Shows the modal for removing selected records
 */
export const useExternalDataTableRecordsList = ({ externalDataTable, fetchExternalDataTable, initialQuery } = {}) => {
  const { fetchExternalDataTableRecordsListWithoutStore } = useExternalDataTableRecord();

  const {
    data: records,
    meta,
    pending,
    query,
    options,
    updateOptions,
    resetQuery,
    fetchList,
  } = useFetchListWithoutStoreWithOptions({
    initialQuery,
    fetchListHandler: (rest) => {
      fetchExternalDataTable?.();

      return fetchExternalDataTableRecordsListWithoutStore({
        ...rest,
        id: unref(externalDataTable)._id,
      });
    },
  });

  const hasStructure = computed(() => !!unref(externalDataTable).column_configs?.length);

  const {
    showImportExternalDataTablesModal,
    showCreateExternalDataTableRecordModal,
    showEditExternalDataTableRecordModal,
    showDuplicateExternalDataTableRecordModal,
    showRemoveExternalDataTableRecordModal,
    showRemoveSelectedExternalDataTableRecordsModal,
  } = useExternalDataTableRecordsModals({
    externalDataTable,
    fetchList,
  });

  onMounted(() => {
    if (hasStructure.value) {
      fetchList();
    }
  });

  return {
    records,
    meta,
    pending,
    query,
    options,
    updateOptions,
    resetQuery,

    hasStructure,

    fetchList,

    showImportExternalDataTablesModal,
    showCreateExternalDataTableRecordModal,
    showEditExternalDataTableRecordModal,
    showDuplicateExternalDataTableRecordModal,
    showRemoveExternalDataTableRecordModal,
    showRemoveSelectedExternalDataTableRecordsModal,
  };
};
