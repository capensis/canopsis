<template>
  <div>
    <c-page-header :name="ROUTES_NAMES.adminApiAccessTokens" />

    <v-card class="ma-4 mt-0">
      <v-tabs v-model="activeTab" align-tabs="center" color="primary">
        <v-tab :value="API_ACCESS_TOKENS_TABS.apiTokens">
          {{ t('apiToken.publicApiTokens') }}
        </v-tab>
        <v-tab v-if="hasReadShareTokenAccess" :value="API_ACCESS_TOKENS_TABS.shareTokens">
          {{ t('apiToken.sharedViewTokens') }}
        </v-tab>
      </v-tabs>

      <v-card-text>
        <v-window v-model="activeTab">
          <v-window-item :value="API_ACCESS_TOKENS_TABS.apiTokens">
            <api-tokens />
          </v-window-item>
          <v-window-item
            v-if="hasReadShareTokenAccess"
            :value="API_ACCESS_TOKENS_TABS.shareTokens"
          >
            <share-tokens />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>

    <c-fab-btn :has-access="hasCreateTokenAccess" @refresh="refresh" @create="showCreateTokenModal">
      <span>{{ createTokenLabel }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { APP_HOST, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import {
  API_ACCESS_TOKENS_TABS,
  MODALS,
  ROUTES_NAMES,
  USER_PERMISSIONS,
} from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useApiToken } from '@/hooks/store/modules/api-token';
import { useShareToken } from '@/hooks/store/modules/share-token';

import ApiTokens from '@/components/other/api-token/api-tokens.vue';
import ShareTokens from '@/components/other/share-token/share-tokens.vue';

export default {
  components: {
    ApiTokens,
    ShareTokens,
  },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const router = useRouter();
    const { fetchApiTokensListWithPreviousParams, createApiToken } = useApiToken();
    const {
      fetchShareTokensListWithPreviousParams,
      createShareToken,
    } = useShareToken();
    const {
      hasCreateAccess: hasCreateShareTokenAccess,
      hasReadAccess: hasReadShareTokenAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.shareToken);
    const activeTab = ref(API_ACCESS_TOKENS_TABS.apiTokens);
    const isApiTokensTab = computed(() => activeTab.value === API_ACCESS_TOKENS_TABS.apiTokens);
    const hasCreateTokenAccess = computed(
      () => isApiTokensTab.value || hasCreateShareTokenAccess.value,
    );
    const createTokenLabel = computed(() => t(
      isApiTokensTab.value ? 'apiToken.create' : 'apiToken.createSharedViewToken',
    ));

    const refresh = () => {
      if (isApiTokensTab.value) {
        fetchApiTokensListWithPreviousParams();
      } else {
        fetchShareTokensListWithPreviousParams();
      }
    };

    const showCreateApiTokenModal = () => {
      modals.show({
        name: MODALS.createApiToken,
        config: {
          action: async (data) => {
            const apiToken = await createApiToken({ data });

            modals.show({
              name: MODALS.copyValue,
              config: {
                title: t('apiToken.createdTitle'),
                label: t('apiToken.publicApiToken'),
                value: apiToken.value,
                warningText: t('apiToken.copyWarning'),
                copySuccessText: t('apiToken.copied'),
              },
            });

            fetchApiTokensListWithPreviousParams().catch(err => console.warn(err));
          },
        },
      });
    };

    const showCreateShareTokenModal = ({ viewId, tabId, view }) => {
      modals.show({
        name: MODALS.createShareToken,
        config: {
          title: t('apiToken.createSharedViewToken'),
          action: async (data) => {
            const shareToken = await createShareToken({ data });

            const { href } = router.resolve({
              name: ROUTES_NAMES.viewKiosk,
              params: { id: viewId, tabId },
              query: { [ROUTER_ACCESS_TOKEN_KEY]: shareToken.value },
            });

            modals.show({
              name: MODALS.copyValue,
              config: {
                title: t('view.shareView', { name: view?.title ?? '' }),
                label: t('view.sharedViewUrl'),
                value: removeTrailingSlashes(`${APP_HOST}${href}`),
                copyTooltip: t('common.copyLink'),
                copySuccessText: t('success.linkCopied'),
              },
            });

            fetchShareTokensListWithPreviousParams().catch(err => console.warn(err));
          },
        },
      });
    };

    const showSelectSharedViewTabModal = () => {
      modals.show({
        name: MODALS.selectViewTab,
        config: {
          title: t('apiToken.selectSharedViewTab'),
          publicOnly: true,
          action: showCreateShareTokenModal,
        },
      });
    };

    const showCreateTokenModal = () => {
      if (isApiTokensTab.value) {
        showCreateApiTokenModal();

        return;
      }

      showSelectSharedViewTabModal();
    };

    return {
      API_ACCESS_TOKENS_TABS,
      ROUTES_NAMES,
      activeTab,
      hasReadShareTokenAccess,
      isApiTokensTab,
      hasCreateTokenAccess,
      createTokenLabel,
      t,
      refresh,
      showCreateTokenModal,
    };
  },
};
</script>
