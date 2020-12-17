<template lang="pug">
  div
    v-textarea(
      v-model="localValue",
      v-validate="'json'",
      :label="label",
      :name="name",
      :rows="rows",
      :error-messages="errorsMessages",
      data-vv-validate-on="none",
      @input="resetValidation"
    )
    v-btn.ml-0(
      :disabled="errors.has(name) || !wasTouched",
      color="primary",
      outline,
      @click="validate"
    ) {{ $t('common.parse') }}
</template>

<script>
import { get, isString } from 'lodash';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [String, Object],
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
      default: 10,
    },
  },
  data() {
    return {
      localValue: this.valueToLocalValue(this.value),
    };
  },
  computed: {
    wasTouched() {
      return get(this.fields, [this.name, 'touched']);
    },

    errorsMessages() {
      return this.errors.collect(this.name)
        .map(error => (error.rule === 'json' ? this.$t('errors.JSONNotValid') : error));
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

    validate() {
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

    resetValidation() {
      if (this.errors.has(this.name)) {
        this.errors.remove(this.name);
      }

      if (!this.wasTouched) {
        this.$validator.flag(this.name, {
          touched: true,
        });
      }
    },
  },
};
</script>
