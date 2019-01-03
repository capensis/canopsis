<template lang="pug">
  v-data-table.elevation-1(
  :headers="headers",
  :items="items",
  :pagination.sync="pagination",
  :total-items="meta.total",
  :rows-per-page-items="rowsPerPageItems",
  :loading="pending"
  )
    template(slot="items", slot-scope="{ item }")
      td {{ item._id }}
      td {{ item.name }}
      td {{ item.type }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import filterEditorResultsMixin from '@/mixins/filter/editor/results';

const { mapActions: entityMapActions } = createNamespacedHelpers('entity');

export default {
  mixins: [filterEditorResultsMixin],
  computed: {
    headers() {
      return [
        {
          text: this.$t('filterEditor.resultsTableHeaders.entity.id'),
          align: 'left',
          sortable: false,
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.entity.name'),
          align: 'left',
          sortable: false,
        },
        {
          text: this.$t('filterEditor.resultsTableHeaders.entity.type'),
          align: 'left',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    ...entityMapActions({
      fetchEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      const { rowsPerPage: limit, page } = this.pagination;

      this.pending = true;
      const { entities, total } = await this.fetchEntitiesListWithoutStore({
        params: {
          limit,
          start: limit * (page - 1),
          _filter: this.filter,
        },
      });

      this.pending = false;
      this.items = entities;
      this.meta = { total };
    },
  },
};
</script>
