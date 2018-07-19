<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.defaultSortColumn') }}
    v-container
      v-text-field(
      :value="value.property",
      @input="inputProperty",
      :placeholder="$t('settings.columnName')"
      )
      v-select(
      :value="value.direction",
      @input="inputDirection",
      :items="directions"
      )
</template>

<script>

/**
* Component to select the default column to sort on settings
*
* @prop {Object} [value] - Object containing the default sort column's name and the sort direction
*
* @event value#input
*/
export default {
  props: {
    value: {
      type: Object,
      default: () => ({
        property: '',
        direction: 'ASC',
      }),
    },
  },
  data() {
    return {
      directions: ['ASC', 'DESC'],
    };
  },
  methods: {
    inputProperty(property) {
      this.$emit('input', { ...this.value, property });
    },
    inputDirection(direction) {
      this.$emit('input', { ...this.value, direction });
    },
  },
};
</script>
