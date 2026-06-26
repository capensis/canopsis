<template>
  <c-layout
    :justify="justifyEnd ? 'end' : 'center'"
    column
  >
    <c-layout
      v-if="splitted"
      align="center"
      justify="space-between"
    >
      <v-btn
        :disabled="disabled"
        class="mr-3"
        @click="shown = true"
      >
        {{ displayLabel }}
      </v-btn>

      <div
        :style="style"
        class="pa-1 text-center"
      >
        {{ colorValue }}
      </div>
    </c-layout>

    <c-layout
      :justify="justifyEnd ? 'end' : 'center'"
      class="color-picker__button__wrapper"
      align="center"
    >
      <v-menu
        v-model="shown"
        :close-on-content-click="false"
        location="bottom"
      >
        <template #activator="{ props: activatorProps }">
          <v-btn
            v-bind="activatorProps"
            :style="style"
            :disabled="disabled"
          >
            {{ displayLabel }}
          </v-btn>
        </template>

        <v-card>
          <v-color-picker
            v-model="colorValue"
            :mode="type"
            hide-inputs
          />
        </v-card>
      </v-menu>
    </c-layout>

    <v-messages
      v-if="errorMessages.length"
      :messages="errorMessages"
      color="error"
    />
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { getMostReadableTextColor } from '@/helpers/color';

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
  },

  emits: ['update:modelValue', 'update:color'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const shown = ref(false);

    const colorValue = computed({
      get: () => props.modelValue ?? props.color,
      set: value => {
        emit('update:modelValue', value);
        emit('update:color', value);
      },
    });

    const displayLabel = computed(() => props.label ?? t('common.selectColor'));

    const style = computed(() => ({
      backgroundColor: colorValue.value,
      color: getMostReadableTextColor(colorValue.value),
    }));

    const errorMessages = computed(() => (
      props.required && !colorValue.value ? [`${props.name} is required`] : []
    ));

    return {
      shown,
      colorValue,
      displayLabel,
      style,
      errorMessages,
    };
  },
};
</script>
