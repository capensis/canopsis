<template>
  <v-row v-show="totalItems" align="center" no-gutters>
    <v-col cols="10">
      <c-pagination
        :page="page"
        :limit="itemsPerPage"
        :total="totalItems"
        :total-visible="totalVisible"
        @update:page="updatePage"
      />
    </v-col>
    <v-spacer />
    <v-col cols="2">
      <c-items-per-page-field
        :model-value="itemsPerPage"
        :items="items"
        class="pa-0"
        hide-details
        @update:model-value="updateItemsPerPage"
      />
    </v-col>
  </v-row>
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
    totalVisible: {
      type: Number,
      required: false,
    },
  },

  emits: ['update:page', 'update:items-per-page', 'update:pagination'],

  setup(props, { emit }) {
    const updatePage = page => emit('update:page', page);
    const updateItemsPerPage = (itemsPerPage) => {
      const page = getPageForNewItemsPerPage(itemsPerPage, props.itemsPerPage, props.page);

      if (page !== props.page) {
        return emit('update:pagination', { page, itemsPerPage });
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
