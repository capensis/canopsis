<template>
  <c-select
    :model-value="selectedType"
    :items="defaultTypes"
    :label="t('modals.createPbehaviorType.fields.type')"
    :disabled="disabled"
    :loading="pending"
    item-title="text"
    item-value="type"
    name="type"
    required
    return-object
    @update:model-value="updateModel"
  />
</template>

<script>
import { computed, onMounted, ref } from 'vue';
import { find } from 'lodash';

import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

export default {
  props: {
    modelValue: {
      type: String,
      default: PBEHAVIOR_TYPE_TYPES.active,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:color', 'sync:color'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchPbehaviorTypesListWithoutStore } = usePbehaviorType();

    const pending = ref(false);
    const defaultTypes = ref([]);

    const selectedType = computed(
      () => find(defaultTypes.value, { type: props.modelValue }) ?? props.modelValue,
    );

    const fetchDefaultTypes = async () => {
      try {
        pending.value = true;

        const { data = [] } = await fetchPbehaviorTypesListWithoutStore({
          params: { default: true },
        });

        defaultTypes.value = data.map(type => ({
          ...type,
          text: t(`modals.createPbehaviorType.canonicalTypes.${type.type}`),
        }));

        const activeType = find(data, { type: props.modelValue });

        if (activeType) {
          emit('sync:color', activeType.color);
        }
      } finally {
        pending.value = false;
      }
    };

    const updateModel = (value = {}) => {
      emit('update:modelValue', value.type);
      emit('update:color', value.color);
    };

    onMounted(fetchDefaultTypes);

    return {
      t,
      pending,
      defaultTypes,
      selectedType,
      updateModel,
    };
  },
};
</script>
