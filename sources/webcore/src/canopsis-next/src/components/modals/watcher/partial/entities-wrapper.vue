<template lang="pug">
  div
    div.mt-2(v-for="watcherEntity in slicedWatcherEntities", :key="watcherEntity.key")
      watcher-entity(
        :watcher-id="watcher._id",
        :is-watcher-on-pbehavior="watcher.active_pb_watcher",
        :entity="watcherEntity",
        :template="entityTemplate",
        :entity-name-field="entityNameField",
        @add:event="$emit('add:event', $event)"
      )
    .float-clear
    table-pagination(
      :total-items="pagination.total",
      :rows-per-page="pagination.limit",
      :page="pagination.page",
      @update:page="updatePage",
      @update:rows-per-page="updateItemsPerPage"
    )
</template>

<script>
import { orderBy } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import TablePagination from '@/components/other/table/table-pagination.vue';

import WatcherEntity from './entity.vue';

export default {
  components: {
    TablePagination,
    WatcherEntity,
  },
  props: {
    watcher: {
      type: Object,
      required: true,
    },
    watcherEntities: {
      type: Array,
      default: () => [],
    },
    entityNameField: {
      type: String,
      default: 'name',
    },
    entityTemplate: {
      type: String,
      default: '',
    },
    itemsPerPage: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
  },
  data() {
    return {
      query: {
        page: 1,
        limit: this.itemsPerPage,
      },
    };
  },
  computed: {
    pagination() {
      return {
        page: this.query.page,
        limit: this.query.limit,
        first: (this.query.page - 1) * this.query.limit,
        last: this.query.page * this.query.limit,
        total: this.orderedWatcherEntities.length,
      };
    },
    orderedWatcherEntities() {
      const preparedEntityNameField = this.entityNameField.replace(/^entity\./, '');

      return orderBy(this.watcherEntities, ['state.val', preparedEntityNameField], ['desc', 'asc']);
    },
    slicedWatcherEntities() {
      const { first, last } = this.pagination;

      return this.orderedWatcherEntities.slice(first, last);
    },
  },
  methods: {
    updatePage(page) {
      this.query.page = page;
    },

    updateItemsPerPage(limit) {
      this.query = {
        page: 1,
        limit,
      };
    },
  },
};
</script>
