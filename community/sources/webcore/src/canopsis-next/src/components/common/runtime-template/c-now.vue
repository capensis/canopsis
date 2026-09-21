<template>
  <span>{{ value }}</span>
</template>

<script>
import { onBeforeUnmount, onMounted, ref } from 'vue';

import { formatNow } from '@/helpers/date/date-intervals';

export default {
  props: {
    format: {
      type: String,
      default: undefined,
    },
    timezone: {
      type: String,
      required: true,
    },
    offset: {
      type: String,
      default: undefined,
    },
  },
  setup(props) {
    const value = ref('');
    let timer;

    const update = () => {
      value.value = formatNow(props);
      timer = setTimeout(update, 1000);
    };

    onMounted(update);
    onBeforeUnmount(() => clearTimeout(timer));

    return { value };
  },
};
</script>
