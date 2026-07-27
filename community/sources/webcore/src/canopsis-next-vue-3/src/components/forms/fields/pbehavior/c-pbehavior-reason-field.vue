<template>
  <c-autocomplete
    :model-value="modelValue"
    :items="items"
    :label="label || t('common.reason')"
    :loading="pending"
    :name="name"
    :required="required"
    :disabled="disabled"
    :return-object="returnObject"
    :clearable="clearable"
    item-title="name"
    item-value="_id"
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { onMounted, ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorReason } from '@/hooks/store/modules/pbehavior-reason';

export default {
  props: {
    modelValue: {
      type: [String, Object],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'reason',
    },
    required: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();
    const { fetchPbehaviorReasonsListWithoutStore } = usePbehaviorReason();

    const items = ref([]);
    const pending = ref(false);

    const fetchList = async () => {
      try {
        pending.value = true;

        const { data = [] } = await fetchPbehaviorReasonsListWithoutStore({
          params: {
            paginate: false,
            with_hidden: true,
          },
        });

        items.value = data;
      } finally {
        pending.value = false;
      }
    };

    onMounted(fetchList);

    return {
      t,
      items,
      pending,
    };
  },
};
</script>
