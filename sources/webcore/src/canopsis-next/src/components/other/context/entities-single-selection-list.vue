<template lang="pug">
  div
    search-field(v-model="searchingText", @submit="fetchList")
    v-data-table(
    :no-data-text="$t('tables.noData')",
    :headers="headers",
    :items="contextEntities",
    :loading="pending",
    item-key="_id"
    )
      template(slot="items", slot-scope="{ item }")
        td
          v-radio-group(:value="value", hide-details, @change="updateModel")
            v-radio(:value="item._id")
        td.text-xs-left {{ item.name }}
        td.text-xs-left {{ item._id}}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import formBaseMixin from '@/mixins/form/base';

import { getContextSearchByText } from '@/helpers/widget-search';

import SearchField from '@/components/forms/fields/search-field.vue';

const { mapGetters, mapActions } = createNamespacedHelpers('entity');

export default {
  components: { SearchField },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: null,
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
      selectedEntity: [],
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
          text: '',
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.name'),
          sortable: false,
        },
        {
          text: this.$t('tables.contextList.id'),
          sortable: false,
        },
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchContextEntities: 'fetchGeneralList',
    }),

    fetchList() {
      this.fetchContextEntities({
        params: {
          _filter: this.filterPreparer(this.searchingText),
        },
      });
    },
  },
};
</script>

