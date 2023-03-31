<template lang="pug">
  v-layout(v-show="totalItems", align-center)
    v-flex(xs10)
      c-pagination(
        :page="page",
        :limit="rowsPerPage",
        :total="totalItems",
        @input="$emit('update:page', $event)"
      )
    v-spacer
    v-flex(xs2)
      c-records-per-page-field.pa-0(
        :value="rowsPerPage",
        :items="rowsPerPageItems",
        hide-details,
        @input="$emit('update:rows-per-page', $event)"
      )
</template>

<script>
import CPagination from './c-pagination.vue';

export default {
  components: {
    CPagination,
  },
  props: {
    rowsPerPageItems: {
      type: Array,
      required: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    page: {
      type: Number,
      required: false,
    },
    rowsPerPage: {
      type: Number,
      required: false,
    },
  },
};
</script>
