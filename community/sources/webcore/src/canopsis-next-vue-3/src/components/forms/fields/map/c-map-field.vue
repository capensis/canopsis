<template>
  <c-select
    :model-value="value"
    :items="items"
    :label="label || $t('common.map')"
    :loading="pending"
    :disabled="disabled"
    :name="name"
    :hide-details="hideDetails"
    item-title="name"
    item-value="_id"
    clearable
    @update:model-value="updateModel"
  />
</template>

<script>
import { MAX_LIMIT } from '@/constants';
import { useMapStore } from '@/stores/modules/map';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  emits: ['input', 'update:value'],
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
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: items } = await useMapStore().fetchListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.items = items;
      this.pending = false;
    },
  },
};
</script>
