<template lang="pug">
  v-menu(
    v-model="opened",
    :close-on-content-click="false",
    content-class="time-picker",
    transition="slide-y-transition",
    max-width="290px",
    right,
    lazy
  )
    div(slot="activator")
      v-text-field(
        ref="textField",
        :value="value",
        :label="label",
        :error="error",
        :error-messages="errorMessages",
        :name="name",
        :hide-details="hideDetails",
        :append-icon="clearable ? 'close' : ''",
        readonly,
        @click:append="clear"
      )
    v-time-picker.time-picker(
      :value="value",
      :opened="opened",
      :color="color",
      :format="format",
      @input="input",
      @change="change"
    )
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

/**
 * Time picker field component
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
      default: '24hr',
    },
    error: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
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
      this.$refs.textField.blur();
    },
  },
};
</script>

<style lang="scss">
  .time-picker {
    .v-picker__body,
    .v-time-picker-clock__item,
    .v-time-picker-clock__item span,
    .v-time-picker-clock__hand {
      z-index: inherit;
    }
  }
</style>
