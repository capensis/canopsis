<template>
  <c-select
    :model-value="value"
    :items="items"
    :label="$t('common.criteria')"
    :loading="pending"
    item-title="label"
    item-value="id"
    hide-details
    return-object
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import { isUndefined } from 'lodash';

import { MAX_LIMIT } from '@/constants';
import { useRatingSettings } from '@/hooks/store/modules/rating-settings';

export default {
  props: {
    value: {
      type: Object,
      required: false,
    },
    mandatory: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

  setup() {
    const { ratingSettingsUpdatedAt, fetchRatingSettingsListWithoutStore } = useRatingSettings();

    return {
      ratingSettingsUpdatedAt,
      fetchRatingSettingsListWithoutStore,
    };
  },

  data() {
    return {
      pending: false,
      items: [],
    };
  },

  watch: {
    ratingSettingsUpdatedAt: 'fetchList',
  },

  mounted() {
    this.fetchList();
  },

  methods: {
    async fetchList() {
      try {
        this.pending = true;

        const { data: items } = await this.fetchRatingSettingsListWithoutStore({
          params: {
            limit: MAX_LIMIT,
            enabled: true,
          },
        });

        this.items = items;

        if (this.mandatory) {
          this.setAvailableValue();
        }
      } finally {
        this.pending = false;
      }
    },

    setAvailableValue() {
      const [firstRatingSetting] = this.items;

      if (!this.items.length) {
        this.$emit('update:value', undefined);
      } else if (isUndefined(this.value) || !this.items.some(({ id }) => id === this.value.id)) {
        this.$emit('update:value', firstRatingSetting);
      }
    },
  },
};
</script>
