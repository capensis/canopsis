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
    :return-object="false",
    :menu-props="{ contentClass: 'c-entity-field__list' }",
    dense,
    combobox,
    @focus="onFocus",
    @blur="onBlur",
    @update:searchInput="debouncedUpdateSearch"
  )
    template(#item="{ item, tile }")
      v-list-tile.c-entity-field--tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-content {{ getItemText(item) }}
        span.ml-4.grey--text {{ item.type }}
    template(#append-item="")
      div.c-entity-field__append(ref="append")
    template(v-if="isMultiply", #selection="{ item, index }")
      v-chip.c-entity-field__chip(small, close, @input="removeItemFromArray(index)")
        span.ellipsis {{ getItemText(item) }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { debounce, isEqual, keyBy, isArray, isString } from 'lodash';

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
      isFocused: false,
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
  },
  beforeDestroy() {
    this.observer.unobserve(this.$refs.append);
  },
  methods: {
    ...entityMapActions({ fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore' }),

    getItemText(item) {
      return isString(item) ? item : item[this.itemText];
    },

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
      this.pageCount = Infinity;
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

    getQuery() {
      return {
        limit: this.limit,
        page: this.query.page,
        search: this.query.search,
        type: this.entityTypes,
      };
    },

    async fetchEntities() {
      this.entitiesPending = true;

      const { data: entities, meta } = await this.fetchContextEntitiesListWithoutStore({
        params: this.getQuery(),
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

  .v-select__selections {
    max-width: calc(100% - 24px);
  }

  &__chip {
    max-width: 100%;

    .v-chip__content {
      max-width: 100%;
    }
  }
}
</style>
