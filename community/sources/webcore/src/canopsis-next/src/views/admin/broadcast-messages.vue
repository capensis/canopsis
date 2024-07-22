<template>
  <c-page
    :creatable="hasCreateAnyBroadcastMessageAccess"
    :create-tooltip="$t('modals.createBroadcastMessage.create.title')"
    @refresh="fetchList"
    @create="showCreateBroadcastMessageModal"
  >
    <broadcast-messages-list
      :broadcast-messages="broadcastMessages"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      :editable="hasUpdateAnyBroadcastMessageAccess"
      :deletable="hasDeleteAnyBroadcastMessageAccess"
      @edit="showEditBroadcastMessageModal"
      @remove="showRemoveBroadcastMessageModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USERS_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useBroadcastMessages } from '@/hooks/store/modules/broadcast-message';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useQueryOptions } from '@/hooks/query/options';
import { useCRUDPermissions } from '@/hooks/auth';

import BroadcastMessagesList from '@/components/other/broadcast-message/broadcast-messages-list.vue';

export default {
  components: { BroadcastMessagesList },
  setup() {
    const broadcastMessages = ref([]);
    const meta = ref({});

    const { t } = useI18n();
    const modals = useModals();

    /**
     * PERMISSIONS
     */
    const {
      hasCreateAccess: hasCreateAnyBroadcastMessageAccess,
      hasUpdateAccess: hasUpdateAnyBroadcastMessageAccess,
      hasDeleteAccess: hasDeleteAnyBroadcastMessageAccess,
    } = useCRUDPermissions(USERS_PERMISSIONS.technical.broadcastMessage);

    /**
     * STORE
     */
    const {
      createBroadcastMessage,
      updateBroadcastMessage,
      removeBroadcastMessage,
      fetchBroadcastMessagesListWithoutStore,
    } = useBroadcastMessages();
    const { callActionWithPopup } = useCallActionWithPopup();

    /**
     * QUERY
     */
    const {
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchBroadcastMessagesListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
          },
        });

        broadcastMessages.value = response.data;
        meta.value = response.meta;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * METHODS
     */
    const showCreateBroadcastMessageModal = () => {
      modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          action: newMessage => callActionWithPopup(
            () => createBroadcastMessage({ data: newMessage }),
            fetchList,
          ),
        },
      });
    };

    const showEditBroadcastMessageModal = (message) => {
      modals.show({
        name: MODALS.createBroadcastMessage,
        config: {
          message,
          title: t('modals.createBroadcastMessage.edit.title'),

          action: newMessage => callActionWithPopup(
            () => updateBroadcastMessage({ id: message._id, data: newMessage }),
            fetchList,
          ),
        },
      });
    };

    const showRemoveBroadcastMessageModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => callActionWithPopup(
            () => removeBroadcastMessage({ id }),
            fetchList,
          ),
        },
      });
    };

    onMounted(() => fetchList());

    return {
      hasCreateAnyBroadcastMessageAccess,
      hasUpdateAnyBroadcastMessageAccess,
      hasDeleteAnyBroadcastMessageAccess,
      broadcastMessages,
      meta,
      pending,
      options,

      updateQuery,
      showCreateBroadcastMessageModal,
      showEditBroadcastMessageModal,
      showRemoveBroadcastMessageModal,
      fetchList,
    };
  },
};
</script>
