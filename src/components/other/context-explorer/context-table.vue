<template lang="pug">
  div
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
          more-infos(:item="item")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { getQueryContext } from '@/helpers/pagination';

import BasicList from '@/components/tables/basic-list.vue';
import MoreInfos from '@/components/other/context-explorer/more-infos.vue';
import getProp from 'lodash/get';

const { mapActions, mapGetters } = createNamespacedHelpers('context');

export default {
  name: 'context-table',
  components: { BasicList, MoreInfos },
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
  watch: {
    $route() {
      this.fetchList();
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getQuery: getQueryContext,
    getProp,
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    fetchList() {
      this.fetchListAction({
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
  .box{
    width: 10%;
    flex: 1;
  }
</style>
