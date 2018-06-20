<template lang="pug">
  div
    context-search
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
    create-entity.fab
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/tables/basic-list.vue';
import ContextSearch from '@/components/other/context-explorer/searching/context-search.vue';
import CreateEntity from '@/components/other/context-explorer/actions/context-fab.vue';
import PaginationMixin from '@/mixins/pagination';
import omit from 'lodash/omit';

const { mapActions, mapGetters } = createNamespacedHelpers('context');

export default {
  components: { BasicList, ContextSearch, CreateEntity },
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
  },
  methods: {
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    getQuery() {
      const query = omit(this.$route.query, ['page']);

      query.limit = this.limit;
      query.start = ((this.$route.query.page - 1) * this.limit) || 0;
      return query;
    },
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
  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
