<template lang="pug">
  div
    c-advanced-data-table(
      :items="entities",
      :headers="headers",
      :loading="pending",
      :pagination.sync="pagination",
      :total-items="entitiesTotalCount",
      :no-data-text="$t('tables.noData')",
      :is-disabled-item="isSelectedEntity",
      advanced-pagination,
      select-all
    )
      template(slot="toolbar", slot-scope="props")
        v-layout(row)
          c-search-field(@submit="props.updateSearch", @clear="props.clearSearch")
      template(slot="name", slot-scope="props")
        span.text-xs-left {{ props.item.name }}
      template(slot="id", slot-scope="props")
        span.text-xs-left {{ props.item._id }}
      template(slot="actions", slot-scope="props")
        v-btn(:disabled="props.disabled", icon, small, @click="$emit('select', [props.item])")
          v-icon(color="primary") add
      template(slot="mass-actions", slot-scope="props")
        v-btn(
          color="primary",
          @click="$emit('select', props.selected)"
        ) {{ $t('contextGeneralTable.addSelection') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { BASIC_ENTITY_TYPES } from '@/constants';

import { localQueryMixin } from '@/mixins/query-local/query';

const { mapActions } = createNamespacedHelpers('entity');

export default {
  mixins: [localQueryMixin],
  props: {
    entitiesIds: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      entitiesTotalCount: 0,
      entities: [],
      query: {
        rowsPerPage: 5,
      },
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.name'), value: 'name', sortable: false },
        { text: this.$t('common.id'), value: '_id', sortable: false },
        { text: this.$t('common.actionsLabel'), value: 'actions', sortable: false },
      ];
    },
  },
  methods: {
    ...mapActions({
      fetchContextEntitiesWithoutStore: 'fetchListWithoutStore',
    }),

    isSelectedEntity({ _id }) {
      return this.entitiesIds.includes(_id);
    },

    search() {
      this.fetchList();
    },

    async fetchList() {
      this.pending = true;

      const { data = [], meta } = await this.fetchContextEntitiesWithoutStore({
        params: {
          ...this.getQuery(),
          type: Object.values(BASIC_ENTITY_TYPES),
        },
      });

      this.entities = data;
      this.entitiesTotalCount = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
