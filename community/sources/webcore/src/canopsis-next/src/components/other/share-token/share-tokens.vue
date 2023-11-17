<template>
  <share-tokens-list
    :share-tokens="shareTokens"
    :total-items="shareTokensMeta.total_count"
    :options.sync="options"
    :pending="shareTokensPending"
    :removable="hasDeleteAnyShareTokenAccess"
    @remove="showRemoveShareTokenModal"
    @remove-selected="showRemoveSelectedShareTokensModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesShareTokenMixin } from '@/mixins/entities/share-token';
import { permissionsTechnicalShareTokenMixin } from '@/mixins/permissions/technical/share-token';
import { localQueryMixin } from '@/mixins/query-local/query';

import ShareTokensList from '@/components/other/share-token/share-tokens-list.vue';

export default {
  components: {
    ShareTokensList,
  },
  mixins: [
    entitiesShareTokenMixin,
    localQueryMixin,
    permissionsTechnicalShareTokenMixin,
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
