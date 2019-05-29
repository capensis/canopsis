<template lang="pug">
  div
    v-toolbar.toolbar.white(dense, flat)
      v-text-field(
      v-model="searchingText",
      :label="$t('common.search')",
      hide-details,
      single-line,
      @keyup.enter="submit"
      )
      v-btn(icon, @click="submit")
        v-icon search
    v-btn.green.white--text(
    v-show="selectedEntities.length",
    @click="$emit('update:selectedIds', selectedEntities)"
    ) Add selection
    v-data-table(
    v-model="selectedEntities",
    :no-data-text="$t('tables.noData')",
    :headers="headers",
    :items="contextEntities",
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

import { getContextSearchByText } from '@/helpers/widget-search';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

export default {
  data() {
    return {
      searchingText: '',
      selectedEntities: [],
    };
  },
  computed: {
    ...mapGetters({
      contextEntities: 'itemsGeneralList',
      pending: 'pendingGeneralList',
    }),

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
      fetchContextEntities: 'fetchGeneralList',
    }),

    submit() {
      this.fetchContextEntities({
        params: {
          _filter: getContextSearchByText(this.searchingText),
        },
      });
    },
  },
};
</script>

