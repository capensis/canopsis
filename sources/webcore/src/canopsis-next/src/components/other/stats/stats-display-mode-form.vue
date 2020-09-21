<template lang="pug">
  div
    v-container(data-test="statsDisplayModeParameters")
      v-select(v-field="form.mode", :items="displayModes")
      v-card(dark, color="secondary")
        v-card-title {{ $t('common.parameters') }}
        v-card-text
          v-card.my-1(
            v-for="criticity in $constants.STATS_CRITICITY",
            :data-test="`statsDisplayMode-${criticity}`",
            :key="criticity",
            color="secondary darken-1"
          )
            v-card-title {{ criticity }}
            v-card-text
              v-layout(align-center)
                v-text-field(
                  v-field="form.parameters.criticityLevels[criticity]",
                  :label="$t('common.value')",
                  type="number",
                  data-test="displayModeValue"
                )
                v-layout
                  v-btn(
                    data-test="displayModeColorPicker",
                    :style="{ backgroundColor: form.parameters.colors[criticity] }",
                    @click="openColorPickerModal(criticity)"
                  ) {{ $t('settings.statsNumbers.selectAColor') }}
</template>

<script>
import { MODALS, STATS_DISPLAY_MODE } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    displayModes() {
      return Object.values(STATS_DISPLAY_MODE);
    },
  },
  methods: {
    openColorPickerModal(level) {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.form.parameters.colors[level],
          type: 'hex',
          action: color => this.updateField(['parameters', 'colors', level], color),
        },
      });
    },
  },
};
</script>
