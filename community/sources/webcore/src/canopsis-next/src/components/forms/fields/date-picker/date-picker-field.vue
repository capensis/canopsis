<template lang="pug">
  v-menu(
    v-model="opened",
    :close-on-content-click="false",
    :disabled="disabled",
    content-class="date-picker",
    transition="slide-y-transition",
    max-width="290px",
    right,
    lazy
  )
    div(slot="activator")
      v-text-field(
        :value="value | date(format, true)",
        :label="label",
        :error="error",
        :error-messages="errorMessages",
        :name="name",
        :disabled="disabled",
        :hide-details="hideDetails",
        :append-icon="clearable ? 'close' : ''",
        :readonly="!disabled",
        @click:append="clear"
      )
    v-date-picker.date-picker(
      :value="value",
      :opened="opened",
      :color="color",
      @input="input",
      @change="change"
    )
</template>

<script>

import formBaseMixin from '@/mixins/form/base';

/**
 * Date picker field component
 */
export default {
  $_veeValidate: {
    value() {
      return this.value;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  mixins: [formBaseMixin],
  props: {
    clearable: {
      type: Boolean,
      default: false,
    },
    value: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: null,
    },
    color: {
      type: String,
      default: 'primary',
    },
    format: {
      type: String,
      default: 'datePicker',
    },
    error: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    errorMessages() {
      if (this.$validator && this.errors && this.name) {
        return this.errors.collect(this.name);
      }

      return [];
    },
  },
  methods: {
    clear() {
      this.updateModel(null);
    },
    input(value) {
      this.updateModel(value);
    },
    change(value) {
      this.$emit('change', value);

      this.opened = false;
    },
  },
};
</script>

<style lang="scss">
  .date-picker {
    .v-picker__body,
    .v-time-picker-clock__item,
    .v-time-picker-clock__item span,
    .v-time-picker-clock__hand {
      z-index: inherit;
    }
  }
</style>
