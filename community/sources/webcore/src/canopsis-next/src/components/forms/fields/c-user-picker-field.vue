<template lang="pug">
  v-autocomplete(
    v-bind="$attrs",
    v-field="value",
    v-validate="rules",
    :items="items",
    :label="label",
    :loading="pending",
    :name="name",
    :error-messages="errors.collect(name)",
    :return-object="returnObject",
    :item-text="itemText",
    :item-value="itemValue"
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('user');

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  props: {
    value: {
      type: [Object, Array, String],
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
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    permission: {
      type: String,
      default: '',
    },
    itemText: {
      type: String,
      default: 'display_name',
    },
    itemValue: {
      type: String,
      default: '_id',
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
      fetchUsersListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const params = { limit: MAX_LIMIT };

      if (this.permission) {
        params.permission = this.permission;
      }

      const { data: items } = await this.fetchUsersListWithoutStore({ params });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
