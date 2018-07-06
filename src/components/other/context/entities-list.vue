<template lang="pug">
  div
    v-layout(justify-space-between, align-center)
      v-flex.ml-4(xs4)
        v-btn(v-show="selected.length", @click.stop="deleteEntities", icon, small)
          v-icon delete
      v-flex(xs2)
        v-btn(icon, @click.prevent="$emit('openSettings')")
          v-icon settings
    context-search
    records-per-page
    basic-list(:items="contextEntities", :selected.sync="selected", :pending="pending")
      loader(slot="loader")
      tr.container(slot="header")
        th.box(v-for="columnProperty in contextProperties")
          span {{ columnProperty.text }}
          list-sorting.blue--text(:column="columnProperty.value")
        th.box
      tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in contextProperties") {{ item.props | get(property.value, property.filter) }}
        td.box
          v-btn(@click.stop="editEntity(item)", icon, small)
            v-icon edit
          v-btn(@click.stop="deleteEntity(item)", icon, small)
            v-icon delete
      tr.container(slot="expandedRow", slot-scope="item")
    pagination(:meta="contextEntitiesMeta", :limit="limit")
    create-entity.fab
</template>

<script>
import omit from 'lodash/omit';
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/tables/basic-list.vue';
import ContextSearch from '@/components/other/context/search/context-search.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import Loader from '@/components/other/context/loader/context-loader.vue';

import paginationMixin from '@/mixins/pagination';
import modalMixin from '@/mixins/modal/modal';
import contextEntityMixin from '@/mixins/context';

import { MODALS } from '@/constants';

import CreateEntity from './actions/context-fab.vue';

const { mapActions } = createNamespacedHelpers('context');

/**
 * Entities list
 *
 * @module context
 *
 * @prop {Array} [contextProperties] - List of entities properties
 *
 * @event openSettings#click
 */
export default {
  components: {
    BasicList,
    ContextSearch,
    RecordsPerPage,
    ListSorting,
    CreateEntity,
    Loader,
  },
  mixins: [
    paginationMixin,
    contextEntityMixin,
    modalMixin,
  ],
  props: {
    contextProperties: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  methods: {
    ...mapActions({
      fetchListAction: 'fetchList',
      remove: 'remove',
    }),
    getQuery() {
      const query = omit(this.$route.query, ['page', 'sort_dir', 'sort_key']);
      query.limit = this.limit;
      query.start = ((this.$route.query.page - 1) * this.limit) || 0;

      if (this.$route.query.sort_key) {
        query.sort = [{
          property: this.$route.query.sort_key,
          direction: this.$route.query.sort_dir ? this.$route.query.sort_dir : 'ASC',
        }];
      }

      return query;
    },
    editEntity(item) {
      this.showModal({
        name: MODALS.createEntity,
        config: {
          item,
        },
      });
    },
    deleteEntity(item) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.remove({ id: item.props._id }),
        },
      });
    },
    deleteEntities() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => Promise.all(this.selected.map(id => this.remove({ id }))),
        },
      });
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
