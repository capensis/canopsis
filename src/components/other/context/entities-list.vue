<template lang="pug">
  div
    v-layout.white(justify-space-between, align-center)
      v-flex
        context-search(:query.sync="query")
      v-flex
        pagination(v-if="hasColumns", :meta="contextEntitiesMeta", :query.sync="query", type="top")
      v-flex
        v-select(
        :label="$t('settings.selectAFilter')",
        :items="viewFilters",
        @input="updateSelectedFilter",
        :value="mainFilter",
        item-text="title",
        item-value="filter",
        return-object,
        clearable
        )
      v-flex.ml-4
        div(v-show="selected.length")
          v-btn(@click.stop="deleteEntities", icon, small)
            v-icon delete
          v-btn(@click.stop="addPbehaviors()", icon, small)
            v-icon pause
      v-flex
        context-fab
    no-columns-table(v-if="!hasColumns")
    div(v-else)
      v-data-table(
      v-model="selected",
      :items="contextEntities",
      :headers="headers",
      :loading="contextEntitiesPending",
      :total-items="contextEntitiesMeta.total",
      :pagination.sync="vDataTablePagination",
      item-key="_id",
      select-all,
      hide-actions,
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="headerCell", slot-scope="props")
          span {{ props.header.text }}
        template(slot="items", slot-scope="props")
          td
            v-checkbox(primary, hide-details, v-model="props.selected")
          td(
          v-for="column in columns",
          @click="props.expanded = !props.expanded"
          )
            ellipsis(
            :text="props.item | get(column.value, null, '')",
            :maxLetters="column.maxLetters"
            )
          td
            v-btn(@click.stop="editEntity(props.item)", icon, small)
              v-icon edit
            v-btn(@click.stop="duplicateEntity(props.item)", icon, small)
              v-icon file_copy
            v-btn(@click.stop="deleteEntity(props.item)", icon, small)
              v-icon delete
            v-btn(@click.stop="addPbehaviors(props.item._id)", icon, small)
              v-icon pause
        template(slot="expand", slot-scope="props")
          more-infos(:item="props.item")
      v-layout.white(align-center)
        v-flex(xs10)
          pagination(:meta="contextEntitiesMeta", :query.sync="query")
        v-flex(xs2)
          records-per-page(:query.sync="query")
</template>

<script>
import omit from 'lodash/omit';
import isString from 'lodash/isString';

import { MODALS, ENTITIES_TYPES } from '@/constants';

import ContextSearch from '@/components/other/context/search/context-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import Ellipsis from '@/components/tables/ellipsis.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';

import modalMixin from '@/mixins/modal/modal';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesWatcherMixin from '@/mixins/entities/watcher';
import filterSelectMixin from '@/mixins/filter-select';

import ContextFab from './actions/context-fab.vue';
import MoreInfos from './more-infos.vue';

/**
 * Entities list
 *
 * @module context
 *
 * @prop {Object} widget - Object representing the widget
 * @prop {Array} columns - List of entities columns
 *
 * @event openSettings#click
 */
export default {
  components: {
    ContextSearch,
    RecordsPerPage,
    MoreInfos,
    Ellipsis,
    ContextFab,
    NoColumnsTable,
  },
  mixins: [
    modalMixin,
    widgetQueryMixin,
    widgetColumnsMixin,
    entitiesContextEntityMixin,
    entitiesWatcherMixin,
    filterSelectMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    headers() {
      if (this.hasColumns) {
        return [...this.columns, { text: '', sortable: false }];
      }

      return [];
    },
  },
  methods: {
    getQuery() {
      const query = omit(this.query, [
        'page',
        'sortKey',
        'sortDir',
        'mainFilter',
        'searchFilter',
        'typesFilter',
      ]);

      query.start = ((this.query.page - 1) * this.query.limit) || 0;

      if (this.query.sortKey) {
        query.sort = [{
          property: this.query.sortKey,
          direction: this.query.sortDir,
        }];
      }

      const filters = ['mainFilter', 'searchFilter', 'typesFilter'].reduce((acc, filterKey) => {
        const queryFilter = isString(this.query[filterKey]) ? JSON.parse(this.query[filterKey]) : this.query[filterKey];

        if (queryFilter) {
          acc.push(queryFilter);
        }

        return acc;
      }, []);

      if (filters.length) {
        query._filter = {
          $and: filters,
        };
      }

      return query;
    },
    editEntity(item) {
      if (item.type === ENTITIES_TYPES.watcher) {
        this.showModal({
          name: MODALS.createWatcher,
          config: {
            title: 'modals.createWatcher.editTitle',
            item,
            action: watcher => this.editWatcherWithPopup(watcher),
          },
        });
      } else {
        this.showModal({
          name: MODALS.createEntity,
          config: {
            title: 'modals.createEntity.editTitle',
            item,
            action: entity => this.updateContextEntityWithPopup(entity),
          },
        });
      }
    },
    duplicateEntity(item) {
      if (item.type === 'watcher') {
        this.showModal({
          name: MODALS.createWatcher,
          config: {
            title: 'modals.createWatcher.duplicateTitle',
            item,
            isDuplicating: true,
            action: watcher => this.duplicateWatcherWithPopup(watcher),
          },
        });
      } else {
        this.showModal({
          name: MODALS.createEntity,
          config: {
            title: 'modals.createEntity.duplicateTitle',
            item,
            isDuplicating: true,
            action: entity => this.duplicateContextEntityWithPopup(entity),
          },
        });
      }
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
    addPbehaviors(itemId) {
      this.showModal({
        name: MODALS.createPbehavior,
        config: {
          itemsType: ENTITIES_TYPES.entity,
          itemsIds: itemId ? [itemId] : this.selected,
        },
      });
    },
    fetchList() {
      if (this.hasColumns) {
        this.fetchContextEntitiesList({
          widgetId: this.widget._id,
          params: this.getQuery(),
        });
      }
    },
    /**
     * Surcharge updateSelectedFilter method from filterSelectMixin
     * to adapt it to context API specification
     */
    async updateSelectedFilter(value = {}) {
      this.createUserPreference({
        userPreference: {
          ...this.userPreference,

          widget_preferences: {
            ...this.userPreference.widget_preferences,

            mainFilter: value,
          },
        },
      });

      this.query = {
        ...this.query,

        mainFilter: value.filter,
      };
    },
  },
};
</script>
