<template>
  <c-autocomplete
    :model-value="value"
    :items="items"
    :label="label || t('common.type')"
    :loading="pending"
    :multiple="multiple"
    :chips="chips"
    :closable-chips="false"
    item-title="name"
    item-value="_id"
    return-object
    disabled
  />
</template>

<script>
import { onMounted, ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

export default {
  props: {
    value: {
      type: [Object, String, Array],
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    const items = ref([]);
    const pending = ref(false);

    const fetchList = async () => {
      try {
        pending.value = true;

        const { data = [] } = await fetchPbehaviorTypesFieldList({
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
