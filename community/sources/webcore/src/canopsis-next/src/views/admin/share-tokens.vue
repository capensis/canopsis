<template lang="pug">
  div
    c-page-header
    v-card-text
      share-tokens-list(
        :share-tokens="shareTokens",
        :total-items="shareTokensMeta.total_count",
        :pagination.sync="pagination",
        :pending="shareTokensPending",
        :removable="hasDeleteAnyShareTokenAccess",
        @remove="showRemoveShareTokenModal",
        @remove-selected="showRemoveSelectedShareTokensModal"
      )
    c-fab-btn(:has-access="false", @refresh="fetchList")
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesShareTokenMixin } from '@/mixins/entities/share-token';
import { permissionsTechnicalShareTokenMixin } from '@/mixins/permissions/technical/share-token';
import { localQueryMixin } from '@/mixins/query-local/query';
import { authMixin } from '@/mixins/auth';

import ShareTokensList from '@/components/other/share-token/exploitation/share-tokens-list.vue';

export default {
  components: {
    ShareTokensList,
  },
  mixins: [
    entitiesShareTokenMixin,
    localQueryMixin,
    permissionsTechnicalShareTokenMixin,
    authMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showRemoveShareTokenModal(shareToken) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeShareToken({ id: shareToken._id });

            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedShareTokensModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id }) => this.removeShareToken({ id: _id })));

            this.$popups.success({ text: this.$t('success.default') });

            await this.fetchList();
          },
        },
      });
    },

    fetchList() {
      return this.fetchShareTokensList({ params: this.getQuery() });
    },
  },
};
</script>
