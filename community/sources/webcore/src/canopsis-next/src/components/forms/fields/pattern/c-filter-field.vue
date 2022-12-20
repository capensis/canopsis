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
              :value="!!item.old_entity_patterns",
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
import { MAX_LIMIT } from '@/constants';

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
    fetchList() {
      if (!this.filtersPending) {
        this.fetchFiltersList({ params: { limit: MAX_LIMIT } });
      }
    },
  },
};
</script>
