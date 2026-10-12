<template>
  <c-layout v-if="replaceable" class="c-password-field--replaceable" column>
    <v-label class="c-password-field__label">
      <span v-if="required && fieldLabel" class="text-error"> * </span>
      {{ fieldLabel }}
    </v-label>

    <c-layout v-if="shownField" align="center" gap="3">
      <c-text-field
        ref="fieldElement"
        v-bind="$attrs"
        :model-value="displayedValue"
        :placeholder="fieldPlaceholder"
        :name="name"
        :type="inputType"
        :append-inner-icon="appendInnerIcon"
        :clearable="false"
        :required="required"
        :readonly="isValueMasked"
        @update:model-value="updateModel"
        @click:append-inner="toggleShownPassword"
        @focus="revealValue"
        @blur="maskValueIfEmpty"
      />

      <v-btn
        icon="delete"
        size="small"
        variant="text"
        @click="toggleShownField"
      />
    </c-layout>

    <div v-else>
      <v-btn color="primary" variant="outlined" @click="toggleShownField">
        {{ buttonLabel || `${t('common.replace')} ${fieldLabel}` }}
      </v-btn>
    </div>
  </c-layout>

  <c-text-field
    v-else
    ref="fieldElement"
    v-bind="$attrs"
    :model-value="displayedValue"
    :label="fieldLabel"
    :name="name"
    :type="inputType"
    :append-inner-icon="appendInnerIcon"
    :clearable="false"
    :required="required"
    :readonly="isValueMasked"
    @update:model-value="updateModel"
    @click:append-inner="toggleShownPassword"
    @focus="revealValue"
    @blur="maskValueIfEmpty"
  />
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

const MASKED_VALUE_PLACEHOLDER = '*'.repeat(30);

const props = defineProps({
  modelValue: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  placeholder: {
    type: String,
    default: '',
  },
  buttonLabel: {
    type: String,
    default: '',
  },
  name: {
    type: String,
    default: 'password',
  },
  required: {
    type: Boolean,
    default: false,
  },
  replaceable: {
    type: Boolean,
    default: false,
  },
  visibility: {
    type: Boolean,
    default: false,
  },
  masked: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const fieldElement = ref(null);

const shownPassword = ref(false);
const shownField = ref(false);
const isValueMasked = ref(props.masked && !props.modelValue);

const fieldLabel = computed(() => props.label || t('common.password'));
const fieldPlaceholder = computed(() => props.placeholder || t('common.password'));
const displayedValue = computed(() => (
  isValueMasked.value ? MASKED_VALUE_PLACEHOLDER : props.modelValue
));

const appendInnerIcon = computed(() => {
  if (!props.visibility) {
    return undefined;
  }

  return shownPassword.value ? 'visibility' : 'visibility_off';
});

const inputType = computed(() => (shownPassword.value ? 'text' : 'password'));

/**
 * Returns the root DOM element of the password text field.
 *
 * @returns {HTMLElement|undefined}
 */
const getFieldEl = () => fieldElement.value?.$el ?? fieldElement.value;

let blurTimeout;

/**
 * Restores the masked placeholder after blur when the value is still empty.
 */
const maskValueIfEmpty = () => {
  window.clearTimeout(blurTimeout);

  blurTimeout = window.setTimeout(() => {
    if (getFieldEl()?.contains(document.activeElement)) {
      return;
    }

    isValueMasked.value = props.masked && !props.modelValue;
  });
};

/**
 * Stops masking the field and shows the real empty value.
 */
const revealValue = () => {
  window.clearTimeout(blurTimeout);
  isValueMasked.value = false;
};

/**
 * Focuses the native input inside the password text field.
 */
const focusField = () => {
  nextTick(() => getFieldEl()?.querySelector?.('input')?.focus());
};

/**
 * Emits the new field value and ignores the masked placeholder.
 *
 * @param {string} value
 */
const updateModel = (value) => {
  revealValue();

  if (value === MASKED_VALUE_PLACEHOLDER) {
    return;
  }

  emit('update:modelValue', value);
};

/**
 * Toggles password visibility and reveals the real value when it was masked.
 */
const toggleShownPassword = () => {
  revealValue();
  shownPassword.value = !shownPassword.value;
  focusField();
};

/**
 * Shows or hides the replaceable password field.
 */
const toggleShownField = () => {
  shownField.value = !shownField.value;

  if (!shownField.value) {
    updateModel('');
  }
};

watch(() => props.masked, masked => isValueMasked.value = masked && !props.modelValue);

onBeforeUnmount(() => window.clearTimeout(blurTimeout));
</script>

<style lang="scss" scoped>
.c-password-field__label {
  transform: scale(0.75);
  transform-origin: left;
}
</style>
