<template>
  <c-select
    :model-value="modelValue"
    :items="availableRoles"
    :label="label || t('common.role')"
    :loading="pending"
    :name="name"
    :disabled="disabled"
    :multiple="multiple"
    :chips="chips"
    :menu-props="menuProps"
    :item-props="getItemProps"
    :required="required"
    item-title="name"
    item-value="_id"
    return-object
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';

import { API_ROUTES } from '@/config';
import { MAX_LIMIT } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import request from '@/services/request';

export default {
  props: {
    modelValue: {
      type: [Object, String, Array],
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'role',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    permission: {
      type: String,
      default: '',
    },
    isDisabledItems: {
      type: Function,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const menuProps = { offsetY: true };
    const items = ref([]);
    const pending = ref(false);

    const availableRoles = computed(() => {
      if (items.value.length) {
        return items.value;
      }

      return Array.isArray(props.modelValue)
        ? props.modelValue
        : [props.modelValue].filter(Boolean);
    });

    const fetchList = async () => {
      try {
        pending.value = true;

        const params = { limit: MAX_LIMIT };

        if (props.permission) {
          params.permission = props.permission;
        }

        const { data = [] } = await request.get(API_ROUTES.roles.list, { params });

        items.value = data;
      } finally {
        pending.value = false;
      }
    };

    const getItemProps = item => ({
      disabled: props.isDisabledItems?.(item),
    });

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    watch(
      () => props.disabled,
      (disabled) => {
        if (!disabled && !items.value.length) {
          fetchList();
        }
      },
    );

    onMounted(() => {
      if (!props.disabled) {
        fetchList();
      }
    });

    return {
      t,
      menuProps,
      pending,
      availableRoles,
      getItemProps,
      updateModel,
    };
  },
};
</script>
