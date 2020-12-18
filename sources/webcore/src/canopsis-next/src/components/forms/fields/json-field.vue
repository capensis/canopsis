<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs12)
      v-textarea(
        v-validate="'json'",
        :value="localValue",
        :label="label",
        :name="name",
        :rows="rows",
        :error-messages="errors.collect(name)",
        data-vv-validate-on="none",
        v-on="listeners"
      )
        v-tooltip(slot="append", v-if="helpText", left)
          v-icon(slot="activator") help
          div(v-html="helpText")
    v-flex(v-if="!validateOnBlur", xs12)
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

export default {
  inject: ['$validator'],
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
      required: true,
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
  },
  data() {
    return {
      localValue: this.valueToLocalValue(this.value),
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
  },
  watch: {
    value(newValue) {
      const newLocalValue = this.valueToLocalValue(newValue);

      if (newLocalValue !== this.localValue) {
        this.localValue = newLocalValue;
        this.$validator.reset({ name: this.name });
      }
    },
  },
  methods: {
    valueToLocalValue(value) {
      try {
        return this.$options.filters.json(value);
      } catch (err) {
        this.$popups.error({ text: this.$t('errors.default') });

        return '{}';
      }
    },

    parse() {
      try {
        const newValue = JSON.parse(this.localValue);

        this.$validator.reset({ name: this.name });

        this.$emit('input', isString(this.value) ? this.localValue : newValue);
      } catch (err) {
        this.errors.add({
          field: this.name,
          msg: this.$t('errors.JSONNotValid'),
          rule: 'json',
        });
      }
    },

    reset() {
      this.localValue = this.valueToLocalValue(this.value);

      this.$validator.reset({ name: this.name });
    },

    resetValidation(value) {
      if (value === this.localValue) {
        return;
      }

      this.localValue = value;

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
