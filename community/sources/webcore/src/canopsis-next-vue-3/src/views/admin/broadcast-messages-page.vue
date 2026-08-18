<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="t('modals.createBroadcastMessage.create.title')"
    @refresh="refresh"
    @create="create"
  >
    <broadcast-messages />
  </c-page>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';

import BroadcastMessages from '@/components/other/broadcast-message/broadcast-messages.vue';

export default {
  components: {
    BroadcastMessages,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const { createBroadcastMessage, fetchBroadcastMessagesListWithPreviousParams } =
      useBroadcastMessages();
    const { hasCreateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.broadcastMessage);

    const refresh = () => fetchBroadcastMessagesListWithPreviousParams();

    const create = () =>
      modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          action: async (data) => {
            await createBroadcastMessage({ data });
            popups.success({ text: t('success.default') });

            await refresh();
          },
        },
      });

    return {
      t,
      hasCreateAccess,
      refresh,
      create,
    };
  },
};
</script>
