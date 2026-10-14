<template>
  <v-row class="c-json-field" density="compact">
    <v-col cols="12">
      <c-payload-textarea-field
        v-if="hasVariables"
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
        <template v-if="helpText || (validateOnBlur && $slots.append)" #append>
          <c-help-icon
            v-if="helpText"
            :text="helpText"
            icon="help"
            left
          />
          <slot
            v-if="validateOnBlur"
            :json-string="localValue"
            :error-message="errorMessage"
            name="append"
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
      <c-layout gap="2">
        <v-btn
          :disabled="!!errorMessage || !wasChanged"
          color="primary"
          variant="outlined"
          @click="parse"
        >
          {{ t('common.parse') }}
        </v-btn>
        <v-btn
          :disabled="!wasChanged"
          color="grey-darken-1"
          variant="outlined"
          @click="reset"
        >
          {{ t('common.reset') }}
        </v-btn>
        <slot
          :json-string="localValue"
          :error-message="errorMessage"
          name="append"
        />
      </c-layout>
    </v-col>
  </v-row>
</template>

<script setup>
import { isString } from 'lodash';
import { computed, ref, toRef, watch } from 'vue';

import { stringifyJson } from '@/helpers/json';
import { convertPayloadToJson } from '@/helpers/payload-json';

import { useField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
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
});

const emit = defineEmits(['update:value']);

const { t } = useI18n();
const { errorMessage, setErrors } = useField(toRef(props, 'name'), undefined, {
  validateOnValueUpdate: false,
});

const setErrorMessage = message => setErrors(message ? [message] : []);

const usesPayload = computed(() => Boolean(props.variables));
const hasVariables = computed(() => Boolean(props.variables?.length));
const validateOnBlur = computed(() => props.validateOn === 'blur');

const valueToLocalValue = (value, throwError = false) => {
  try {
    return usesPayload.value ? convertPayloadToJson(value) : stringifyJson(value);
  } catch (err) {
    console.error(err);

    if (throwError) {
      throw err;
    }

    return '{}';
  }
};

const localValue = ref(valueToLocalValue(props.value));
const sourceValue = ref(localValue.value);

const wasChanged = computed(() => valueToLocalValue(props.value) !== localValue.value);

watch(
  () => props.value,
  (newValue) => {
    let newLocalValue = '{}';
    let hasError = false;

    try {
      newLocalValue = valueToLocalValue(newValue, true);
    } catch {
      newLocalValue = newValue;
      hasError = true;
    } finally {
      sourceValue.value = newLocalValue;

      if (newLocalValue !== localValue.value) {
        localValue.value = newLocalValue;
      }

      setErrorMessage(hasError ? t('errors.JSONNotValid') : undefined);
    }
  },
);

const parse = () => {
  try {
    const payload = usesPayload.value
      ? convertPayloadToJson(localValue.value)
      : JSON.parse(localValue.value);

    setErrorMessage(undefined);

    if (props.value === payload) {
      localValue.value = valueToLocalValue(payload);
    }

    emit(
      'update:value',
      usesPayload.value || !isString(props.value) ? payload : localValue.value,
    );
  } catch {
    setErrorMessage(t('errors.JSONNotValid'));
  }
};

const handleBlur = () => {
  if (validateOnBlur.value) {
    parse();
  }
};

const reset = () => {
  localValue.value = valueToLocalValue(props.value);
  setErrorMessage(undefined);
};

const resetValidation = (value) => {
  localValue.value = value;

  if (errorMessage.value) {
    setErrorMessage(undefined);
  }

  if (value === sourceValue.value) {
    setErrorMessage(undefined);
  }
};
</script>

<style lang="scss" scoped>
.c-json-field {
  width: 100%;
}
</style>
