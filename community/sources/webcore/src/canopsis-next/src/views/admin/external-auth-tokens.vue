<template>
  <c-page
    :creatable="hasCreateAnyExternalAuthTokensAccess"
    :create-tooltip="$t('modals.createExternalAuthToken.create.title')"
    @refresh="fetchList"
    @create="showCreateExternalAuthTokenModal"
  >
    <external-auth-tokens-list
      :tokens="tokens"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      :editable="hasUpdateAnyExternalAuthTokensAccess"
      :deletable="hasDeleteAnyExternalAuthTokensAccess"
      @edit="showEditExternalAuthTokenModal"
      @remove="showRemoveExternalAuthTokenModal"
      @update:options="updateOptions"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useWebhookTokenRule } from '@/hooks/store/modules/webhook-token-rule';
import { useCallActionWithPopup } from '@/hooks/actions/call';
import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import ExternalAuthTokensList from '@/components/other/external-auth-token/external-auth-tokens-list.vue';

export default {
  components: { ExternalAuthTokensList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();

    /**
     * PERMISSIONS
     */
    const {
      hasCreateAccess: hasCreateAnyExternalAuthTokensAccess,
      hasUpdateAccess: hasUpdateAnyExternalAuthTokensAccess,
      hasDeleteAccess: hasDeleteAnyExternalAuthTokensAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.externalAuthTokens);

    /**
     * STORE
     */
    const {
      createWebhookTokenRule,
      updateWebhookTokenRule,
      removeWebhookTokenRule,
      fetchWebhookTokenRulesListWithoutStore,
    } = useWebhookTokenRule();

    const { callActionWithPopup } = useCallActionWithPopup();

    const {
      data: tokens,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: params => fetchWebhookTokenRulesListWithoutStore({ params: { ...params, with_flags: true } }),
    });

    /**
     * METHODS
     */
    const showCreateExternalAuthTokenModal = () => {
      modals.show({
        name: MODALS.createExternalAuthToken,
        config: {
          action: newExternalAuthToken => callActionWithPopup(
            () => createWebhookTokenRule({ data: newExternalAuthToken }),
            fetchList,
          ),
        },
      });
    };

    const showEditExternalAuthTokenModal = (externalAuthToken) => {
      modals.show({
        name: MODALS.createExternalAuthToken,
        config: {
          externalAuthToken,
          title: t('modals.createExternalAuthToken.edit.title'),

          action: newExternalAuthToken => callActionWithPopup(
            () => updateWebhookTokenRule({ id: externalAuthToken._id, data: newExternalAuthToken }),
            fetchList,
          ),
        },
      });
    };

    const showRemoveExternalAuthTokenModal = (externalAuthToken) => {
      modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          title: t('modals.confirmationPhrase.deleteExternalAuthToken.title'),
          text: t('modals.confirmationPhrase.deleteExternalAuthToken.text'),
          phraseText: t('modals.confirmationPhrase.deleteExternalAuthToken.phraseText'),
          phrase: externalAuthToken.name,
          action: () => callActionWithPopup(
            () => removeWebhookTokenRule({ id: externalAuthToken._id }),
            fetchList,
          ),
        },
      });
    };

    onMounted(fetchList);

    return {
      hasCreateAnyExternalAuthTokensAccess,
      hasUpdateAnyExternalAuthTokensAccess,
      hasDeleteAnyExternalAuthTokensAccess,
      tokens,
      meta,
      pending,
      options,

      updateOptions,
      showCreateExternalAuthTokenModal,
      showEditExternalAuthTokenModal,
      showRemoveExternalAuthTokenModal,
      fetchList,
    };
  },
};
</script>
