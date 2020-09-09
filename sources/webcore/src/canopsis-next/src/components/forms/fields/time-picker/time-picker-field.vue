<template lang="pug">
  v-menu(
    ref="menu",
    v-model="opened",
    content-class="date-time-picker",
    transition="slide-y-transition",
    max-width="290px",
    :close-on-content-click="false",
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
    v-time-picker(
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
 * Date picker field component
 *
 * @warning If you want to use validation on the field you shouldn't use `v-field`
 *
 * @prop {Date} [value=null] - Date value
 * @prop {Boolean} [clearable=false] - if it is true input field will be have cross button with clear event on click
 * @prop {string} [label=''] - Label of the input field
 * @prop {string} [name=null] - Name property in the validation object
 * @prop {Boolean} [roundHours=false] - Deny to change minutes it will be only 0
 *
 * @event value#input
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
