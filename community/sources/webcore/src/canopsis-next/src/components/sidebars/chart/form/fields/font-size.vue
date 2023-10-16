<template>
  <widget-settings-item :title="$t('settings.chart.fontSize')">
    <v-radio-group
      class="pt-0 mt-0"
      v-model="enabled"
      :name="name"
      hide-details
    >
      <v-radio
        :value="false"
        :label="$t('settings.chart.auto')"
        color="primary"
      />
      <v-radio
        :value="true"
        :label="$t('settings.chart.manual')"
        color="primary"
      />
    </v-radio-group>
    <c-number-field
      v-if="enabled"
      v-field="value"
      required
    />
  </widget-settings-item>
</template>

<script>
import { isUndefined } from 'lodash';

import { NUMBERS_CHART_DEFAULT_FONT_SIZE } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
    name: {
      type: String,
      default: 'font_size',
    },
  },
  computed: {
    enabled: {
      get() {
        return !isUndefined(this.value);
      },
      set(value) {
        this.updateModel(value ? NUMBERS_CHART_DEFAULT_FONT_SIZE : undefined);
      },
    },
  },
};
</script>
