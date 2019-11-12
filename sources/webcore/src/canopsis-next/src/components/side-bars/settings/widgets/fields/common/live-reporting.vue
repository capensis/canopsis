<template lang="pug">
  settings-button-field(
    :isEmpty="isValueEmpty",
    @create="showEditLiveReportingModal",
    @edit="showEditLiveReportingModal",
    @delete="removeLiveReporting"
  )
    .subheading(slot="title", data-test="liveReporting") {{ $t('settings.liveReporting.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import formBaseMixin from '@/mixins/form/base';

import SettingsButtonField from '../partials/button-field.vue';

export default {
  components: { SettingsButtonField },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isValueEmpty() {
      const { tstart, tstop } = this.value || {};

      return !tstart && !tstop;
    },
  },
  methods: {
    showEditLiveReportingModal() {
      this.$modals.show({
        name: MODALS.editLiveReporting,
        config: {
          ...this.value,

          action: value => this.updateModel(value),
        },
      });
    },

    removeLiveReporting() {
      this.updateModel({});
    },
  },
};
</script>
