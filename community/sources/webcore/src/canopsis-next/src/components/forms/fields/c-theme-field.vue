<template>
  <c-select
    :model-value="modelValue"
    :items="themes"
    :label="label || t('common.theme')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :hide-details="hideDetails"
    :clearable="clearable"
    :menu-props="menuProps"
    :required="required"
    item-title="name"
    item-value="_id"
    @update:model-value="updateModel"
  />
</template>

<script>
import { onMounted, ref } from 'vue';

import { API_ROUTES } from '@/config';
import { MAX_LIMIT } from '@/constants';

import request from '@/services/request';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: [Object, String],
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'theme',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const menuProps = { offsetY: true };
    const themes = ref([]);
    const pending = ref(false);

    const fetchList = async () => {
      try {
        pending.value = true;

        const { data = [] } = await request.get(API_ROUTES.themes.list, {
          params: { limit: MAX_LIMIT },
        });

        themes.value = data;
      } finally {
        pending.value = false;
      }
    };

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    onMounted(fetchList);

    return {
      t,
      menuProps,
      themes,
      pending,
      updateModel,
    };
  },
};
</script>
