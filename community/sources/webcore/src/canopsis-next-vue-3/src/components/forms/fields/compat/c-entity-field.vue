<template>
  <c-lazy-search-field
    :value="value"
    v-bind="$attrs"
    :label="selectLabel"
    :loading="pending"
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
    with-type
    @update:value="updateValue"
    @fetch="fetchEntities"
    @fetch:more="fetchMoreEntities"
    @update:search="updateSearch"
  >
    <template v-if="$slots.selection" #selection="slotProps">
      <slot name="selection" v-bind="slotProps" />
    </template>
    <template v-if="$slots.item" #item="slotProps">
      <slot name="item" v-bind="slotProps" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { isArray, keyBy, pick, isEqual } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { BASIC_ENTITY_TYPES } from '@/constants';
import { useEntity } from '@/hooks/store/modules/entity';

export default {
  inheritAttrs: false,
  emits: ['input', 'update:value', 'update:modelValue'],
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
        return this.$t('common.entity', this.value.length);
      }

      return this.$t('common.entity');
    },
  },
  watch: {
    entityTypes(entityTypes, oldEntityTypes) {
      if (!isEqual(entityTypes, oldEntityTypes)) {
        this.query.page = 1;

        this.updateValue(isArray(this.value) ? [] : '');
        this.fetchEntities();
      }
    },
  },
  methods: {
    updateValue(value) {
      this.$emit('input', value);
      this.$emit('update:value', value);
      this.$emit('update:modelValue', value);
    },

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        type: this.entityTypes,
      };
    },

    async fetchEntities() {
      const { fetchContextEntitiesListWithoutStore } = useEntity();

      try {
        this.pending = true;

        const { data: entities, meta } = await fetchContextEntitiesListWithoutStore({
          params: this.getQuery(),
        });

        this.pageCount = meta.page_count;

        const currentEntities = this.getCurrentEntities();

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

    getCurrentEntities() {
      if (!this.returnObject || typeof this.value === 'string') {
        return pick(this.entitiesById, isArray(this.value) ? this.value : [this.value]);
      }

      return keyBy(this.value, '_id');
    },
  },
};
</script>
