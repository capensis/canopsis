<template lang="pug">
  div
    v-toolbar.toolbar.white(dense, flat)
      search-field(
        v-model="searchingText",
        :label="$t('common.search')",
        @submit="search"
      )
    v-btn.green.white--text(
      v-show="selectedEntities.length",
      @click="$emit('update:selectedIds', selectedEntities)"
    ) {{ $t('contextGeneralTable.addSelection') }}
    v-data-table(
      v-model="selectedEntities",
      :no-data-text="$t('tables.noData')",
      :headers="headers",
      :items="entities",
      :loading="pending",
      item-key="_id",
      select-all
    )
      template(slot="items", slot-scope="props")
        td
          v-checkbox(
            v-model="props.selected",
            :value="props._id",
            primary,
            hide-details
          )
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item._id}}
        td
          v-btn(icon, @click="$emit('update:selectedIds', [props.item])")
            v-icon add
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { getContextWidgetSearchByText } from '@/helpers/entities-search';

import SearchField from '@/components/forms/fields/search-field.vue';

const { mapActions } = createNamespacedHelpers('entity');

export default {
  components: { SearchField },
  data() {
    return {
      pending: false,
      searchingText: '',
      entities: [],
      selectedEntities: [],
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
          sortable: false,
        },
      ];
    },
  },
  methods: {
    ...mapActions({
      fetchContextEntitiesWithoutStore: 'fetchListWithoutStore',
    }),

    async search() {
      this.pending = true;

      const { entities = [] } = await this.fetchContextEntitiesWithoutStore({
        params: {
          _filter: getContextWidgetSearchByText(this.searchingText),
        },
      });

      this.entities = entities;
      this.pending = false;
    },
  },
};
</script>

