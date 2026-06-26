<template>
  <c-select
    :model-value="modelValue"
    :items="languages"
    :disabled="disabled"
    :label="label"
    :name="name"
    :menu-props="menuProps"
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed } from 'vue';

import { LOCALES } from '@/config';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
    value: {
      type: String,
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'language',
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const menuProps = { offsetY: true };
    const languages = computed(() => Object.values(LOCALES));

    const updateModel = value => {
      emit('update:modelValue', value);
    };

    return {
      menuProps,
      languages,
      updateModel,
    };
  },
};
</script>
