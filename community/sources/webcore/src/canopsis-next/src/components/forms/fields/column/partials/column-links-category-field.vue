<template lang="pug">
  v-combobox(
    v-field="value",
    :items="categories",
    :label="$tc('common.category')",
    :return-object="false",
    :loading="pending",
    item-text="value",
    item-value="value"
  )
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

    fetchList() {
      try {
        this.pending = true;

        const params = { limit: MAX_LIMIT, type: LINK_RULE_TYPES.alarm };

        const { data } = this.fetchLinkCategoriesWithoutStore({ params });

        this.categories = data;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
