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
      loader(v-if="contextEntitiesPending")
      div(v-else)
        v-data-table(
        v-model="selected",
        :items="contextEntities",
        :headers="properties",
        :total-items="contextEntitiesMeta.total",
        :pagination.sync="vDataTablePagination",
        item-key="_id",
        select-all,
        hide-actions,
        )
          template(slot="headerCell", slot-scope="props")
            span {{ props.header.text }}
          template(slot="items", slot-scope="props")
            td
              v-checkbox(primary, hide-details, v-model="props.selected")
            td(
            v-for="prop in properties",
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
            pagination(:meta="contextEntitiesMeta", :query.sync="query")
          v-flex(xs2)
            records-per-page(:query.sync="query")
</template>

<script>
import omit from 'lodash/omit';

import ContextSearch from '@/components/other/context/search/context-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import Loader from '@/components/other/context/loader/context-loader.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';

import queryMixin from '@/mixins/query';
import modalMixin from '@/mixins/modal/modal';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import { MODALS } from '@/constants';

import CreateEntity from './actions/context-fab.vue';
import MoreInfos from './more-infos.vue';

/**
 * Entities list
 *
 * @module context
 *
 * @prop {Object} widget - Object representing the widget
 * @prop {Array} properties - List of entities properties
 *
 * @event openSettings#click
 */
export default {
  components: {
    ContextSearch,
    RecordsPerPage,
    CreateEntity,
    MoreInfos,
    Loader,
    Ellipsis,
  },
  mixins: [
    queryMixin,
    modalMixin,
    entitiesContextEntityMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    properties: {
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
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
        'selectedTypes',
      ]);

      query.start = ((this.query.page - 1) * this.query.limit) || 0;

      if (this.query.sortKey) {
        query.sort = [{
          property: this.query.sortKey,
          direction: this.query.sortDir,
        }];
      }

      if (!query._filter) {
        const selectedTypes = this.userPreference.widget_preferences.selectedTypes || [];

        if (selectedTypes.length) {
          query._filter = JSON.stringify({
            $or: selectedTypes.map(type => ({ type })),
          });
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
          action: () => Promise.all(this.selected.map(item => this.removeContextEntity({ id: item._id }))),
        },
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

