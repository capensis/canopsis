<template lang="pug">
  v-layout.white(v-show="totalItems", align-center)
    v-flex(xs10)
      c-pagination(
        :page="page",
        :limit="rowsPerPage",
        :total="totalItems",
        @input="$emit('update:page', $event)"
      )
    v-spacer
    v-flex(xs2)
      c-records-per-page(
        :value="rowsPerPage",
        :items="rowsPerPageItems",
        @input="$emit('update:rows-per-page', $event)"
      )
</template>

<script>
import CPagination from './c-pagination.vue';
import CRecordsPerPage from './c-records-per-page.vue';

export default {
  components: {
    CPagination,
    CRecordsPerPage,
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
