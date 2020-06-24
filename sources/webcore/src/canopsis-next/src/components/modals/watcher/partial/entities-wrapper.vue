<template lang="pug">
  div
    div.mt-2(v-for="watcherEntity in slicedWatcherEntities", :key="watcherEntity._id")
      watcher-entity(
        :watcherId="watcher.entity_id",
        :isWatcherOnPbehavior="watcher.active_pb_watcher",
        :entity="watcherEntity",
        :template="entityTemplate",
        :entityNameField="entityNameField",
        @add:event="$emit('add:event', $event)"
      )
    .float-clear
    v-layout.white(v-if="pagination.total", align-center)
      v-flex(xs10)
        pagination(
          :page="pagination.page",
          :limit="pagination.limit",
          :total="pagination.total",
          @input="updatePage"
        )
      v-spacer
      v-flex(xs2)
        records-per-page(:value="pagination.limit", @input="updateItemsPerPage")
</template>

<script>
import { orderBy } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';

import Pagination from '@/components/tables/pagination.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

import WatcherEntity from './entity.vue';

export default {
  components: {
    Pagination,
    RecordsPerPage,
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
