<template>
  <v-tooltip location="left">
    <template #activator="{ props }">
      <v-btn
        v-bind="props"
        color="secondary-lighten-2"
        class="text-white"
        size="default"
        icon
        @click.stop="showCreateShareTokenModal"
      >
        <v-icon icon="share" size="small" />
      </v-btn>
    </template>

    <span>{{ t('common.shareLink') }}</span>
  </v-tooltip>
</template>

<script>
import { useRouter } from 'vue-router';

import { APP_HOST, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import { MODALS, ROUTES_NAMES } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useShareToken } from '@/hooks/store/modules/share-token';

export default {
  props: {
    view: {
      type: Object,
      required: true,
    },
    tab: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const router = useRouter();
    const { createShareToken } = useShareToken();

    const showCreateShareTokenModal = () => {
      modals.show({
        name: MODALS.createShareToken,
        config: {
          action: async (data) => {
            const shareToken = await createShareToken({ data });

            const { href } = router.resolve({
              name: ROUTES_NAMES.viewKiosk,
              params: {
                id: props.view._id,
                tabId: props.tab._id,
              },
              query: { [ROUTER_ACCESS_TOKEN_KEY]: shareToken.value },
            });

            const url = removeTrailingSlashes(`${APP_HOST}${href}`);

            modals.show({
              name: MODALS.shareView,
              config: {
                title: t('view.shareView', { name: props.view.title }),
                url,
              },
            });
          },
        },
      });
    };

    return {
      t,
      showCreateShareTokenModal,
    };
  },
};
</script>
