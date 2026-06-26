<template>
  <broadcast-messages-list
    :broadcast-messages="broadcastMessages"
    :pending="pending"
    :options="options"
    :total-items="meta.total_count"
    :editable="hasUpdateAccess"
    :deletable="hasDeleteAccess"
    @edit="showEditBroadcastMessageModal"
    @remove="showRemoveBroadcastMessageModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';

import BroadcastMessagesList from './broadcast-messages-list.vue';

export default {
  components: {
    BroadcastMessagesList,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const {
      broadcastMessages,
      broadcastMessagesMeta,
      broadcastMessagesPending,
      fetchBroadcastMessagesList,
      updateBroadcastMessage,
      removeBroadcastMessage,
    } = useBroadcastMessages();
    const {
      hasUpdateAccess,
      hasDeleteAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.broadcastMessage);

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        sortBy: [],
        sortDesc: [],
        search: '',
      },
      onUpdate: fetchQuery => fetchBroadcastMessagesList({
        params: getQueryForList(fetchQuery),
      }),
    });

    const showEditBroadcastMessageModal = (message) => modals.show({
      name: MODALS.createBroadcastMessage,
      config: {
        message,
        title: t('modals.createBroadcastMessage.edit.title'),
        action: async (data) => {
          await updateBroadcastMessage({ id: message._id, data });
          popups.success({ text: t('success.default') });

          await fetchList();
        },
      },
    });

    const showRemoveBroadcastMessageModal = (id) => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          await removeBroadcastMessage({ id });
          popups.success({ text: t('success.default') });

          await fetchList();
        },
      },
    });

    onMounted(fetchList);

    return {
      broadcastMessages,
      meta: broadcastMessagesMeta,
      pending: broadcastMessagesPending,
      options,
      hasUpdateAccess,
      hasDeleteAccess,
      updateOptions,
      showEditBroadcastMessageModal,
      showRemoveBroadcastMessageModal,
    };
  },
};
</script>
