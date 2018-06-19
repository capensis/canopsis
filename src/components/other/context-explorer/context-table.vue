<template lang="pug">
  div
    basic-list(:items="items")
      tr.container(slot="header")
        th.box(v-for="columnProperty in contextProperties")
          span {{ columnProperty.text }}
          th.box
      tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in contextProperties") {{ item.props | get(property.value, property.filter) }}
        td.box
          v-btn(@click.stop="deleteEntity(item)", icon, small)
            v-icon delete
      tr.container(slot="expandedRow", slot-scope="item")
    pagination(:meta="meta", :limit="limit")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/tables/basic-list.vue';
import PaginationMixin from '@/mixins/pagination';
import omit from 'lodash/omit';

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
  },
  methods: {
    ...mapActions({
      fetchListAction: 'fetchList',
      remove: 'remove',
    }),
    getQuery() {
      const query = omit(this.$route.query, ['page']);

      query.limit = this.limit;
      query.start = ((this.$route.query.page - 1) * this.limit) || 0;
      return query;
    },
    deleteEntity(item) {
      this.remove({ ids: `"${item.props._id}"` });
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
</style>
