<template>
  <c-combobox
    :model-value="modelValue"
    :items="categories"
    :label="$t('common.category')"
    :loading="pending"
    clearable
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { mapActions } from 'pinia';

import { LINK_RULE_TYPES, MAX_LIMIT } from '@/constants';
import { useLinkRuleStore } from '@/stores/modules/link-rule';

export default {
  props: {
    modelValue: {
      type: String,
      default: '',
    },
  },

  emits: ['update:modelValue'],

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
    ...mapActions(useLinkRuleStore, {
      fetchLinkRuleCategoriesWithoutStore: 'fetchLinkCategoriesWithoutStore',
    }),

    async fetchList() {
      try {
        this.pending = true;

        const params = { limit: MAX_LIMIT, type: LINK_RULE_TYPES.alarm };
        const { categories = [] } = await this.fetchLinkRuleCategoriesWithoutStore({ params });

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
