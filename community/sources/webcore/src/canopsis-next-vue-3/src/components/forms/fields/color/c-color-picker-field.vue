<template>
  <v-input :hide-details="hideDetails">
    <div
      :class="[
        'd-flex flex-column align-stretch',
        justifyEnd ? 'justify-end' : 'justify-center',
        { 'w-100': fillWidth },
      ]"
    >
      <div v-if="splitted" class="d-flex align-center justify-space-between">
        <v-btn :disabled="disabled" class="mr-3" @click="shown = true">
          {{ displayLabel }}
        </v-btn>

        <div :style="style" class="pa-1 text-center">
          {{ colorValue }}
        </div>
      </div>

      <div
        :class="[
          'color-picker__button__wrapper d-flex align-center',
          justifyEnd ? 'justify-end' : 'justify-center',
          { 'w-100': fillWidth },
        ]"
      >
        <v-menu v-model="shown" :close-on-content-click="false" location="bottom">
          <template #activator="{ props: activatorProps }">
            <slot
              :props="activatorProps"
              :color="colorValue"
              :label="displayLabel"
              :style="style"
              :disabled="disabled"
              name="activator"
            >
              <c-btn-with-error
                v-bind="activatorProps"
                :error="buttonError"
                :size="size"
                :style="style"
                :disabled="disabled"
              >
                <v-icon v-if="icon">{{ icon }}</v-icon>
                <span v-else-if="!hideLabel">{{ displayLabel }}</span>
              </c-btn-with-error>
            </slot>
          </template>

          <v-card>
            <v-color-picker v-model="colorValue" :mode="type" show-swatches />
          </v-card>
        </v-menu>
      </div>
    </div>
  </v-input>
</template>

<script>
import { computed, ref, watch } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';

import { useField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    label: {
      type: String,
      required: false,
    },
    color: {
      type: String,
      default: '',
    },
    icon: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: 'hex',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    splitted: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'color',
    },
    required: {
      type: Boolean,
      default: false,
    },
    justifyEnd: {
      type: Boolean,
      default: false,
    },
    hideLabel: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    size: {
      type: String,
      default: undefined,
    },
    fillWidth: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:color'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const shown = ref(false);
    const fieldName = computed(() => props.name || props.label || 'color');
    const fieldRules = computed(() => ({
      required: props.required,
    }));
    const { value, errors, setValue, validate } = useField(fieldName, fieldRules, {
      initialValue: props.modelValue ?? props.color,
      validateOnMount: false,
      validateOnValueUpdate: false,
    });

    const colorValue = computed({
      get: () => props.modelValue ?? props.color ?? value.value,
      set: (value) => {
        setValue(value, false);
        emit('update:modelValue', value);
        emit('update:color', value);
      },
    });

    const displayLabel = computed(() => props.label ?? t('common.selectColor'));
    const buttonError = computed(() => errors.value[0] ?? '');

    const style = computed(() => ({
      backgroundColor: colorValue.value,
      color: getMostReadableTextColor(
        colorValue.value,
        props.icon ? { level: 'AA', size: 'large' } : undefined,
      ),
    }));

    watch(
      () => props.modelValue,
      (nextValue) => {
        if (nextValue !== undefined) {
          setValue(nextValue, false);
        }
      },
    );

    watch(
      () => props.color,
      (nextValue) => {
        if (props.modelValue === undefined) {
          setValue(nextValue, false);
        }
      },
    );

    return {
      shown,
      colorValue,
      displayLabel,
      style,
      buttonError,
      validate,
    };
  },
};
</script>
