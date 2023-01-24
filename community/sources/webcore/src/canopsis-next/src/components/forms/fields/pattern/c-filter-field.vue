<template lang="pug">
  v-autocomplete(
    v-field="value",
    :items="filters",
    :label="label || $t('common.filters')",
    :loading="filtersPending",
    :disabled="disabled",
    :name="name",
    item-text="name",
    item-value="_id",
    hide-details,
    clearable
  )
    template(#item="{ item, tile }")
      v-list-tile(v-bind="tile.props", v-on="tile.on")
        v-list-tile-content
          v-list-tile-title.v-list-badge__tile__title
            v-badge(
              :value="isOldPattern(item)",
              color="error",
              overlap
            )
              template(#badge="")
                v-tooltip(top)
                  template(#activator="{ on: badgeTooltipOn }")
                    v-icon(v-on="badgeTooltipOn", color="white") priority_high
                  span {{ $t('pattern.oldPatternTooltip') }}
              span {{ item.name }}
</template>

<script>
import { MAX_LIMIT, OLD_PATTERNS_FIELDS } from '@/constants';

import { isOldPattern } from '@/helpers/pattern';

import { entitiesFilterMixin } from '@/mixins/entities/filter';

export default {
  mixins: [entitiesFilterMixin],
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
      default: 'filter',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    isOldPattern(item) {
      return isOldPattern(item, [OLD_PATTERNS_FIELDS.entity]);
    },

    fetchList() {
      if (!this.filtersPending) {
        this.fetchFiltersList({ params: { limit: MAX_LIMIT } });
      }
    },
  },
};
</script>
