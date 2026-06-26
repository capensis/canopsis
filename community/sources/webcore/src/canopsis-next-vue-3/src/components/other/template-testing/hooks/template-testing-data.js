import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

/**
 * Hook for managing template testing data modals
 *
 * @param {Object} [options={}]
 * @param {number} [options.type]
 * @param {Function} [options.refresh]
 * @returns {Object}
 */
export const useTemplateDataModals = ({ type, refresh } = {}) => {
  const { t } = useI18n();
  const modals = useModals();
  const {
    createTemplateTestingData,
    updateTemplateTestingData,
    removeTemplateTestingData,
  } = useTemplateTesting();

  const showCreateTemplateTestingDataModal = () => modals.show({
    name: MODALS.createTemplateTestingData,
    config: {
      templateTestingData: { type },
      title: t('modals.createTemplateData.title'),
      action: async (newTemplateTestingData) => {
        await createTemplateTestingData({ data: newTemplateTestingData });

        return refresh?.();
      },
    },
  });

  const showEditTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
    name: MODALS.createTemplateTestingData,
    config: {
      templateTestingData,
      title: t('modals.createTemplateTestingData.edit.title'),
      action: async (newTemplateTestingData) => {
        await updateTemplateTestingData({
          id: templateTestingData._id,
          data: newTemplateTestingData,
        });

        return refresh?.();
      },
    },
  });

  const showRemoveTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
    name: MODALS.confirmationPhrase,
    config: {
      phrase: templateTestingData.name,
      title: t('modals.confirmationPhrase.templateTestingData.title'),
      text: t('modals.confirmationPhrase.templateTestingData.text'),
      phraseText: t('modals.confirmationPhrase.templateTestingData.phraseText'),
      action: async () => {
        await removeTemplateTestingData({ id: templateTestingData._id });

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
