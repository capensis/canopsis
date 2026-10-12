<template>
  <v-select
    v-field="value"
    :items="languages"
    :disabled="disabled"
    :label="label"
    :menu-props="menuProps"
  />
</template>

<script>
import { computed } from 'vue';

import { useComponentInstance } from '@/hooks/vue';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'language',
    },
  },
  setup() {
    const menuProps = { offsetY: true };

    const instance = useComponentInstance();

    const languages = computed(() => Object.keys(instance.$i18n.messages));

    return {
      menuProps,

      languages,
    };
  },
};
</script>
