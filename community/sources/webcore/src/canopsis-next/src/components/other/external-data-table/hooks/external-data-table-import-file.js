import { mapValues } from 'lodash';
import { computed, ref } from 'vue';

import { CSV_SEPARATORS, IMPORT_STATUSES, IMPORT_PREVIEW_STATUSES } from '@/constants';

import {
  externalDataTableColumnConfigsToForm,
  formToExternalDataTableColumnTags,
  formToExternalDataTableColumnConfigs,
} from '@/helpers/entities/external-data-table/form';

import { usePendingHandler } from '@/hooks/query/pending';
import { useValidator } from '@/hooks/validator/validator';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { usePolling, useFilePolling } from '@/hooks/polling';
import { useExternalDataTableImport } from '@/hooks/store/modules/external-data-table-import';

export const useExternalDataTableImportFile = ({ config }) => {
  const validator = useValidator();

  const {
    createExternalDataTableImport,
    fetchExternalDataTableImportStatus,
    fetchExternalDataTableImportData,
    previewExternalDataTableImport,
    completeExternalDataTableImport,
  } = useExternalDataTableImport();

  const separator = ref(CSV_SEPARATORS.comma);
  const activeImportFileId = ref(null);
  const form = ref({});
  const needPreview = ref(false);

  const isReadyForComplete = computed(() => !Object.values(form.value).some(value => value.messages.length));

  const {
    data: records,
    meta,
    pending,
    resetQuery,
    options,
    updateOptions,
    fetchList,
  } = useFetchListWithoutStoreWithOptions({
    fetchListHandler: rest => fetchExternalDataTableImportData({
      ...rest,
      id: activeImportFileId.value,
    }),
  });

  const { poll: importFile } = useFilePolling({
    createHandler: createExternalDataTableImport,
    fetchHandler: fetchExternalDataTableImportStatus,
    completedStatus: IMPORT_STATUSES.completed,
    failedStatus: IMPORT_STATUSES.failed,
  });

  const { poll: pollPreview } = usePolling({
    startHandler: () => previewExternalDataTableImport({
      id: activeImportFileId.value,
      data: {
        column_configs: formToExternalDataTableColumnConfigs(form.value),
      },
    }),
    processHandler: async (_, resolve) => {
      const response = await fetchExternalDataTableImportStatus({
        id: activeImportFileId.value,
      });

      if ([IMPORT_PREVIEW_STATUSES.completed, IMPORT_PREVIEW_STATUSES.failed].includes(response.status)) {
        return resolve(response);
      }

      return response;
    },
  });

  const toggleOnNeedPreview = () => needPreview.value = true;
  const toggleOffNeedPreview = () => needPreview.value = false;

  const { pending: updatingPreview, handler: updatePreview } = usePendingHandler(async () => {
    const newForm = mapValues(form.value, value => ({ ...value, messages: [], rows: [] }));

    try {
      const { error_info: errorInfo = {} } = await pollPreview();

      Object.entries(errorInfo).forEach(([key, value]) => newForm[key] = {
        ...newForm[key],
        ...value,
      });
    } catch (err) {
      console.error(err);
    } finally {
      await fetchList();

      form.value = newForm;

      if (isReadyForComplete.value) {
        toggleOffNeedPreview();
      }
    }
  });

  const {
    pending: uploading,
    handler: chooseFile,
  } = usePendingHandler(async ([file] = []) => {
    try {
      resetQuery();

      activeImportFileId.value = null;

      if (!file) {
        return;
      }

      validator.errors.clear();

      const data = {
        separator: separator.value,
        file,
      };

      const {
        _id: id,
        column_configs: columnConfigs = [],
      } = await importFile({ id: config.value.externalDataTable._id, data });

      activeImportFileId.value = id;

      form.value = externalDataTableColumnConfigsToForm(columnConfigs, true);

      fetchList();
    } catch (err) {
      if (!err.file) {
        throw err;
      }

      validator.errors.add({
        field: 'file',
        msg: err.file,
      });
    }
  });

  const completeImport = async (afterSubmitCallback) => {
    await completeExternalDataTableImport({
      id: activeImportFileId.value,
      data: {
        column_tags: formToExternalDataTableColumnTags(form.value),
      },
    });

    await afterSubmitCallback?.(activeImportFileId.value);
  };

  return {
    // State
    separator,
    activeImportFileId,
    form,
    needPreview,

    // Computed
    isReadyForComplete,

    // Data
    records,
    meta,
    pending,
    options,

    // Loading states
    updatingPreview,
    uploading,

    // Methods
    updateOptions,
    updatePreview,
    chooseFile,
    completeImport,
    toggleOnNeedPreview,
  };
};
