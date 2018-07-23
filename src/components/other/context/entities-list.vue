<template lang="pug">
  v-container
    v-layout.white(justify-space-between, align-center)
      v-flex(xs12, md4)
        context-search(:query.sync="query")
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
          :pagination.sync="vDataTablePagination",
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
              ellipsis(
                :text="$options.filters.get(props.item,prop.value) || ''",
                :maxLetters="prop.maxLetters"
              )
            td
              v-btn(@click.stop="editEntity(props.item)", icon, small)
                v-icon edit
              v-btn(@click.stop="deleteEntity(props.item)", icon, small)
                v-icon delete
          template(slot="expand", slot-scope="props")
            more-infos(:item="props")
        v-layout.white(align-center)
          v-flex(xs10)
            pagination(:meta="meta", :query.sync="query")
          v-flex(xs2)
            records-per-page(:query.sync="query")
</template>

<script>
import find from 'lodash/find';
import omit from 'lodash/omit';
import { createNamespacedHelpers } from 'vuex';

import ContextSearch from '@/components/other/context/search/context-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import Loader from '@/components/other/context/loader/context-loader.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';

import paginationMixin from '@/mixins/pagination';
import modalMixin from '@/mixins/modal/modal';
import contextEntityMixin from '@/mixins/context/list';
import AddInfoObject from '@/components/other/context/actions/manage-info-object.vue';
import { MODALS, ENTITIES_TYPES } from '@/constants';

import CreateEntity from './actions/context-fab.vue';
import MoreInfos from './more-infos.vue';

const { mapGetters } = createNamespacedHelpers('entity');
const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

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
    AddInfoObject,
    ContextSearch,
    RecordsPerPage,
    CreateEntity,
    MoreInfos,
    Loader,
    Ellipsis,
  },
  mixins: [
    paginationMixin,
    contextEntityMixin,
    modalMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
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
      queryPrefix: 'context',
    };
  },
  computed: {
    ...mapGetters(['items', 'meta', 'pending']),

    ...entitiesMapGetters(['getItem']),

    userPreference() {
      return this.getItem(ENTITIES_TYPES.userPreference, `${this.widget.id}_root`); // TODO: fix it
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.$route.query, ['page', 'sort_dir', 'sort_key']);
      query.limit = this.query.limit;
      query.start = ((this.query.page - 1) * this.query.limit) || 0;

      if (this.query.sort_key) {
        query.sort = [{
          property: this.query.sort_key,
          direction: this.query.sort_dir ? this.query.sort_dir : 'ASC',
        }];
      }

      // TODO: fix it
      if (this.userPreference) {
        const filter = find(this.userPreference.widget_preferences.user_filters, { title: 'default_type_filter' });

        if (filter) {
          query._filter = filter.filter;
        } else {
          delete query._filter;
        }
      }

      return query;
    },
    editEntity(item) {
      this.showModal({
        name: MODALS.createEntity,
        config: {
          title: this.$t('modals.createEntity.editTitle'),
          item,
        },
      });
    },
    deleteEntity(item) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeContextEntity({ id: item._id }),
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

