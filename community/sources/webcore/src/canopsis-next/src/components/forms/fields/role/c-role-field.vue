<template>
  <component
    :is="component"
    v-field="value"
    v-validate="rules"
    :items="availableRoles"
    :label="label || $tc('common.role')"
    :loading="pending"
    :name="name"
    :error-messages="errors.collect(name)"
    :disabled="disabled"
    :multiple="multiple"
    :chips="chips"
    :small-chips="chips"
    item-text="name"
    item-value="_id"
    return-object
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import { isArray, isObject } from 'lodash';

import { MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('role');

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: [Object, String, Array],
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
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    permission: {
      type: String,
      default: '',
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
      if (!this.items.length) {
        if (isArray(this.value)) {
          return this.value;
        }

        if (isObject(this.value)) {
          return [this.value];
        }
      }

      return this.items;
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

      const params = { limit: MAX_LIMIT };

      if (this.permission) {
        params.permission = this.permission;
      }

      const { data: items } = await this.fetchRolesListWithoutStore({ params });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
