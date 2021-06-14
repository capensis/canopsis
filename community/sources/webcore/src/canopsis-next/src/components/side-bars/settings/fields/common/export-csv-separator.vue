<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      v-layout(row)
        v-select(v-field="value", :items="separators")
</template>

<script>
import { EXPORT_CSV_SEPARATORS } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: EXPORT_CSV_SEPARATORS.comma,
    },
    title: {
      type: String,
      required: true,
    },
  },
  computed: {
    separators() {
      return Object.values(EXPORT_CSV_SEPARATORS);
    },
  },
};
</script>
