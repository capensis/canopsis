<template lang="pug">
  c-select-field.c-entity-field(
    v-field="value",
    v-validate="'required'",
    :search-input="query.search",
    :label="selectLabel",
    :loading="entitiesPending",
    :items="entities",
    :name="name",
    :item-text="itemText",
    :item-value="itemValue",
    :multiple="isMultiply",
    :deletable-chips="isMultiply",
    :small-chips="isMultiply",
    :error-messages="errors.collect(name)",
    :disabled="disabled",
    :menu-props="{ contentClass: 'c-entity-field__list' }",
    dense,
    autocomplete,
    @focus="onFocus",
    @blur="onBlur",
    @update:searchInput="debouncedUpdateSearch"
  )
    template(#item="{ item, tile }")
      v-list-tile.c-entity-field--tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-content {{ item[itemText] }}
        span.ml-4.grey--text(v-if="shownType") {{ item.type }}
    template(#append-item="")
      div.c-entity-field__append(ref="append")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { debounce, isEqual, keyBy, isArray } from 'lodash';

import { BASIC_ENTITY_TYPES } from '@/constants';

import { PAGINATION_LIMIT } from '@/config';

const { mapActions: entityMapActions } = createNamespacedHelpers('entity');

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String],
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
      type: String,
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
  },
  data() {
    return {
      entitiesById: {},
      entitiesPending: false,
      pageCount: Infinity,

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

    shownType() {
      return this.entityTypes.length !== 1;
    },

    isMultiply() {
      return isArray(this.value);
    },

    selectLabel() {
      if (this.label) {
        return this.label;
      }

      if (this.isMultiply) {
        return this.$tc('common.entity', this.value.length);
      }

      return this.$tc('common.entity');
    },
  },
  watch: {
    query: {
      deep: true,
      handler(newQuery, prevQuery) {
        if (this.isFocused && !isEqual(newQuery, prevQuery)) {
          this.fetchEntities();
        }
      },
    },
  },
  created() {
    this.debouncedUpdateSearch = debounce(this.updateSearch, 300);
  },
  mounted() {
    this.observer = new IntersectionObserver(this.intersectionHandler);

    this.observer.observe(this.$refs.append);

    if (this.value) {
      this.fetchEntities(this.value);
    }
  },
  beforeDestroy() {
    this.observer.unobserve(this.$refs.append);
  },
  methods: {
    ...entityMapActions({ fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore' }),

    intersectionHandler(entries) {
      const [entry] = entries;

      if (entry.isIntersecting && this.pageCount > this.query.page) {
        this.query = {
          ...this.query,
          page: this.query.page + 1,
        };
      }
    },

    updateSearch(value) {
      this.query = {
        page: 1,
        search: value,
      };
    },

    onFocus() {
      this.isFocused = true;

      if (!this.entities.length) {
        this.fetchEntities();
      }
    },

    onBlur() {
      this.isFocused = false;
    },

    getParams(ids) {
      const params = {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        filter: { type: { $in: this.entityTypes } },
      };

      if (ids) {
        params.filter._id = { $in: isArray(ids) ? ids : [ids] };
      }

      return params;
    },

    async fetchEntities(ids) {
      this.entitiesPending = true;

      const { data: entities, meta } = await this.fetchContextEntitiesListWithoutStore({
        params: this.getParams(ids),
      });

      this.pageCount = meta.page_count;

      this.entitiesById = {
        ...this.entitiesById,
        ...keyBy(entities, this.itemValue),
      };
      this.entitiesPending = false;
    },
  },
};
</script>

<style lang="scss">
.c-entity-field {
  &__list .v-list {
    position: relative;
  }

  &__append {
    position: absolute;
    pointer-events: none;
    right: 0;
    bottom: 0;
    left: 0;
    height: 300px;
  }
}
</style>
