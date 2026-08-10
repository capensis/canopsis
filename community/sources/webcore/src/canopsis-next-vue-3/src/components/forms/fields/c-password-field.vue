<script setup>
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';

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
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const shownPassword = ref(false);
const shownField = ref(false);

const fieldLabel = computed(() => props.label || t('common.password'));
const fieldPlaceholder = computed(() => props.placeholder || t('common.password'));

const appendInnerIcon = computed(() => {
  if (!props.visibility) {
    return undefined;
  }

  return shownPassword.value ? 'visibility' : 'visibility_off';
});

const inputType = computed(() => (shownPassword.value ? 'text' : 'password'));

const updateModel = value => emit('update:modelValue', value);
const toggleShownPassword = () => (shownPassword.value = !shownPassword.value);
const toggleShownField = () => {
  shownField.value = !shownField.value;

  if (!shownField.value) {
    updateModel('');
  }
};
</script>

<template>
  <c-layout v-if="replaceable" class="c-password-field--replaceable" column>
    <v-label class="c-password-field__label">
      <span v-if="required && fieldLabel" class="text-error"> * </span>
      {{ fieldLabel }}
    </v-label>

    <c-layout v-if="shownField" align="center" gap="3">
      <c-text-field
        v-bind="$attrs"
        :model-value="modelValue"
        :placeholder="fieldPlaceholder"
        :name="name"
        :type="inputType"
        :append-inner-icon="appendInnerIcon"
        :clearable="false"
        :required="required"
        @update:model-value="updateModel"
        @click:append-inner="toggleShownPassword"
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
    v-bind="$attrs"
    :model-value="modelValue"
    :label="fieldLabel"
    :name="name"
    :type="inputType"
    :append-inner-icon="appendInnerIcon"
    :clearable="false"
    :required="required"
    @update:model-value="updateModel"
    @click:append-inner="toggleShownPassword"
  />
</template>

<style lang="scss" scoped>
.c-password-field__label {
  transform: scale(0.75);
  transform-origin: left;
}
</style>
