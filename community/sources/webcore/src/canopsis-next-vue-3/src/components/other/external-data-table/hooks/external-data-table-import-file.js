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

import { useForm } from '@/hooks/form';
import { useFilePolling, usePolling } from '@/hooks/polling';
import { usePendingHandler } from '@/hooks/query/pending';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useExternalDataTableImport } from '@/hooks/store/modules/external-data-table-import';

/**
 * Composable for managing external data table import form state and validation.
 * Handles form data, error states, and form-related operations.
 *
 * @param {Object} formContext Form context created by `useForm`.
 * @returns {Object} Form state management utilities.
 * @returns {import('vue').Ref<Object>} returns.form Reactive import column form.
 * @returns {import('vue').Ref<Object>} returns.originalForm Last previewed column form.
 * @returns {import('vue').Ref<boolean>} returns.needPreview Whether a preview refresh is required.
 * @returns {import('vue').ComputedRef<boolean>} returns.hasErrors Whether any column has preview error messages.
 * @returns {Function} returns.toggleOnNeedPreview Enables preview mode.
 * @returns {Function} returns.toggleOffNeedPreview Disables preview mode.
 * @returns {Function} returns.setForm Sets form data.
 * @returns {Function} returns.clearErrors Clears column preview errors and form validation errors.
 */
const useExternalDataTableImportForm = (formContext) => {
  const form = ref({});
  const originalForm = ref({});
  const needPreview = ref(true);

  /**
   * Whether any column currently has preview error messages.
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const hasErrors = computed(() =>
    Object.values(form.value).some(value => value?.messages?.length),
  );

  /**
   * Enables preview mode so the next submit refreshes the import preview.
   */
  const toggleOnNeedPreview = () => (needPreview.value = true);

  /**
   * Disables preview mode so the next submit completes the import.
   */
  const toggleOffNeedPreview = () => (needPreview.value = false);

  /**
   * Sets the import column form.
   *
   * @param {Object} newForm New form data to set.
   * @param {boolean} [onlyForm=false] Whether to update only the editable form and keep the original snapshot.
   */
  const setForm = (newForm, onlyForm = false) => {
    form.value = newForm;

    if (!onlyForm) {
      originalForm.value = newForm;
    }
  };

  /**
   * Clears column preview errors and form validation errors.
   */
  const clearErrors = () => {
    setForm(
      mapValues(form.value, value => ({ ...value, messages: [], rows: [] })),
      true,
    );
    formContext.resetErrors();
  };

  return {
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

/**
 * Composable for handling external data table file upload operations.
 * Manages file selection, upload process, and initial configuration setup.
 *
 * @param {Object} params Configuration parameters.
 * @param {import('vue').ComputedRef<Object>} params.config External data table modal configuration.
 * @param {Object} params.formContext Form context created by `useForm`.
 * @param {Function} params.setForm Function to set form data.
 * @param {Function} params.toggleOnNeedPreview Function to enable preview mode.
 * @returns {Object} File upload utilities and state.
 * @returns {import('vue').Ref<string>} returns.separator CSV separator.
 * @returns {import('vue').Ref<string|null>} returns.activeImportFileId Active import job identifier.
 * @returns {import('vue').Ref<boolean>} returns.uploading Whether a file upload is in progress.
 * @returns {Function} returns.chooseFile Uploads a file or resets import state when the file is cleared.
 */
const useExternalDataTableFileUpload = ({
  config,
  formContext,
  setForm,
  toggleOnNeedPreview,
}) => {
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

  /**
   * Handles file selection for an external data table import.
   * Clears current import state, uploads the selected file when present, and loads column configs.
   *
   * @param {File[]|File} [selectedFiles=[]] Selected files from the file input.
   * @param {Function} resetQuery Function to reset the preview query state.
   * @param {Function} fetchList Function to fetch the preview records list.
   */
  const chooseFileHandler = async (selectedFiles = [], resetQuery, fetchList) => {
    try {
      resetQuery?.();

      activeImportFileId.value = null;
      setForm({});
      formContext.resetErrors();
      toggleOnNeedPreview();

      const files = Array.isArray(selectedFiles) ? selectedFiles : [selectedFiles].filter(Boolean);
      const [file] = files;

      if (!file) {
        return;
      }

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

      formContext.setFieldTouched('file', true);
      formContext.setFieldError('file', err.file ?? err.fail_reason);
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

/**
 * Composable for managing external data table import preview functionality.
 * Handles preview generation, polling, and error information processing.
 *
 * @param {Object} params Configuration parameters.
 * @param {import('vue').Ref<string|null>} params.activeImportFileId Active import job identifier.
 * @param {import('vue').Ref<Object>} params.form Import column form.
 * @param {Function} params.clearErrors Function to clear form errors.
 * @param {Function} params.fetchList Function to fetch the preview records list.
 * @param {Function} params.toggleOffNeedPreview Function to disable preview mode.
 * @param {import('vue').ComputedRef<boolean>} params.hasErrors Whether the form has preview errors.
 * @param {Function} params.setForm Function to set form data.
 * @returns {Object} Preview utilities.
 * @returns {Function} returns.updatePreview Polls preview status and refreshes preview data.
 */
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

  /**
   * Updates the preview for the imported file by polling error information and refreshing the data list.
   */
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

/**
 * Composable for handling external data table import completion operations.
 * Manages the final submission and completion of the import process.
 *
 * @param {Object} params Configuration parameters.
 * @param {import('vue').Ref<string|null>} params.activeImportFileId Active import job identifier.
 * @param {import('vue').Ref<Object>} params.form Import column form.
 * @returns {Object} Import completion utilities.
 * @returns {Function} returns.completeImport Completes the import and runs an optional callback.
 */
const useExternalDataTableImportCompletion = ({ activeImportFileId, form }) => {
  const { completeExternalDataTableImport } = useExternalDataTableImport();

  /**
   * Completes the external data table import by submitting the final column tags
   * and executing an optional callback after successful completion.
   *
   * @param {Function} [afterSubmitCallback] Optional callback executed with the import file id.
   */
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
 * @param {Object} params Configuration parameters.
 * @param {import('vue').ComputedRef<Object>|Object} params.config External data table modal configuration.
 * @returns {Object} Complete import functionality interface.
 */
export const useExternalDataTableImportFile = ({ config }) => {
  const formContext = useForm();
  const { fetchExternalDataTableImportData } = useExternalDataTableImport();

  const {
    form,
    originalForm,
    needPreview,
    hasErrors,
    toggleOnNeedPreview,
    toggleOffNeedPreview,
    clearErrors,
    setForm,
  } = useExternalDataTableImportForm(formContext);

  const {
    separator,
    activeImportFileId,
    uploading,
    chooseFile: chooseFileBase,
  } = useExternalDataTableFileUpload({
    config,
    form,
    formContext,
    setForm,
    toggleOnNeedPreview,
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

  /**
   * Selects an import file or clears preview records when the current file is removed.
   *
   * @param {File[]|File} files Selected files from the file input.
   */
  const chooseFile = (files) => {
    const selectedFiles = Array.isArray(files) ? files : [files].filter(Boolean);

    if (!selectedFiles.length) {
      records.value = [];
      meta.value = {};
    }

    return chooseFileBase(files, resetQuery, fetchList);
  };

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

  /**
   * Marks a preview refresh as required and validates the form on the next tick.
   */
  const validate = () => {
    toggleOnNeedPreview();

    nextTick(() => formContext.validate());
  };

  return {
    formContext,
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
