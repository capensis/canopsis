<template lang="pug">
  v-container
    v-layout.white(justify-space-between, align-center)
      v-flex(xs12, md4)
        context-search
      v-flex.ml-4(xs4)
        v-btn(v-show="selected.length", @click.stop="deleteEntities", icon, small)
          v-icon delete
      v-flex(xs2)
        v-btn(icon, @click.prevent="$emit('openSettings')")
          v-icon settings
    transition(name="fade", mode="out-in")
      loader(v-if="pending")
      div(v-else)
        v-data-table(
          v-model="selected",
          :items="contextEntities",
          :headers="contextProperties",
          item-key="_id",
          :total-items="meta.total",
          :pagination.sync="pagination",
          select-all,
          hide-actions,
        )
          template(slot="headerCell", slot-scope="props")
              span {{ props.header.text }}
          template(slot="items", slot-scope="props")
            td
              v-checkbox(primary, hide-details, v-model="props.selected")
            td(
              v-for="prop in contextProperties",
              @click="props.expanded = !props.expanded"
            )
              ellipsis(:text="$options.filters.get(props.item,prop.value) || ''",
                       :maxLetters="prop.maxLetters || MAX_LETTERS")
            td
              v-btn(@click.stop="deleteEntity(props.item)", icon, small)
                v-icon delete
          template(slot="expand", slot-scope="props")
        v-layout.white(align-center)
          v-flex(xs10)
            pagination(:meta="meta", :limit="limit")
          v-flex(xs2)
            records-per-page
        create-entity.fab
</template>

<script>
import omit from 'lodash/omit';
import { createNamespacedHelpers } from 'vuex';

import ContextSearch from '@/components/other/context/search/context-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import Loader from '@/components/other/context/loader/context-loader.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';

import paginationMixin from '@/mixins/pagination';
import modalMixin from '@/mixins/modal/modal';
import contextEntityMixin from '@/mixins/context';

import { MODALS } from '@/constants';
import { MAX_LETTERS } from '@/config';

import CreateEntity from './actions/context-fab.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('context');

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
    ContextSearch,
    RecordsPerPage,
    CreateEntity,
    Loader,
    Ellipsis,
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
      pagination: {},
      MAX_LETTERS,
    };
  },
  computed: {
    ...mapGetters(['items', 'meta', 'pending']),
  },
  watch: {
    pagination: {
      handler(e) {
        this.$router.push({
          query: {
            ...this.$route.query,
            sort_key: e.sortBy,
            sort_dir: e.descending ? 'DESC' : 'ASC',
          },
        });
      },
    },
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
    deleteEntity(item) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.remove({ id: item._id }),
        },
      });
    },
    deleteEntities() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => Promise.all(this.selected.map(item => this.remove({ id: item._id }))),
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
.fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
.fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>

