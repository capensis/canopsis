<template>
  <share-tokens-list
    :share-tokens="shareTokens"
    :total-items="shareTokensMeta.total_count"
    :options="options"
    :pending="shareTokensPending"
    :removable="hasDeleteAccess"
    @update:options="updateOptions"
    @remove="showRemoveShareTokenModal"
    @remove-selected="showRemoveSelectedShareTokensModal"
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
import { useShareToken } from '@/hooks/store/modules/share-token';

import ShareTokensList from './share-tokens-list.vue';

export default {
  components: {
    ShareTokensList,
  },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { shareTokens, shareTokensMeta, shareTokensPending, fetchShareTokensList, removeShareToken } =
      useShareToken();
    const { hasDeleteAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.shareToken);

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
      onUpdate: fetchQuery => fetchShareTokensList({
        params: getQueryForList(fetchQuery),
      }),
    });

    const showRemoveShareTokenModal = (shareToken) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('shareToken.deleteConfirmation', { name: shareToken.description }),
          action: async () => {
            await removeShareToken({ id: shareToken._id });

            popups.success({ text: t('success.default') });

            await fetchList();
          },
        },
      });
    };

    const showRemoveSelectedShareTokensModal = (selected) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('shareToken.deleteSelectedConfirmation', { count: selected.length }),
          action: async () => {
            await Promise.all(selected.map(({ _id }) => removeShareToken({ id: _id })));

            popups.success({ text: t('success.default') });

            await fetchList();
          },
        },
      });
    };

    onMounted(fetchList);

    return {
      shareTokens,
      shareTokensMeta,
      shareTokensPending,
      options,
      hasDeleteAccess,
      updateOptions,
      showRemoveShareTokenModal,
      showRemoveSelectedShareTokensModal,
    };
  },
};
</script>
