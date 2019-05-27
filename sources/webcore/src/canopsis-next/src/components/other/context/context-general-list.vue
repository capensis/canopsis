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
    :select-all="!single",
    item-key="_id"
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
  props: {
    single: {
      type: Boolean,
      default: false,
    },
    filterPreparer: {
      type: Function,
      default: getContextSearchByText,
    },
    initialSearchingText: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      searchingText: this.initialSearchingText,
      selectedEntities: [],
    };
  },
  computed: {
    ...mapGetters({
      contextEntities: 'itemsGeneralList',
      pending: 'pendingGeneralList',
    }),

    headers() {
      const headers = [
        {
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
      ];

      if (!this.single) {
        return [
          ...headers,

          {
            text: this.$t('common.actionsLabel'),
            sortable: false,
          },
        ];
      }

      return [
        {
          text: '',
          sortable: false,
        },

        ...headers,
      ];
    },
  },
  mounted() {
    if (this.single) {
      this.submit();
    }
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchGeneralList',
    }),

    submit() {
      this.fetchContextEntities({
        params: {
          _filter: this.filterPreparer(this.searchingText),
        },
      });
    },
  },
};
</script>

