<template lang="pug">
  v-tooltip(left)
    v-btn(
      slot="activator",
      fab,
      dark,
      small,
      @click.stop="showCreateShareTokenModal"
    )
      v-icon link
    span {{ $t('common.shareLink') }}
</template>

<script>
import { APP_HOST, ROUTER_ACCESS_TOKEN_KEY } from '@/config';

import { MODALS, ROUTES_NAMES } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';
import { writeTextToClipboard } from '@/helpers/clipboard';

import { entitiesShareTokenMixin } from '@/mixins/entities/share-token';

export default {
  mixins: [entitiesShareTokenMixin],
  methods: {
    showCreateShareTokenModal() {
      this.$modals.show({
        name: MODALS.createShareToken,
        config: {
          action: async (data) => {
            const shareToken = await this.createShareToken({ data });
            const { query, params } = this.$route;

            const { href } = this.$router.resolve(
              {
                name: ROUTES_NAMES.viewKiosk,
                params: {
                  id: params.id,
                  tabId: query.tabId,
                },
                query: { [ROUTER_ACCESS_TOKEN_KEY]: shareToken.value },
              },
              this.$route,
            );

            const url = removeTrailingSlashes(`${APP_HOST}${href}`);

            try {
              await writeTextToClipboard(url);

              this.$popups.success({ text: this.$t('success.pathCopied') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },
  },
};
</script>
