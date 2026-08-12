<template>
  <c-select
    v-model="selectedCategory"
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
  />
</template>

<script>
import { isString } from 'lodash';

import { MAX_LIMIT } from '@/constants';
import { formBaseMixin } from '@/mixins/form';
import { useEntityCategoryStore } from '@/stores/modules/entity-category';

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

  computed: {
    selectedCategory: {
      get() {
        if (!isString(this.category)) {
          return this.category;
        }

        return this.items.find(({ _id: id }) => id === this.category) ?? this.category;
      },
      set(category) {
        this.updateModel(category);
      },
    },
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
