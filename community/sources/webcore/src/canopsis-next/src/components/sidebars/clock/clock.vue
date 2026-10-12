<template>
  <widget-settings
    :submitting="submitting"
    :dirty="hasChanges"
    divider
    @submit="submit"
  >
    <field-title v-model="form.title" />

    <widget-settings-item :title="t('clock.fields.display')">
      <c-select
        v-model="form.parameters.display"
        :items="displayItems"
        :label="t('clock.fields.display')"
        item-title="text"
        item-value="value"
        hide-details
      />
    </widget-settings-item>

    <widget-settings-item :title="t('clock.fields.size')">
      <c-select
        v-model="form.parameters.size"
        :items="sizeItems"
        :label="t('clock.fields.size')"
        item-title="text"
        item-value="value"
        hide-details
      />
    </widget-settings-item>

    <widget-settings-item :title="t('clock.fields.spacing')">
      <c-select
        v-model="form.parameters.spacing"
        :items="spacingItems"
        :label="t('clock.fields.spacing')"
        item-title="text"
        item-value="value"
        hide-details
      />
    </widget-settings-item>

    <template v-if="hasDigitalSettings">
      <field-switcher
        v-model="form.parameters.is_12_hour_format"
        :title="t('clock.fields.use12HourFormat')"
      />
      <field-switcher
        v-model="form.parameters.show_seconds"
        :title="t('clock.fields.showSeconds')"
      />
    </template>

    <field-switcher v-model="form.parameters.show_date" :title="t('clock.fields.showDate')" />
    <field-switcher
      v-model="form.parameters.show_utc_offset"
      :title="t('clock.fields.showUtcOffset')"
    />
    <field-switcher
      v-model="form.parameters.show_timezone"
      :title="t('clock.fields.showTimezone')"
    />

    <widget-settings-group :title="t('clock.fields.clocks')">
      <clocks-field v-model="form.parameters.clocks" />
    </widget-settings-group>

    <widget-colors-form v-model="form.parameters.colors" />
  </widget-settings>
</template>

<script setup>
import { computed, toRef } from 'vue';

import {
  CLOCK_DISPLAY_TYPES,
  CLOCK_SIZES,
  CLOCK_SPACING_TYPES,
  SIDE_BARS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useWidgetSettings } from '@/hooks/widget/settings';

import FieldSwitcher from '@/components/sidebars/form/fields/switcher.vue';
import FieldTitle from '@/components/sidebars/form/fields/title.vue';
import WidgetColorsForm from '@/components/sidebars/form/fields/widget-colors-form.vue';
import WidgetSettingsGroup from '@/components/sidebars/partials/widget-settings-group.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import WidgetSettings from '@/components/sidebars/partials/widget-settings.vue';

import ClocksField from './form/fields/clocks-field.vue';

defineOptions({
  name: SIDE_BARS.clockSettings,
});

const props = defineProps({
  sidebar: {
    type: Object,
    required: true,
  },
});

const { t } = useI18n();

const sidebar = toRef(props, 'sidebar');

const { form, hasChanges, submitting, close, submit } = useWidgetSettings({ sidebar });

const displayItems = computed(() => Object.values(CLOCK_DISPLAY_TYPES).map(value => ({
  value,
  text: t(`clock.display.${value}`),
})));

const sizeItems = computed(() => Object.values(CLOCK_SIZES).map(value => ({
  value,
  text: t(`clock.size.${value}`),
})));

const spacingItems = computed(() => Object.values(CLOCK_SPACING_TYPES).map(value => ({
  value,
  text: t(`clock.spacing.${value}`),
})));

const hasDigitalSettings = computed(() => [
  CLOCK_DISPLAY_TYPES.digital,
  CLOCK_DISPLAY_TYPES.sevenSegments,
].includes(form.value.parameters.display));

defineExpose({
  close,
});
</script>
