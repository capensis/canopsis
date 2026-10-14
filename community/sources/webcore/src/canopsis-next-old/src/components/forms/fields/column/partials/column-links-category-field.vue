<template>
  <v-combobox
    v-field="value"
    :items="categories"
    :label="$tc('common.category')"
    :return-object="false"
    :loading="pending"
    item-text="value"
    item-value="value"
    clearable
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { LINK_RULE_TYPES, MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('linkRule');

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
  },
  data() {
    return {
      pending: false,
      categories: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions(['fetchLinkCategoriesWithoutStore']),

    async fetchList() {
      try {
        this.pending = true;

        const params = { limit: MAX_LIMIT, type: LINK_RULE_TYPES.alarm };

        const { categories = [] } = await this.fetchLinkCategoriesWithoutStore({ params });

        if (categories?.length) {
          this.categories = categories.filter(category => !!category);
        }
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
