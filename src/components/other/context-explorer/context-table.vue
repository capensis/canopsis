<template lang="pug">
  div
    records-per-page
    basic-list(:items="contextEntities")
      tr.container(slot="header")
        th.box(v-for="columnProperty in contextProperties")
          span {{ columnProperty.text }}
          list-sorting(:column="columnProperty.value", class="blue--text")
        th.box
      tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in contextProperties") {{ item.props | get(property.value, property.filter) }}
        td.box
    pagination(:meta="contextEntitiesMeta", :limit="limit")
    create-entity.fab
</template>

<script>
import BasicList from '@/components/tables/basic-list.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
import CreateEntity from '@/components/other/context-explorer/actions/context-fab.vue';
import PaginationMixin from '@/mixins/pagination';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import omit from 'lodash/omit';
import ContextEntityMixin from '@/mixins/context';

export default {
  components: {
    BasicList,
    RecordsPerPage,
    ListSorting,
    CreateEntity,
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
      const query = omit(this.$route.query, ['page', 'sort_dir', 'sort_key']);
      query.limit = this.limit;
      query.start = ((this.$route.query.page - 1) * this.limit) || 0;

      if (this.$route.query.sort_key) {
        query.sort = [{
          property: this.$route.query.page.sort_key,
          direction: this.$route.query.page.sort_dir ? this.$route.query.page.sort_dir : 'ASC',
        }];
      }

      return query;
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

  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
