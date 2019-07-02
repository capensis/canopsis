<template lang="pug">
  v-toolbar.white(dense, flat)
    v-text-field(
      :value="value",
      :label="$t('common.search')",
      hide-details,
      single-line,
      @keyup.enter="submit",
      @input="$emit('input', $event)"
    )
    v-tooltip(bottom)
      v-btn(slot="activator", icon, @click="submit")
        v-icon search
      span {{ $t('common.searchSubmit') }}
    v-tooltip(bottom)
      v-btn(slot="activator", icon, @click="clear")
        v-icon clear
      span {{ $t('common.searchClear') }}
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
      this.$emit('submit');
    },
  },
};
</script>
