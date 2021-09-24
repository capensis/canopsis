<template lang="pug">
  v-autocomplete(
    v-field="value",
    v-validate="rules",
    :items="items",
    :label="label",
    :loading="pending",
    :name="name",
    :error-messages="errors.collect(name)",
    item-text="name",
    item-value="_id",
    return-object
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('role');

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'role',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      items: [],
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchRolesListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { data: items } = await this.fetchRolesListWithoutStore({ params: { limit: MAX_LIMIT } });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
