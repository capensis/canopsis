<template lang="pug">
  div
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header")
          th.box(v-for="columnProperty in contextProperty")
            span {{ columnProperty.text }}
            th.box
        tr.container(slot="row" slot-scope="item")
          td.box(v-for="columnProperty in contextProperty") {{ getProp(item.props, columnProperty.value) }}
          td.box
        tr.container(slot="expandedRow", slot-scope="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { PAGINATION_LIMIT } from '@/config';
import getQuery from '@/helpers/pagination';
import BasicList from '@/components/BasicComponent/basic-list.vue';
import getProp from 'lodash/get';

const { mapActions, mapGetters } = createNamespacedHelpers('context');
export default {
  name: 'context-table',
  components: { BasicList },
  props: {
    contextProperty: {
      type: Array,
      default() {
        return [];
      },
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'pending',
    ]),
  },
  watch: {
    $route() {
      this.fetchList();
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getQuery,
    getProp,
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    fetchList() {
      this.fetchListAction({
        params: { sort: 'ASC' },
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
  .box{
    width: 10%;
    flex: 1;
  }
</style>
