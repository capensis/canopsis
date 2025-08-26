import { mapValues } from 'lodash';
import { computed, nextTick, ref } from 'vue';

import { CSV_SEPARATORS, IMPORT_STATUSES, IMPORT_PREVIEW_STATUSES } from '@/constants';

import {
  externalDataTableColumnConfigsToForm,
  formToExternalDataTableColumnConfigs,
  formToExternalDataTableColumnTags,
} from '@/helpers/entities/external-data-table/form';

import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useExternalDataTableImport } from '@/hooks/store/modules/external-data-table-import';
import { useFilePolling, usePolling } from '@/hooks/polling';
import { usePendingHandler } from '@/hooks/query/pending';
import { useValidator } from '@/hooks/validator/validator';

/**
 * Composable for managing external data table import form state and validation.
 * Handles form data, error states, and form-related operations.
 *
 * @returns {Object} Form state management utilities
 * @returns {Object} returns.validator - Validator instance for form validation
 * @returns {Object} returns.form - Reactive form data reference
 * @returns {Object} returns.needPreview - Reactive flag indicating if preview is needed
 * @returns {Object} returns.hasErrors - Computed property indicating if form has validation errors
 * @returns {Function} returns.toggleOnNeedPreview - Function to enable preview mode
 * @returns {Function} returns.toggleOffNeedPreview - Function to disable preview mode
 * @returns {Function} returns.setForm - Function to set form data
 * @returns {Function} returns.clearErrors - Function to clear all form errors and validation messages
 */
const useExternalDataTableImportForm = () => {
  const validator = useValidator();

  const form = ref({});
  const needPreview = ref(false);

  /**
   * Computed property that checks if the form has any validation errors.
   *
   * @returns {boolean} True if any form field has error messages
   */
  const hasErrors = computed(() => Object.values(form.value).some(value => value?.messages?.length));

  /**
   * Enables the preview mode for the data table import.
   * Sets the needPreview flag to true, indicating that a preview is required.
   */
  const toggleOnNeedPreview = () => needPreview.value = true;

  /**
   * Disables the preview mode for the data table import.
   * Sets the needPreview flag to false, indicating that no preview is needed.
   */
  const toggleOffNeedPreview = () => needPreview.value = false;

  /**
   * Sets the form data.
   *
   * @param {Object} newForm - New form data to set
   */
  const setForm = (newForm) => {
    form.value = newForm;
  };

  /**
   * Clears all form errors and validation messages.
   * Resets both form field messages and validator errors.
   */
  const clearErrors = () => {
    setForm(mapValues(form.value, value => ({ ...value, messages: [], rows: [] })));
    validator.errors.clear();
  };

  return {
    validator,

    form,
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
 * @param {Object} params - Configuration parameters
 * @param {Object} params.config - External data table configuration
 * @param {Object} params.validator - Validator instance
 * @param {Function} params.setForm - Function to set form data
 * @returns {Object} File upload utilities and state
 */
const useExternalDataTableFileUpload = ({
  config,
  validator,
  setForm,
}) => {
  const {
    createExternalDataTableImport,
    fetchExternalDataTableImportStatus,
  } = useExternalDataTableImport();

  const separator = ref(CSV_SEPARATORS.comma);
  const activeImportFileId = ref(null);

  const { poll: importFile } = useFilePolling({
    createHandler: createExternalDataTableImport,
    fetchHandler: fetchExternalDataTableImportStatus,
    completedStatus: IMPORT_STATUSES.completed,
    failedStatus: IMPORT_STATUSES.failed,
  });

  /**
   * Handles the file selection and import process for external data table.
   * Uploads the selected file, processes column configurations, and updates the form state.
   *
   * @param {File[]} [files=[]] - Array of selected files, typically from file input
   * @param {Function} resetQuery - Function to reset query state
   * @param {Function} fetchList - Function to fetch data list
   */
  const chooseFileHandler = async ([file] = [], resetQuery, fetchList) => {
    try {
      resetQuery?.();

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

      setForm(externalDataTableColumnConfigsToForm(columnConfigs, true));

      fetchList?.();
    } catch (err) {
      if (!err.file) {
        throw err;
      }

      validator.errors.add({
        field: 'file',
        msg: err.file,
      });
    }
  };

  const {
    pending: uploading,
    handler: chooseFile,
  } = usePendingHandler(chooseFileHandler);

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
 * @param {Object} params - Configuration parameters
 * @param {Object} params.activeImportFileId - Reference to active import file ID
 * @param {Object} params.form - Form state reference
 * @param {Function} params.clearErrors - Function to clear form errors
 * @param {Function} params.fetchList - Function to fetch data list
 * @param {Function} params.toggleOffNeedPreview - Function to disable preview mode
 * @param {Object} params.hasErrors - Computed property indicating if form has errors
 * @param {Function} params.setForm - Function to set form data
 * @returns {Object} Preview utilities and state
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
  const {
    fetchExternalDataTableImportStatus,
    previewExternalDataTableImport,
  } = useExternalDataTableImport();

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

  /**
   * Updates the preview data for the imported file by polling for error information
   * and refreshing the data list. Handles form state updates and preview mode management.
   */
  const updatePreview = async () => {
    clearErrors();

    const newForm = mapValues(form.value, value => ({ ...value }));

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
 * @param {Object} params - Configuration parameters
 * @param {Object} params.activeImportFileId - Reference to active import file ID
 * @param {Object} params.form - Form state reference
 * @returns {Object} Import completion utilities
 */
const useExternalDataTableImportCompletion = ({ activeImportFileId, form }) => {
  const { completeExternalDataTableImport } = useExternalDataTableImport();

  /**
   * Completes the external data table import process by submitting the final configuration
   * and executing an optional callback after successful completion.
   *
   * @param {Function} [afterSubmitCallback] - Optional callback function to execute after import completion
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
 * @param {Object} params - Configuration parameters
 * @param {Object} params.config - External data table configuration
 * @returns {Object} Complete import functionality interface
 */
export const useExternalDataTableImportFile = ({ config }) => {
  const { fetchExternalDataTableImportData } = useExternalDataTableImport();

  // Form state management
  const {
    validator,
    form,
    needPreview,
    hasErrors,
    toggleOnNeedPreview,
    toggleOffNeedPreview,
    clearErrors,
    setForm,
  } = useExternalDataTableImportForm();

  // File upload functionality
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

  // Data fetching with pagination
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

  /**
   * Validates the form and triggers preview mode.
   * Sets preview requirement and validates all form fields in the next tick.
   */
  const validate = () => {
    toggleOnNeedPreview();

    nextTick(() => validator.validateAll());
  };

  return {
    validator,

    // State
    separator,
    activeImportFileId,
    form,
    needPreview,

    // Computed
    hasErrors,

    // Data
    records,
    meta,
    pending,
    options,

    // Loading states
    uploading,

    // Methods
    updateOptions,
    updatePreview,
    chooseFile,
    completeImport,
    validate,
  };
};
