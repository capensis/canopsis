<template lang="pug">
  v-select(
    v-field="value",
    :items="items",
    :label="$t('common.criteria')",
    hide-details,
    item-text="label",
    item-value="id",
    return-object
  )
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
        if (isUndefined(this.value)) {
          if (this.items.length) {
            const [firstRatingSetting] = this.items;

            this.updateModel(firstRatingSetting);
          }
        } else if (!this.items.some(({ id }) => id === this.value)) {
          this.updateModel(undefined);
        }
      }

      this.pending = true;
    },
  },
};
</script>
