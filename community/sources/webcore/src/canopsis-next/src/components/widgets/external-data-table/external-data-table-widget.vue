<template>
  <v-layout>
    <c-advanced-data-table
      items=""
      headers=""
    />
  </v-layout>
</template>

<script>
import CAdvancedDataTable from '@/components/common/table/c-advanced-data-table.vue';

export default {
  components: { CAdvancedDataTable },
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  setup() {},
};
</script>
