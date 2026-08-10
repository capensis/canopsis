<template>
  <c-select
    :model-value="category"
    :label="$t('common.category')"
    :loading="pending"
    :items="items"
    :name="name"
    :hide-details="hideDetails"
    :clearable="!required"
    class="mt-0"
    item-title="name"
    item-value="_id"
    return-object
    @update:model-value="updateModel"
  />
</template>

<script>
import { MAX_LIMIT } from '@/constants';
import { useEntityCategoryStore } from '@/stores/modules/entity-category';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],

  props: {
    category: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'category',
    },
    addable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['input', 'update:category'],

  data() {
    return {
      pending: false,
      items: [],
    };
  },

  mounted() {
    this.fetchList();
  },

  methods: {
    async fetchList() {
      this.pending = true;

      const { data: items } = await useEntityCategoryStore().fetchCategoriesListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
