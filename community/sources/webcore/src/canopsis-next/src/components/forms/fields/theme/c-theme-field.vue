<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :items="themes"
    :label="label || $tc('common.theme')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :error-messages="errors.collect(name)"
    :hide-details="hideDetails"
    :clearable="clearable"
    item-text="name"
    item-value="_id"
  />
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';

import { useTheme } from '@/hooks/store/modules/theme';
import { usePendingHandler } from '@/hooks/query/pending';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'map',
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
  setup(props) {
    const themes = ref([]);

    const { fetchThemesListWithoutStore } = useTheme();

    const rules = computed(() => ({ required: props.required }));

    const fetchListHandler = async () => {
      try {
        const { data: themesData } = await fetchThemesListWithoutStore({ params: { limit: MAX_LIMIT } });

        themes.value = themesData;
      } catch (err) {
        console.error(err);
      }
    };

    const { pending, handler: fetchList } = usePendingHandler(fetchListHandler);

    onMounted(() => {
      fetchList();
    });

    return {
      pending,
      themes,
      rules,
    };
  },
};
</script>
