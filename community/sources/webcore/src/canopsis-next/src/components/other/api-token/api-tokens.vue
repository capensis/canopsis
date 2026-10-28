<template>
  <api-tokens-list
    :api-tokens="apiTokens"
    :total-items="apiTokensMeta.total_count"
    :options="options"
    :pending="apiTokensPending"
    @update:options="updateOptions"
    @remove="showRemoveApiTokenModal"
    @remove-selected="showRemoveSelectedApiTokensModal"
  />
</template>

<script>
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useApiToken } from '@/hooks/store/modules/api-token';

import ApiTokensList from './api-tokens-list.vue';

export default {
  components: { ApiTokensList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const { apiTokens, apiTokensMeta, apiTokensPending, fetchApiTokensList, removeApiToken } = useApiToken();
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
      },
      onUpdate: fetchQuery => fetchApiTokensList({
        params: {
          ...getQueryForList(fetchQuery),
          active: true,
        },
      }),
    });

    const removeTokens = async (tokens) => {
      await Promise.all(tokens.map(({ _id }) => removeApiToken({ id: _id })));

      popups.success({ text: t('success.default') });

      await fetchList();
    };

    const showRemoveApiTokenModal = (apiToken) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('apiToken.deleteConfirmation', { name: apiToken.name }),
          action: () => removeTokens([apiToken]),
        },
      });
    };

    const showRemoveSelectedApiTokensModal = (selected) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          text: t('apiToken.deleteSelectedConfirmation', { count: selected.length }),
          action: () => removeTokens(selected),
        },
      });
    };

    onMounted(fetchList);

    return {
      apiTokens,
      apiTokensMeta,
      apiTokensPending,
      options,
      updateOptions,
      showRemoveApiTokenModal,
      showRemoveSelectedApiTokensModal,
    };
  },
};
</script>
