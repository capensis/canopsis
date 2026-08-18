<template>
  <v-row class="c-json-field" density="compact">
    <v-col cols="12">
      <c-payload-textarea-field
        v-if="variables?.length"
        v-bind="$attrs"
        :value="localValue"
        :label="label"
        :name="name"
        :rows="rows"
        :auto-grow="autoGrow"
        :readonly="readonly"
        :disabled="disabled"
        :variables="variables"
        @update:value="resetValidation"
        @blur="handleBlur"
      >
        <template v-if="helpText || $slots.append" #append>
          <c-help-icon
            v-if="helpText"
            :text="helpText"
            icon="help"
            left
          />
        </template>
      </c-payload-textarea-field>
      <c-textarea
        v-else
        v-bind="$attrs"
        :model-value="localValue"
        :label="label"
        :name="name"
        :rows="rows"
        :auto-grow="autoGrow"
        :readonly="readonly"
        :disabled="disabled"
        @update:model-value="resetValidation"
        @blur="handleBlur"
      >
        <template v-if="helpText" #append-inner>
          <c-help-icon :text="helpText" icon="help" left />
        </template>
      </c-textarea>
    </v-col>
    <v-col v-if="!validateOnBlur && !readonly" cols="12">
      <c-layout>
        <v-btn
          :disabled="!!errorMessage || !wasChanged"
          color="primary"
          variant="outlined"
          @click="parse"
        >
          {{ $t('common.parse') }}
        </v-btn>
        <v-btn
          :disabled="!wasChanged"
          class="ml-2"
          color="grey-darken-1"
          variant="outlined"
          @click="reset"
        >
          {{ $t('common.reset') }}
        </v-btn>
        <slot :json-string="localValue" name="append" />
      </c-layout>
    </v-col>
  </v-row>
</template>

<script>
import { isString } from 'lodash';
import { computed } from 'vue';

import { stringifyJson } from '@/helpers/json';
import { convertPayloadToJson } from '@/helpers/payload-json';

import { useField } from '@/hooks/form';

export default {
  inheritAttrs: false,
  emits: ['update:value'],
  props: {
    value: {
      type: [Object, Array, String],
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
      validator: value => ['blur', 'button'].includes(value),
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
  setup(props) {
    const fieldName = computed(() => props.name);
    const { errorMessage, setErrors } = useField(fieldName, undefined, {
      validateOnValueUpdate: false,
    });

    const setErrorMessage = message => setErrors(message ? [message] : []);

    return {
      errorMessage,
      setErrorMessage,
    };
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
      return this.valueToLocalValue(this.value) !== this.localValue;
    },
  },
  watch: {
    value(newValue) {
      let newLocalValue = '{}';
      let hasError = false;

      try {
        newLocalValue = this.valueToLocalValue(newValue, true);
      } catch {
        newLocalValue = newValue;
        hasError = true;
      } finally {
        this.sourceValue = newLocalValue;

        if (newLocalValue !== this.localValue) {
          this.localValue = newLocalValue;
        }

        this.setErrorMessage(hasError ? this.$t('errors.JSONNotValid') : undefined);
      }
    },
  },
  methods: {
    valueToLocalValue(value, throwError = false) {
      try {
        return this.variables ? convertPayloadToJson(value) : stringifyJson(value);
      } catch (err) {
        console.error(err);

        if (throwError) {
          throw err;
        }

        return '{}';
      }
    },

    handleBlur() {
      if (this.validateOnBlur) {
        this.parse();
      }
    },

    parse() {
      try {
        const payload = this.variables
          ? convertPayloadToJson(this.localValue)
          : JSON.parse(this.localValue);

        this.setErrorMessage(undefined);

        if (this.value === payload) {
          this.localValue = this.valueToLocalValue(payload);
        }

        this.$emit(
          'update:value',
          this.variables || !isString(this.value) ? payload : this.localValue,
        );
      } catch {
        this.setErrorMessage(this.$t('errors.JSONNotValid'));
      }
    },

    reset() {
      this.localValue = this.valueToLocalValue(this.value);
      this.setErrorMessage(undefined);
    },

    resetValidation(value) {
      this.localValue = value;

      if (this.errorMessage) {
        this.setErrorMessage(undefined);
      }

      if (value === this.sourceValue) {
        this.setErrorMessage(undefined);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.c-json-field {
  width: 100%;
}
</style>
