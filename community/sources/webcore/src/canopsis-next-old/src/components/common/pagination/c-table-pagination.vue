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
        @input="updatePage"
      />
    </v-flex>
    <v-spacer />
    <v-flex xs2>
      <c-items-per-page-field
        :value="itemsPerPage"
        :items="items"
        class="pa-0"
        hide-details
        @input="updateItemsPerPage"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { getPageForNewItemsPerPage } from '@/helpers/pagination';

export default {
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
    emitInput: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const updatePage = page => emit('update:page', page);
    const updateItemsPerPage = (itemsPerPage) => {
      const page = getPageForNewItemsPerPage(itemsPerPage, props.itemsPerPage, props.page);

      if (page !== props.page) {
        return emit('input', { page, itemsPerPage });
      }

      return emit('update:items-per-page', itemsPerPage);
    };

    return {
      updatePage,
      updateItemsPerPage,
    };
  },
};
</script>
