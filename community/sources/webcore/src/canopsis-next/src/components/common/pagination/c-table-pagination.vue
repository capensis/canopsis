<template>
  <v-layout
    v-show="totalItems"
    align-center
  >
    <v-flex xs10>
      <c-pagination
        :page="page"
        :limit="itemsPerPage"
        :total="totalItems"
        @input="$emit('update:page', $event)"
      />
    </v-flex>
    <v-spacer />
    <v-flex xs2>
      <c-items-per-page-field
        class="pa-0"
        :value="itemsPerPage"
        :items="items"
        hide-details
        @input="$emit('update:items-per-page', $event)"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import CPagination from './c-pagination.vue';

export default {
  components: {
    CPagination,
  },
  props: {
    items: {
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
    itemsPerPage: {
      type: Number,
      required: false,
    },
  },
};
</script>
