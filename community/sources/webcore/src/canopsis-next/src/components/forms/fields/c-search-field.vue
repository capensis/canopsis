<template lang="pug">
  v-toolbar.white(dense, flat)
    v-text-field(
      :value="localValue",
      :label="$t('common.search')",
      data-test="searchingTextField",
      hide-details,
      single-line,
      @keydown.enter.prevent="submit",
      @input="input"
    )
    v-tooltip(bottom)
      v-btn(slot="activator", data-test="submitSearchButton", icon, @click="submit")
        v-icon search
      span {{ $t('common.search') }}
    v-tooltip(bottom)
      v-btn(slot="activator", data-test="clearSearchButton", icon, @click="clear")
        v-icon clear
      span {{ $t('common.clearSearch') }}
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
  data() {
    return {
      localValue: this.value,
    };
  },
  watch: {
    value(newValue) {
      if (newValue !== this.localValue) {
        this.localValue = newValue;
      }
    },
  },
  methods: {
    input(value) {
      this.localValue = value;

      this.$emit('input', value);
    },
    clear() {
      this.localValue = '';

      this.$emit('input', '');
      this.$emit('clear');
    },
    submit() {
      this.$emit('submit', this.localValue);
    },
  },
};
</script>
