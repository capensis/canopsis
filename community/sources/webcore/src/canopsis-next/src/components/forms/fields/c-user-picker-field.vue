<template lang="pug">
  v-combobox(
    v-field="value",
    v-validate="'required'",
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

import formBaseMixin from '@/mixins/form/base';

const { mapActions } = createNamespacedHelpers('user');

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Object,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'user',
    },
  },
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
    ...mapActions({
      fetchUsersListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { data: items } = await this.fetchUsersListWithoutStore({ params: { limit: MAX_LIMIT } });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
