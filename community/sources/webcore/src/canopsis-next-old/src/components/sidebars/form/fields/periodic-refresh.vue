<template>
  <widget-settings-item
    :title="$t('settings.periodicRefresh')"
    optional
  >
    <periodic-refresh-field
      v-field="form.periodic_refresh"
      :name="name"
    />
    <live-watching-field
      v-if="withLiveWatching"
      v-field="form.liveWatching"
    />
  </widget-settings-item>
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { formMixin } from '@/mixins/form';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

import LiveWatchingField from './live-watching.vue';

export default {
  inject: ['$validator'],
  components: {
    PeriodicRefreshField,
    WidgetSettingsItem,
    LiveWatchingField,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      required: false,
    },
    withLiveWatching: {
      type: Boolean,
      default: false,
    },
  },
  created() {
    if (!this.form?.periodic_refresh?.unit) {
      this.updateField('periodic_refresh.unit', TIME_UNITS.second);
    }
  },
};
</script>
