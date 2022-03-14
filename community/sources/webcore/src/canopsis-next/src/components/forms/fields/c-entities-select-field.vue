<template lang="pug">
  v-autocomplete.c-entities-select-field(
    v-field="value",
    v-validate="'required'",
    :search-input.sync="searchInput",
    :label="selectLabel",
    :loading="loading",
    :items="items",
    :name="name",
    :item-text="itemText",
    :item-value="itemValue",
    :multiple="isMultiply",
    :deletable-chips="isMultiply",
    :small-chips="isMultiply",
    :error-messages="errors.collect(name)"
  )
    template(#item="{ item, tile }")
      v-list-tile.c-entities-select-field--tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-content {{ item | get(itemText) }}
</template>

<script>
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
    search: {
      type: String,
      default: null,
    },
    items: {
      type: Array,
      default: () => [],
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
      default: 'name',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      searchInput: this.search,
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
    searchInput() {
      this.$emit('update:search', this.searchInput);
    },
  },
};
</script>

<style scoped lang="scss">
.c-entities-select-field {
  &--tile {
    & /deep/ .v-list__tile {
      height: 36px;
    }
  }
}
</style>
