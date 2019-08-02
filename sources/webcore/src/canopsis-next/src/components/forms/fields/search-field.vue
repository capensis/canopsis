<template lang="pug">
  v-toolbar.white(dense, flat)
    v-text-field(
    :value="value",
    :label="$t('common.search')",
    data-test="searchingTextField"
    hide-details,
    single-line,
    @keyup.enter="submit",
    @input="$emit('input', $event)"
    )
    v-btn(data-test="submitSearchButton", icon, @click="submit")
      v-icon search
    v-btn(data-test="clearSearchButton", icon, @click="clear")
      v-icon clear
    slot
</template>

<script>
/**
 * Search component
 */
export default {
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  methods: {
    clear() {
      this.$emit('input', '');
      this.$emit('clear');
    },
    submit() {
      this.$emit('submit', this.value);
    },
  },
};
</script>
