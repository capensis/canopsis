<template>
  <v-tooltip left>
    <template #activator="{ on }">
      <v-btn
        color="secondary lighten-2"
        fab
        dark
        small
        v-on="on"
        @click.stop="showCreateShareTokenModal"
      >
        <v-icon small>
          share
        </v-icon>
      </v-btn>
    </template>
    <span>{{ $t('common.shareLink') }}</span>
  </v-tooltip>
</template>

<script>
import { APP_HOST, ROUTER_ACCESS_TOKEN_KEY } from '@/config';
import { MODALS, ROUTES_NAMES } from '@/constants';

import { removeTrailingSlashes } from '@/helpers/url';

import { entitiesShareTokenMixin } from '@/mixins/entities/share-token';

export default {
  mixins: [entitiesShareTokenMixin],
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
  methods: {
    showCreateShareTokenModal() {
      this.$modals.show({
        name: MODALS.createShareToken,
        config: {
          action: async (data) => {
            const shareToken = await this.createShareToken({ data });

            const { href } = this.$router.resolve(
              {
                name: ROUTES_NAMES.viewKiosk,
                params: {
                  id: this.view._id,
                  tabId: this.tab._id,
                },
                query: { [ROUTER_ACCESS_TOKEN_KEY]: shareToken.value },
              },
              this.$route,
            );

            const url = removeTrailingSlashes(`${APP_HOST}${href}`);

            this.$modals.show({
              name: MODALS.shareView,
              config: {
                title: this.$t('view.shareView', { name: this.view.title }),
                url,
              },
            });
          },
        },
      });
    },
  },
};
</script>
