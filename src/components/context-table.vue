<template lang="pug">
  div(v-if="!pending")
    basic-list(:items="items")
      tr.container(slot="header")
        th.box(v-for="columnProperty in contextProperties")
          span {{ columnProperty.text }}
          th.box
      tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in contextProperties") {{ item.props | get(property.value, property.filter) }}
        td.box
      tr.container(slot="expandedRow", slot-scope="item")
    pagination(:meta="meta", :limit="limit")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/basic-component/basic-list.vue';
import PaginationMixin from '@/mixins/pagination';

const { mapActions, mapGetters } = createNamespacedHelpers('context');

export default {
  name: 'context-table',
  components: { BasicList },
  mixins: [PaginationMixin],
  props: {
    contextProperties: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'pending',
    ]),
    queries() {
      const queries = {};
      queries.limit = this.limit;
      queries.start = ((this.$route.query.page - 1) * this.limit) || 0;
      return queries;
    },
  },
  methods: {
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
  },
};
</script>

<style scoped>
  th {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  td {
    overflow-wrap: break-word;
  }
  .container {
    display: flex;
  }
  .box{
    width: 10%;
    flex: 1;
  }
</style>
