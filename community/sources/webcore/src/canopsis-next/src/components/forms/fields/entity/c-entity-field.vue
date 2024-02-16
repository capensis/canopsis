<template>
  <c-lazy-search-field
    v-field="value"
    :label="selectLabel"
    :loading="pending"
    :limit="limit"
    :items="entities"
    :name="name"
    :has-more="hasMoreEntities"
    :item-text="itemText"
    :item-value="itemValue"
    :item-disabled="itemDisabled"
    :disabled="disabled"
    :required="required"
    :clearable="clearable"
    :autocomplete="autocomplete"
    :return-object="returnObject"
    @fetch="fetchEntities"
    @fetch:more="fetchMoreEntities"
    @update:search="updateSearch"
  />
</template>

<script>
import { isArray, keyBy, pick } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { BASIC_ENTITY_TYPES } from '@/constants';
import { PAGINATION_LIMIT } from '@/config';

import { formArrayMixin } from '@/mixins/form';

const { mapActions: entityMapActions } = createNamespacedHelpers('entity');

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Object],
      default: '',
    },
    name: {
      type: String,
      default: 'entities',
    },
    label: {
      type: String,
      required: false,
    },
    itemText: {
      type: [String, Function],
      default: '_id',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    entityTypes: {
      type: Array,
      default: () => Object.values(BASIC_ENTITY_TYPES),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    autocomplete: {
      type: Boolean,
      default: false,
    },
    itemDisabled: {
      type: [String, Array, Function],
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      entitiesById: {},
      pending: false,
      pageCount: 1,

      query: {
        page: 1,
        search: null,
      },
    };
  },
  computed: {
    entities() {
      return Object.values(this.entitiesById);
    },

    hasMoreEntities() {
      return this.pageCount > this.query.page;
    },

    selectLabel() {
      if (this.label) {
        return this.label;
      }

      if (isArray(this.value)) {
        return this.$tc('common.entity', this.value.length);
      }

      return this.$tc('common.entity');
    },
  },
  methods: {
    ...entityMapActions({ fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore' }),

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        type: this.entityTypes,
      };
    },

    async fetchEntities() {
      try {
        this.pending = true;

        const { data: entities, meta } = await this.fetchContextEntitiesListWithoutStore({
          params: this.getQuery(),
        });

        this.pageCount = meta.page_count;

        const currentEntities = this.returnObject
          ? keyBy(this.value, '_id')
          : pick(this.entitiesById, isArray(this.value) ? this.value : [this.value]);

        this.entitiesById = {
          ...(this.query.page !== 1 ? this.entitiesById : {}),
          ...keyBy(entities, '_id'),
          ...currentEntities,
        };
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },

    fetchMoreEntities() {
      this.query.page += 1;

      this.fetchEntities();
    },

    updateSearch(search) {
      this.query.search = search;
      this.query.page = 1;

      this.fetchEntities();
    },
  },
};
</script>
