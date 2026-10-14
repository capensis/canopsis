<template>
  <v-autocomplete
    v-field="value"
    v-validate="rules"
    v-bind="$attrs"
    :items="items"
    :label="label"
    :loading="pending"
    :name="name"
    :error-messages="errors.collect(name)"
    :return-object="returnObject"
    :item-text="itemText"
    :item-value="itemValue"
  />
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useUser } from '@/hooks/store/modules/user';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
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
  setup(props) {
    const items = ref([]);

    const rules = computed(() => ({
      required: props.required,
    }));

    const { fetchUsersListWithoutStore } = useUser();

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const params = { limit: MAX_LIMIT };

      if (props.permission) {
        params.permission = props.permission;
      }

      const { data } = await fetchUsersListWithoutStore({ params });

      items.value = data;
    });

    onMounted(fetchList);

    return {
      pending,
      items,

      rules,

      fetchList,
    };
  },
};
</script>
