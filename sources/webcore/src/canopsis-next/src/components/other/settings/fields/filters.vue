<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.filters')}}
    v-container
      v-select(
      :label="$t('settings.selectAFilter')",
      :items="filters",
      :value="value",
      @change="updateField",
      item-text="title",
      item-value="title",
      clearable
      )
</template>

<script>
export default {
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  methods: {
    updateField(value) {
      this.$emit('input', this.filters.find(v => v.title === value));
    },
  },
};
</script>
