<template>
  <c-page
    :creatable="hasCreateAnyIconAccess"
    :create-tooltip="$t('modals.createIcon.create.title')"
    @create="showCreateIconModal"
    @refresh="fetchIconsListWithPreviousParams"
  >
    <template #header>
      {{ $tc('common.icon', 2) }}
    </template>
    <icons />
  </c-page>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useIcon } from '@/hooks/store/modules/icon';
import { useCallActionWithPopup } from '@/hooks/actions/call';

import Icons from '@/components/other/icons/icons.vue';

export default {
  components: { Icons },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const { fetchIconsListWithPreviousParams, createIcon } = useIcon();
    const { hasCreateAccess: hasCreateAnyIconAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.icon);
    const { callActionWithPopup } = useCallActionWithPopup();

    const showCreateIconModal = () => modals.show({
      name: MODALS.createIcon,
      config: {
        title: t('modals.createIcon.create.title'),
        action: newIcon => callActionWithPopup(
          () => createIcon({ data: newIcon }),
          fetchIconsListWithPreviousParams,
          t('modals.createIcon.create.success'),
        ),
      },
    });

    return {
      hasCreateAnyIconAccess,
      showCreateIconModal,
      fetchIconsListWithPreviousParams,
    };
  },
};
</script>
