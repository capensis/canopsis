<template lang="pug">
  div
    v-layout(row)
      search-field(@submit="updateSearchHandler", @clear="clearSearchHandler")
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
      template(slot="items", slot-scope="props")
        tr
          td
            v-checkbox-functional(v-model="props.selected", primary, hide-details)
          td {{ props.item.name }}
</template>

<script>
import { isEqual, omit } from 'lodash';

import { PLANNING_TABS } from '@/constants';

import entitiesPbehaviorExceptionMixin from '@/mixins/entities/pbehavior/exceptions';
import pbehaviorQueryMixin from '@/mixins/pbehavior/query';

import SearchField from '@/components/forms/fields/search-field.vue';

export default {
  components: { SearchField },
  mixins: [
    entitiesPbehaviorExceptionMixin,
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
    queryId: {
      type: String,
      default: PLANNING_TABS.exceptions,
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
      this.fetchPbehaviorExceptionsList({ params: this.query });
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
