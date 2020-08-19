<template lang="pug">
  div
    v-layout(row)
      search-field(@submit="updateSearchHandler", @clear="clearSearchHandler")
    v-data-table(
      v-field="exceptions",
      :headers="headers",
      :items="pbehaviorDatesExceptions",
      :loading="pbehaviorDatesExceptionsPending",
      :total-items="pbehaviorDatesExceptionsMeta.total_count",
      :pagination.sync="query",
      item-key="_id",
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

import entitiesPbehaviorDatesExceptionMixin from '@/mixins/entities/pbehavior/dates-exceptions';
import pbehaviorQueryMixin from '@/mixins/pbehavior/query';

import SearchField from '@/components/forms/fields/search-field.vue';

export default {
  components: { SearchField },
  mixins: [
    entitiesPbehaviorDatesExceptionMixin,
    pbehaviorQueryMixin,
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
      this.fetchPbehaviorDatesExceptionsList({ params: this.getQuery() });
    },
  },
};
</script>
