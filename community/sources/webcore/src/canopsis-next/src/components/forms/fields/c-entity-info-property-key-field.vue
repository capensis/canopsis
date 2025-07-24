<template>
  <c-lazy-search-field
    :value="value"
    :label="label"
    :name="name"
    :required="required"
    :items="items"
    :loading="pending"
    item-text="value"
    item-value="value"
    @input="$emit('input', $event)"
    @search="fetchList"
  />
</template>

<script>
import { ref, onMounted } from 'vue';

import { useService } from '@/hooks/store/modules/service';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const items = ref([]);
    const pending = ref(false);

    const { fetchInfosKeysWithoutStore } = useService();

    const fetchList = async (search = '') => {
      try {
        pending.value = true;

        const { data } = await fetchInfosKeysWithoutStore({
          params: {
            search,
            limit: 50,
          },
        });

        items.value = data || [];
      } catch (err) {
        console.error(err);
        items.value = [];
      } finally {
        pending.value = false;
      }
    };

    onMounted(() => {
      fetchList();
    });

    return {
      items,
      pending,
      fetchList,
    };
  },
};
</script>
