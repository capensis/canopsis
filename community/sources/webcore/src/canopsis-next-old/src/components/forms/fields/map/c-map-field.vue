<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :items="items"
    :label="label || $tc('common.map')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :error-messages="errors.collect(name)"
    :hide-details="hideDetails"
    item-text="name"
    item-value="_id"
    clearable
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

const { mapActions } = createNamespacedHelpers('map');

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
      default: 'map',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
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
      fetchMapsListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      this.pending = true;

      const { data: items } = await this.fetchMapsListWithoutStore({ params: { limit: MAX_LIMIT } });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
