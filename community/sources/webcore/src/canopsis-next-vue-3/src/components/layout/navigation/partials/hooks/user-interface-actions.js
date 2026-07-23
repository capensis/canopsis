import { useModals } from '@/hooks/modals';
import { useI18n } from '@/hooks/i18n';
import { MODALS } from '@/constants';

export const useUserInterfaceActions = () => {
  const modals = useModals();
  const { t } = useI18n();

  const showUserInterfaceModal = () =>
    modals.show({
      name: MODALS.info,
      config: {
        title: t('common.parameters'),
        text: t('errors.default'),
      },
    });

  return {
    showUserInterfaceModal,
  };
};
