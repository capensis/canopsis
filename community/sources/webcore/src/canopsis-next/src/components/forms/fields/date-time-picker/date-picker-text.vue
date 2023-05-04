<template lang="pug">
  v-text-field(
    :class="contentClass",
    :value="value | date(format)",
    :label="label",
    :error="error",
    :hide-details="hideDetails",
    disabled,
    readonly
  )
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
      type: [String, Object],
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
