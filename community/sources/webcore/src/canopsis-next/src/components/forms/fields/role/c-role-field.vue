<template>
  <component
    :is="component"
    v-field="value"
    v-validate="rules"
    :items="availableRoles"
    :label="label || $tc('common.role')"
    :loading="pending"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :item-disabled="isDisabledItems"
    :multiple="multiple"
    :chips="chips"
    :small-chips="chips"
    item-text="name"
    item-value="_id"
    return-object
  />
</template>

<script>
import { ref, computed, watch, onMounted } from 'vue';
import { isArray, isObject } from 'lodash';

import { MAX_LIMIT } from '@/constants';

import { usePendingHandler } from '@/hooks/query/pending';
import { useRole } from '@/hooks/store/modules/role';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Object, String, Array],
      required: false,
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
    autocomplete: {
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
  setup(props) {
    const items = ref([]);

    const rules = computed(() => ({
      required: props.required,
    }));

    const component = computed(() => (props.autocomplete ? 'v-autocomplete' : 'v-select'));

    const availableRoles = computed(() => {
      if (!items.value.length) {
        if (isArray(props.value)) {
          return props.value;
        }

        if (isObject(props.value)) {
          return [props.value];
        }
      }

      return items.value;
    });

    const { fetchRolesListWithoutStore } = useRole();

    const { pending, handler: fetchList } = usePendingHandler(async () => {
      const params = { limit: MAX_LIMIT };

      if (props.permission) {
        params.permission = props.permission;
      }

      const { data } = await fetchRolesListWithoutStore({ params });

      items.value = data;
    });

    watch(() => props.disabled, (newValue) => {
      if (!newValue && !items.value.length) {
        fetchList();
      }
    });

    onMounted(() => {
      if (!props.disabled) {
        fetchList();
      }
    });

    return {
      pending,
      items,
      component,
      availableRoles,
      rules,

      fetchList,
    };
  },
};
</script>
