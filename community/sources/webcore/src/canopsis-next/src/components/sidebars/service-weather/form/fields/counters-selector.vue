<template>
  <widget-settings-item :title="$t('settings.counters')">
    <v-layout column>
      <c-enabled-field
        v-field="value.pbehavior_enabled"
        :label="$t('settings.pbehaviorCounters')"
        hide-details
      />
      <c-pbehavior-type-field
        v-field="value.pbehavior_types"
        :required="value.pbehavior_enabled"
        :disabled="!value.pbehavior_enabled"
        :max="$constants.PBEHAVIOR_COUNTERS_LIMIT"
        chips
        multiple
      />
      <c-enabled-field
        v-field="value.state_enabled"
        :label="$t('settings.entityStateCounters')"
        hide-details
      />
      <c-service-weather-state-counter-field
        v-field="value.state_types"
        :required="value.state_enabled"
        :disabled="!value.state_enabled"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [entitiesFieldPbehaviorFieldTypeMixin],
  props: {
    value: {
      type: Object,
      required: false,
    },
  },
  mounted() {
    this.fetchFieldPbehaviorTypesList({
      params: {
        paginate: false,
        with_hidden: true,
        types: [
          PBEHAVIOR_TYPE_TYPES.inactive,
          PBEHAVIOR_TYPE_TYPES.maintenance,
          PBEHAVIOR_TYPE_TYPES.pause,
        ],
      },
    });
  },
};
</script>
