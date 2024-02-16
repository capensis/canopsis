<template>
  <v-menu
    v-model="opened"
    :close-on-content-click="false"
    :disabled="disabled"
    content-class="date-picker"
    transition="slide-y-transition"
    max-width="290px"
    right
  >
    <template #activator="{ on }">
      <v-text-field
        v-on="on"
        :class="contentClass"
        :value="value | date(format)"
        :label="label"
        :placeholder="placeholder"
        :error="error"
        :error-messages="errorMessages"
        :name="name"
        :disabled="disabled"
        :hide-details="hideDetails"
        :append-icon="clearable ? 'close' : ''"
        :readonly="!disabled"
        @click:append="clear"
      >
        <template #append="">
          <slot name="append" />
        </template>
      </v-text-field>
    </template>
    <v-date-picker
      class="date-picker"
      v-field="value"
      :opened="opened"
      :color="color"
      :min="min"
      :max="max"
      :allowed-dates="allowedDates"
      @change="change"
    />
  </v-menu>
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

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
    value: {
      type: [String, Number],
      default: null,
    },
    label: {
      type: String,
      default: '',
    },
    placeholder: {
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
      default: 'short',
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
    min: {
      type: [String, Number],
      required: false,
    },
    max: {
      type: [String, Number],
      required: false,
    },
    allowedDates: {
      type: Function,
      required: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    contentClass: {
      type: [String, Object, Array],
      required: false,
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
