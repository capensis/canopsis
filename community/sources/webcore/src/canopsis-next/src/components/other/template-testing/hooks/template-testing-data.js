import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useTemplateData } from '@/hooks/store/modules/template-data';

/**
 * Hook for managing template testing data modals
 *
 * @param {Function} [refresh] - Optional callback function to refresh data after modal actions
 * @returns {Object} Object containing modal functions
 * @returns {Function} returns.showCreateTemplateTestingDataModal - Function to show create template data modal
 * @returns {Function} returns.showEditTemplateTestingDataModal - Function to show edit template data modal
 * @returns {Function} returns.showRemoveTemplateTestingDataModal - Function to show remove template data modal
 */
export const useTemplateDataModals = (refresh) => {
  const { t } = useI18n();
  const modals = useModals();
  const {
    createTemplateData,
    updateTemplateData,
    removeTemplateData,
  } = useTemplateData();

  /**
   * Shows the modal for creating new template testing data
   */
  const showCreateTemplateTestingDataModal = () => modals.show({
    name: MODALS.createTemplateTestingData,
    config: {
      title: t('modals.createTemplateData.title'),
      action: async (newTemplateTestingData) => {
        await createTemplateData({ data: newTemplateTestingData });

        return refresh?.();
      },
    },
  });

  /**
   * Shows the modal for editing existing template testing data
   *
   * @param {Object} [templateTestingData={}] - The template testing data object to edit
   * @param {string} templateTestingData._id - The ID of the template data
   * @param {string} templateTestingData.name - The name of the template data
   */
  const showEditTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
    name: MODALS.createTemplateTestingData,
    config: {
      templateTestingData,
      title: t('modals.createTemplateTestingData.edit.title'),
      action: async (newTemplateTestingData) => {
        await updateTemplateData({ id: templateTestingData._id, data: newTemplateTestingData });

        return refresh?.();
      },
    },
  });

  /**
   * Shows the confirmation modal for removing template testing data
   *
   * @param {Object} [templateTestingData={}] - The template testing data object to remove
   * @param {string} templateTestingData._id - The ID of the template data
   * @param {string} templateTestingData.name - The name of the template data (used as confirmation phrase)
   */
  const showRemoveTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
    name: MODALS.confirmationPhrase,
    config: {
      phrase: templateTestingData.name,
      title: t('modals.confirmationPhrase.templateTestingData.title'),
      text: t('modals.confirmationPhrase.templateTestingData.text'),
      phraseText: t('modals.confirmationPhrase.templateTestingData.phraseText'),
      action: async () => {
        await removeTemplateData({ id: templateTestingData._id });

        return refresh?.();
      },
    },
  });

  return {
    showCreateTemplateTestingDataModal,
    showEditTemplateTestingDataModal,
    showRemoveTemplateTestingDataModal,
  };
};
