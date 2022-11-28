<template lang="pug">
  div
    v-layout(row)
      c-search-field(@submit="updateSearchHandler", @clear="clearSearchHandler")
    v-data-table(
      v-field="exceptions",
      :headers="headers",
      :items="pbehaviorExceptions",
      :loading="pbehaviorExceptionsPending",
      :total-items="pbehaviorExceptionsMeta.total_count",
      :pagination.sync="pagination",
      item-key="_id",
      select-all
    )
      template(#items="{ selected, item }")
        tr
          td
            v-checkbox-functional(v-model="selected", primary, hide-details)
          td {{ item.name }}
</template>

<script>
import { omit } from 'lodash';

import entitiesPbehaviorExceptionMixin from '@/mixins/entities/pbehavior/exceptions';
import { localQueryMixin } from '@/mixins/query-local/query';

export default {
  mixins: [
    entitiesPbehaviorExceptionMixin,
    localQueryMixin,
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
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.fetchPbehaviorExceptionsList({ params: this.getQuery() });
    },

    updateSearchHandler(search) {
      this.$emit('update:pagination', { ...this.pagination, search });
    },

    clearSearchHandler() {
      this.$emit('update:pagination', omit(this.pagination, ['search']));
    },
  },
};
</script>
