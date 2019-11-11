<template lang="pug">
  div
    v-toolbar.toolbar.white(dense, flat)
      v-text-field(
        data-test="searchEntity",
        label="Search",
        v-model="searchingText",
        hide-details,
        single-line,
        @keyup.enter="submit"
      )
      v-btn(data-test="searchEntityButton", icon, @click="submit")
        v-icon search
    v-btn.green.white--text(
      data-test="addCollectionEntities",
      v-show="selectedEntities.length",
      @click="$emit('update:selectedIds', selectedEntities)"
    ) Add selection
    v-data-table(
      data-test="contextEntitiesTable",
      :no-data-text="$t('tables.noData')",
      :headers="headers",
      :items="contextEntities",
      :loading="pending",
      v-model="selectedEntities",
      select-all,
      item-key="_id"
    )
      template(slot="items", slot-scope="props")
        td
          v-checkbox(
            :data-test="`contextRowCheckbox-${props.item._id}`",
            v-model="props.selected",
            :value="props._id",
            primary,
            hide-details
          )
        td.text-xs-left {{ props.item.name }}
        td.text-xs-left {{ props.item._id}}
        td
          v-btn(
            :data-test="`contextRowAdd-${props.item._id}`",
            icon,
            @click="$emit('update:selectedIds', [props.item])"
          )
            v-icon add
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { getContextWidgetSearchByText } from '@/helpers/entities-search';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

export default {
  data() {
    return {
      searchingText: '',
      selectedEntities: [],
      headers: [
        {
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
      ],
    };
  },
  computed: {
    ...mapGetters({
      contextEntities: 'itemsGeneralList',
      pending: 'pendingGeneralList',
    }),
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchGeneralList',
    }),
    submit() {
      this.fetchContextEntities({
        params: {
          _filter: getContextWidgetSearchByText(this.searchingText),
        },
      });
    },
  },
};
</script>

