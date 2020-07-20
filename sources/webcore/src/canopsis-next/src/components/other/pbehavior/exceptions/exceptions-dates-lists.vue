<template lang="pug">
  div
    v-layout(row)
      search-field(@submit="handleSearch", @clear="handleSearchClear")
    v-data-table(
      v-field="exceptions",
      :headers="headers",
      :items="pbehaviorExceptions",
      :loading="pbehaviorExceptionsPending",
      :total-items="pbehaviorExceptionsMeta.total_count",
      :pagination.sync="query",
      item-key="id",
      select-all
    )
      template(slot="items", slot-scope="props")
        tr
          td
            v-checkbox-functional(v-model="props.selected", primary, hide-details)
          td {{ props.item.name }}
</template>

<script>
import { isEqual } from 'lodash';

import exceptionMixin from '@/mixins/entities/pbehavior/exception';
import exceptionQueryMixin from '@/mixins/pbehavior/exception/query';

import SearchField from '@/components/forms/fields/search-field.vue';

export default {
  components: { SearchField },
  mixins: [
    exceptionMixin,
    exceptionQueryMixin,
  ],
  model: {
    prop: 'exceptions',
    event: 'input',
  },
  props: {
    exceptions: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    headers() {
      return [{ text: this.$t('common.name'), value: 'name', sortable: false }];
    },
  },
  watch: {
    query(query, oldQuery) {
      if (!isEqual(query, oldQuery)) {
        this.fetchList();
      }
    },
  },

  methods: {
    async fetchList() {
      this.fetchPbehaviorExceptionsList({ params: this.getQuery() });
    },
  },
};
</script>
