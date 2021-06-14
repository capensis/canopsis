<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs12)
      v-textarea(
        v-validate="rule",
        v-on="listeners",
        :value="localValue",
        :label="label",
        :name="name",
        :rows="rows",
        :auto-grow="autoGrow",
        :box="box",
        :readonly="readonly",
        :outline="outline",
        :disabled="disabled",
        :error-messages="errors.collect(name)",
        data-vv-validate-on="none"
      )
        v-tooltip(v-if="helpText", slot="append", left)
          v-icon(slot="activator") help
          div(v-html="helpText")
    v-flex(v-if="!validateOnBlur && !readonly", xs12)
      v-btn.ml-0(
        :disabled="errors.has(name) || !wasChanged",
        color="primary",
        outline,
        @click="parse"
      ) {{ $t('common.parse') }}
      v-btn(
        :disabled="!wasChanged",
        color="grey darken-1",
        outline,
        @click="reset"
      ) {{ $t('common.reset') }}
</template>

<script>
import { get, isString } from 'lodash';
import { Validator } from 'vee-validate';

import { PAYLOAD_VARIABLE_REGEXP } from '@/constants';

import { convertPayloadToJson } from '@/helpers/payload-json';

import { isValidJson } from '@/plugins/validator/helpers/is-valid-json';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String],
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'json',
    },
    rows: {
      type: [Number, String],
      default: 5,
    },
    validateOn: {
      type: String,
      default: 'blur',
      validate: value => ['blur', 'button'].includes(value),
    },
    helpText: {
      type: String,
      default: '',
    },
    variables: {
      type: [Boolean, Array],
      default: false,
    },
    autoGrow: {
      type: Boolean,
      default: false,
    },
    box: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    outline: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const localValue = this.valueToLocalValue(this.value);

    return {
      localValue,

      sourceValue: localValue,
    };
  },
  computed: {
    validateOnBlur() {
      return this.validateOn === 'blur';
    },

    wasChanged() {
      return get(this.fields, [this.name, 'changed']);
    },

    listeners() {
      const listeners = {
        input: this.resetValidation,
      };

      if (this.validateOnBlur) {
        listeners.blur = this.parse;
      }

      return listeners;
    },

    rule() {
      return this.variables ? 'payload' : 'json';
    },
  },
  watch: {
    value(newValue) {
      const newLocalValue = this.valueToLocalValue(newValue);
      this.sourceValue = newLocalValue;

      if (newLocalValue !== this.localValue) {
        this.localValue = newLocalValue;
        this.$validator.reset({ name: this.name });
      }
    },
  },
  created() {
    this.$validator.extend('payload', {
      getMessage: () => this.$t('errors.JSONNotValid'),
      /**
       * Function for check json payload with variables is valid
       *
       * @param {string} json
       * @return {boolean}
       */
      validate: (json) => {
        try {
          const string = json.replace(new RegExp(PAYLOAD_VARIABLE_REGEXP), '""');

          return isValidJson(string);
        } catch (e) {
          return false;
        }
      },
    });
  },
  methods: {
    valueToLocalValue(value) {
      try {
        return this.variables
          ? convertPayloadToJson(value, 4)
          : this.$options.filters.json(value);
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });

        return '{}';
      }
    },

    parse() {
      try {
        const payload = this.variables
          ? convertPayloadToJson(this.localValue)
          : JSON.parse(this.localValue);

        this.$validator.reset({ name: this.name });

        if (this.value === payload) {
          this.localValue = this.valueToLocalValue(payload);
        }

        if (this.variables) {
          this.$emit('input', payload);
        } else {
          this.$emit('input', isString(this.value) ? this.localValue : payload);
        }
      } catch (err) {
        this.errors.add({
          field: this.name,
          msg: this.$t('errors.JSONNotValid'),
          rule: this.rule,
        });
      }
    },

    reset() {
      this.localValue = this.valueToLocalValue(this.value);

      this.$validator.reset({ name: this.name });
    },

    resetValidation(value) {
      this.localValue = value;

      if (value === this.sourceValue) {
        this.$validator.reset({ name: this.name });
        return;
      }

      if (this.errors.has(this.name)) {
        this.errors.remove(this.name);
      }

      if (!this.wasChanged) {
        this.$validator.flag(this.name, {
          touched: true,
          changed: true,
        });
      }
    },
  },
};
</script>
