<template lang="pug">
  div
    records-per-page
    basic-list(:items="contextEntities", :needExpand="false")
      tr.container(slot="header")
        th.box(v-for="columnProperty in contextProperties")
          span {{ columnProperty.text }}
          list-sorting(:column="columnProperty.value", class="blue--text")
        th.box
      tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in contextProperties") {{ item.props | get(property.value, property.filter) }}
        td.box
    pagination(:meta="contextEntitiesMeta", :limit="limit")
</template>

<script>
import BasicList from '@/components/tables/basic-list.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
import PaginationMixin from '@/mixins/pagination';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import omit from 'lodash/omit';
import ContextEntityMixin from '@/mixins/context';

export default {
  name: 'context-table',
  components: {
    BasicList,
    RecordsPerPage,
    ListSorting,
  },
  mixins: [
    PaginationMixin,
    ContextEntityMixin,
  ],
  props: {
    contextProperties: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.$route.query, ['page']);
      const data = {};
      data.limit = this.limit;
      data.start = ((this.$route.query.page - 1) * this.limit) || 0;

      if (query.sort_key) {
        data.sort = [{
          property: query.sort_key,
          direction: query.sort_dir ? query.sort_dir : 'ASC',
        }];
      }

      return data;
    },
    fetchList() {
      this.fetchContextEntities({
        params: this.getQuery(),
      });
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

  .box {
    width: 10%;
    flex: 1;
  }
</style>
