<template>
  <v-select
    v-field="value"
    :items="items"
    :label="$t('common.criteria')"
    item-text="label"
    item-value="id"
    hide-details
    return-object
  />
</template>

<script>
import { isUndefined } from 'lodash';

import { MAX_LIMIT } from '@/constants';

import { formMixin } from '@/mixins/form';
import { entitiesRatingSettingsMixin } from '@/mixins/entities/rating-settings';

export default {
  mixins: [formMixin, entitiesRatingSettingsMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
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

      this.pending = true;
    },

    setAvailableValue() {
      const [firstRatingSetting] = this.items;

      if (!this.items.length) {
        this.updateModel(undefined);
      } else if (isUndefined(this.value) || !this.items.some(({ id }) => id === this.value.id)) {
        this.updateModel(firstRatingSetting);
      }
    },
  },
};
</script>
