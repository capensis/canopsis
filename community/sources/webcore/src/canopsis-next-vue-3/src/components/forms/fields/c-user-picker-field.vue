<template>
  <c-autocomplete
    :model-value="value"
    v-bind="$attrs"
    :items="items"
    :label="label"
    :loading="pending"
    :name="name"
    :required="required"
    :return-object="returnObject"
    :item-title="itemText"
    :item-value="itemValue"
    @update:model-value="updateValue"
  />
</template>

<script>
import { ref, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';
import { usePendingHandler } from '@/hooks/query/pending';
import { useUser } from '@/hooks/store/modules/user';

export default {
  inheritAttrs: false,
  emits: ['update:value', 'update:modelValue'],
  props: {
    value: {
      type: [Object, Array, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'user',
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    permission: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'display_name',
    },
    itemValue: {
      type: String,
      default: '_id',
    },
  },
  setup(props, { emit }) {
    const items = ref([]);

    const { fetchUsersListWithoutStore } = useUser();

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const params = { limit: MAX_LIMIT };

      if (props.permission) {
        params.permission = props.permission;
      }

      const { data } = await fetchUsersListWithoutStore({ params });

      items.value = data;
    });

    const updateValue = (value) => {
      emit('update:value', value);
      emit('update:modelValue', value);
    };

    onMounted(fetchList);

    return {
      pending,
      items,
      fetchList,
      updateValue,
    };
  },
};
</script>
