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
    :menu-props="{ contentClass: 'c-entity-field__list' }",
    dense,
    autocomplete,
    @focus="onFocus",
    @update:searchInput="updateSearch"
  )
    template(#item="{ item, tile }")
      v-list-tile.c-entity-field--tile(ref="a", v-bind="tile.props", v-on="tile.on")
        v-list-tile-content {{ item[itemText] }}
        span.ml-4.grey--text {{ item.type }}
    template(#append-item="")
      div.c-entity-field__append(ref="append")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { debounce, isEqual } from 'lodash';
import { PAGINATION_LIMIT } from '@/config';
import { BASIC_ENTITY_TYPES } from '@/constants';

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
  },
  data() {
    return {
      entities: [],
      entitiesPending: false,
      pageCount: Infinity,

      query: {
        page: 1,
        search: null,
      },
    };
  },
  computed: {
    isMultiply() {
      return Array.isArray(this.value);
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
        if (!isEqual(newQuery, prevQuery)) {
          this.debouncedFetchEntities();
        }
      },
    },
  },
  created() {
    this.debouncedFetchEntities = debounce(this.fetchEntities, 300);
  },
  mounted() {
    this.observer = new IntersectionObserver(this.intersectionHandler);

    this.observer.observe(this.$refs.append);
  },
  beforeDestroy() {
    this.observer.unobserve(this.$refs.append);
  },
  methods: {
    ...entityMapActions({ fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore' }),

    intersectionHandler(entries) {
      const [entry] = entries;

      if (entry.isIntersecting && this.pageCount >= this.query.page) {
        this.query.page += 1;
      }
    },

    updateSearch(value) {
      this.query.page = 1;
      this.query.search = value;
    },

    onFocus() {
      if (!this.entities.length) {
        this.fetchEntities();
      }
    },

    async fetchEntities() {
      this.entitiesPending = true;

      const { data: entities, meta } = await this.fetchContextEntitiesListWithoutStore({
        params: {
          limit: this.limit,
          page: this.query.page,
          search: this.query.search,
          filter: { type: { $in: this.entityTypes } },
        },
      });

      this.pageCount = meta.page_count;

      this.entities.push(...entities);
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
