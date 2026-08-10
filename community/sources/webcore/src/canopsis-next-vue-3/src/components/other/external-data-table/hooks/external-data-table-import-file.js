import { mapValues } from 'lodash';
import { computed, nextTick, ref } from 'vue';

import {
  CSV_SEPARATORS,
  IMPORT_STATUSES,
  EXTERNAL_DATA_TABLE_IMPORT_PREVIEW_STATUSES,
} from '@/constants';
import {
  externalDataTableColumnConfigsToForm,
  formToExternalDataTableColumnConfigs,
  formToExternalDataTableColumnTags,
  addPriorityColumnForRegexpTypes,
} from '@/helpers/entities/external-data-table/form';
import { useFilePolling, usePolling } from '@/hooks/polling';
import { usePendingHandler } from '@/hooks/query/pending';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useExternalDataTableImport } from '@/hooks/store/modules/external-data-table-import';
import { useValidator } from '@/hooks/validator/validator';

const useExternalDataTableImportForm = () => {
  const validator = useValidator();

  const form = ref({});
  const originalForm = ref({});
  const needPreview = ref(true);

  const hasErrors = computed(() =>
    Object.values(form.value).some(value => value?.messages?.length),
  );

  const toggleOnNeedPreview = () => (needPreview.value = true);

  const toggleOffNeedPreview = () => (needPreview.value = false);

  const setForm = (newForm, onlyForm = false) => {
    form.value = newForm;

    if (!onlyForm) {
      originalForm.value = newForm;
    }
  };

  const clearErrors = () => {
    setForm(
      mapValues(form.value, value => ({ ...value, messages: [], rows: [] })),
      true,
    );
    validator.errors.clear();
  };

  return {
    validator,

    form,
    originalForm,
    needPreview,

    hasErrors,

    toggleOnNeedPreview,
    toggleOffNeedPreview,
    setForm,
    clearErrors,
  };
};

const useExternalDataTableFileUpload = ({ config, validator, setForm }) => {
  const { createExternalDataTableImport, fetchExternalDataTableImportStatus } =
    useExternalDataTableImport();

  const separator = ref(CSV_SEPARATORS.comma);
  const activeImportFileId = ref(null);

  const { poll: importFile } = useFilePolling({
    createHandler: createExternalDataTableImport,
    fetchHandler: fetchExternalDataTableImportStatus,
    completedStatus: IMPORT_STATUSES.completed,
    failedStatus: IMPORT_STATUSES.failed,
  });

  const chooseFileHandler = async (selectedFiles = [], resetQuery, fetchList) => {
    try {
      resetQuery?.();

      activeImportFileId.value = null;

      const files = Array.isArray(selectedFiles) ? selectedFiles : [selectedFiles].filter(Boolean);
      const [file] = files;

      if (!file) {
        return;
      }

      validator.errors.clear();

      const data = {
        separator: separator.value,
        file,
      };

      const { _id: id, column_configs: columnConfigs = [] } = await importFile({
        id: config.value.externalDataTable._id,
        data,
      });

      activeImportFileId.value = id;

      setForm(externalDataTableColumnConfigsToForm(columnConfigs, true));

      fetchList?.();
    } catch (err) {
      if (!err.file && !err.fail_reason) {
        throw err;
      }

      validator.errors.add({
        field: 'file',
        msg: err.file ?? err.fail_reason,
      });
    }
  };

  const { pending: uploading, handler: chooseFile } = usePendingHandler(chooseFileHandler);

  return {
    separator,
    activeImportFileId,
    uploading,
    chooseFile,
  };
};

const useExternalDataTableImportPreview = ({
  activeImportFileId,
  form,
  clearErrors,
  fetchList,
  toggleOffNeedPreview,
  hasErrors,
  setForm,
}) => {
  const { fetchExternalDataTableImportStatus, previewExternalDataTableImport } =
    useExternalDataTableImport();

  const { poll: pollPreview } = usePolling({
    startHandler: () =>
      previewExternalDataTableImport({
        id: activeImportFileId.value,
        data: {
          column_configs: formToExternalDataTableColumnConfigs(form.value),
        },
      }),
    processHandler: async (_, resolve) => {
      const response = await fetchExternalDataTableImportStatus({
        id: activeImportFileId.value,
      });

      if (
        [
          EXTERNAL_DATA_TABLE_IMPORT_PREVIEW_STATUSES.completed,
          EXTERNAL_DATA_TABLE_IMPORT_PREVIEW_STATUSES.failed,
        ].includes(response.status)
      ) {
        return resolve(response);
      }

      return response;
    },
  });

  const updatePreview = async () => {
    clearErrors();

    let newForm = mapValues(form.value, value => ({ ...value }));

    try {
      const { error_info: errorInfo = {} } = await pollPreview();

      Object.entries(errorInfo).forEach(
        ([key, value]) =>
          (newForm[key] = {
            ...newForm[key],
            ...value,
          }),
      );
    } catch (err) {
      console.error(err);
    } finally {
      await fetchList();

      newForm = addPriorityColumnForRegexpTypes(newForm);

      setForm(newForm);

      if (!hasErrors.value) {
        toggleOffNeedPreview();
      }
    }
  };

  return {
    updatePreview,
  };
};

const useExternalDataTableImportCompletion = ({ activeImportFileId, form }) => {
  const { completeExternalDataTableImport } = useExternalDataTableImport();

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
    completeImport,
  };
};

/**
 * Main composable for external data table import functionality.
 * Orchestrates file upload, form management, preview, and completion operations.
 *
 * @param {Object} params - Configuration parameters
 * @param {Object} params.config - External data table configuration
 * @returns {Object} Complete import functionality interface
 */
export const useExternalDataTableImportFile = ({ config }) => {
  const { fetchExternalDataTableImportData } = useExternalDataTableImport();

  const {
    validator,
    form,
    originalForm,
    needPreview,
    hasErrors,
    toggleOnNeedPreview,
    toggleOffNeedPreview,
    clearErrors,
    setForm,
  } = useExternalDataTableImportForm();

  const {
    separator,
    activeImportFileId,
    uploading,
    chooseFile: chooseFileBase,
  } = useExternalDataTableFileUpload({
    config,
    form,
    validator,
    setForm,
  });

  const {
    data: records,
    meta,
    pending,
    resetQuery,
    options,
    updateOptions,
    fetchList,
  } = useFetchListWithoutStoreWithOptions({
    fetchListHandler: rest =>
      fetchExternalDataTableImportData({
        ...rest,
        id: activeImportFileId.value,
      }),
  });

  const chooseFile = files => chooseFileBase(files, resetQuery, fetchList);

  const { updatePreview } = useExternalDataTableImportPreview({
    activeImportFileId,
    form,
    clearErrors,
    fetchList,
    toggleOffNeedPreview,
    hasErrors,
    setForm,
  });

  const { completeImport } = useExternalDataTableImportCompletion({
    activeImportFileId,
    form,
  });

  const validate = () => {
    toggleOnNeedPreview();

    nextTick(() => validator.validateAll());
  };

  return {
    validator,

    separator,
    activeImportFileId,
    form,
    originalForm,
    needPreview,

    hasErrors,

    records,
    meta,
    pending,
    options,

    uploading,

    updateOptions,
    updatePreview,
    chooseFile,
    completeImport,
    validate,
  };
};
