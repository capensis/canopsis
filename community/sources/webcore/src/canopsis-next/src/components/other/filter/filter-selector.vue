<template lang="pug">
  v-select(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :disabled="disabled",
    :always-dirty="!!lockedFilter",
    :item-text="itemText",
    :item-value="itemValue"
  )
    template(#selections="{ items }")
      v-tooltip(v-if="lockedFilter", top)
        template(#activator="{ on }")
          v-chip(v-on="on", small)
            span {{ getItemText(lockedFilter) }}
            v-icon.ml-2(small) lock
        span {{ $t('settings.lockedFilter') }}
      v-chip(v-if="items.length", small, close, @input="cancelFilter") {{ getItemText(items[0]) }}

    template(#item="{ parent, item, tile }")
      v-list-tile(v-bind="tile.props", v-on="tile.on", :disabled="item.active")
        v-list-tile-content
          v-list-tile-title
            span {{ item.title }}
            v-icon.ml-2(
              v-if="!hideIcon",
              :color="tile.props.value ? parent.color : ''",
              small
            ) {{ item.is_private ? 'person' : 'lock' }}
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  props: {
    value: {
      type: String,
      default: () => null,
    },
    lockedValue: {
      type: String,
      required: false,
    },
    filters: {
      type: Array,
      default: () => [],
    },
    lockedFilters: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'title',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
    hideIcon: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    lockedFilter() {
      return this.lockedFilters.find(this.isLockedFilter);
    },

    preparedFilters() {
      const preparedFilters = [...this.filters];

      if (!this.lockedFilters.length) {
        return preparedFilters;
      }

      if (preparedFilters.length) {
        preparedFilters.push({ divider: true });
      }

      preparedFilters.push(
        ...this.lockedFilters.map(filter => ({
          ...filter,
          active: this.isLockedFilter(filter),
        })),
      );

      return preparedFilters;
    },
  },
  methods: {
    getItemText(item) {
      return item[this.itemText];
    },

    getItemValue(item) {
      return item[this.itemValue];
    },

    isLockedFilter(filter) {
      return this.getItemValue(filter) === this.lockedValue;
    },

    cancelFilter() {
      this.updateModel('');
    },
  },
};
</script>
