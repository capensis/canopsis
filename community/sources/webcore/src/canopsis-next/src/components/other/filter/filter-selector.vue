<template lang="pug">
  v-select(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :disabled="disabled",
    :always-dirty="!!lockedFilter",
    item-text="title",
    item-value="_id",
    clearable
  )
    template(#selections="{ items }")
      v-tooltip(v-if="lockedFilter", top)
        template(#activator="{ on }")
          v-chip(v-on="on", small) {{ lockedFilter.title }}
        span {{ $t('settings.lockedFilter') }}
      span {{ items.map(item => item.title).join(', ') }}

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
export default {
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
    isLockedFilter(filter) {
      return filter._id === this.lockedValue;
    },
  },
};
</script>
