<template lang="pug">
  div
    c-advanced-data-table.white(
      :headers="headers",
      :items="templates",
      :loading="pending",
      :total-items="totalItems",
      :pagination.sync="pagination",
      advanced-pagination,
      search
    )
      template(#last_modified="{ item }") {{ item.last_modified | date }}
      template(#actions="{ item }")
        v-layout(row)
          c-action-btn(
            type="edit",
            @click.stop="edit(item)"
          )
          c-action-btn(
            type="delete",
            @click.stop="remove(item)"
          )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { localQueryMixin } from '@/mixins/query-local/query';

const { mapActions } = createNamespacedHelpers('view/widget');

export default {
  mixins: [localQueryMixin],
  data() {
    return {
      pending: false,
      templates: [],
      totalItems: 0,
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.templateName'), value: 'name', sortable: false },
        { text: this.$t('common.type'), value: 'type', sortable: false },
        { text: this.$t('common.lastModifiedOn'), value: 'last_modified', sortable: false },
        { text: this.$t('common.lastModifiedBy'), value: 'last_modified_by', sortable: false },
        { text: this.$t('common.actionsLabel'), sortable: false },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions([
      'fetchWidgetTemplatesWithoutStore',
      'createWidgetTemplate',
      'updateWidgetTemplate',
      'removeWidgetTemplate',
    ]),

    edit() {},
    remove() {},

    async fetchList() {
      this.pending = true;

      const { data, meta } = await this.fetchWidgetTemplatesWithoutStore({ params: this.getQuery() });

      this.templates = data;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
