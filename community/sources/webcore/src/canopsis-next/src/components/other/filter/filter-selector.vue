<template lang="pug">
  v-select(
    v-field="value",
    :items="preparedFilters",
    :label="label",
    :disabled="disabled",
    item-text="title",
    item-value="_id",
    clearable
  )
    template(#item="{ parent, item, tile }")
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
    preparedFilters() {
      const preparedFilters = [...this.filters];

      if (!this.lockedFilters.length) {
        return preparedFilters;
      }

      if (preparedFilters.length) {
        preparedFilters.push({ divider: true });
      }

      preparedFilters.push(...this.lockedFilters);

      return preparedFilters;
    },
  },
};
</script>
