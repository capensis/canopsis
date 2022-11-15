<template lang="pug">
  component(
    v-field="value",
    v-validate="rules",
    :is="component",
    :items="availableRoles",
    :label="label || $tc('common.role')"
    :loading="pending",
    :name="name",
    :error-messages="errors.collect(name)",
    :disabled="disabled",
    item-text="name",
    item-value="_id",
    return-object
  )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isObject } from 'lodash';

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
    disabled: {
      type: Boolean,
      default: false,
    },
    autocomplete: {
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
    component() {
      if (this.autocomplete) {
        return 'v-autocomplete';
      }

      return 'v-select';
    },

    availableRoles() {
      return !this.items.length && isObject(this.value) ? [this.value] : this.items;
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
  watch: {
    disabled(value) {
      if (!value && !this.items.length) {
        this.fetchList();
      }
    },
  },
  mounted() {
    if (!this.disabled) {
      this.fetchList();
    }
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
