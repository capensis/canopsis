<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('settings.statsNumbers.displayMode') }}
    v-card-text
      v-container
        v-select(:items="displayModes", v-model="form.mode")
        v-card(dark, color="secondary")
          v-card-title {{ $t('common.parameters') }}
          v-card-text
            v-card.my-1(
              v-for="criticity in $constants.STATS_CRITICITY",
              :key="criticity",
              color="secondary darken-1"
            )
              v-card-title {{ criticity }}
              v-card-text
                v-layout(align-center)
                  v-text-field(
                    type="number",
                    :label="$t('common.value')",
                    v-model="form.parameters.criticityLevels[criticity]"
                  )
                  v-layout
                    v-btn(
                      :style="{ backgroundColor: form.parameters.colors[criticity] }",
                      @click="openColorPickerModal(criticity)"
                    ) {{ $t('settings.statsNumbers.selectAColor') }}
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, STATS_DISPLAY_MODE, STATS_DISPLAY_MODE_PARAMETERS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.statsDisplayMode,
  mixins: [modalInnerMixin],
  data() {
    const { displayMode } = this.modal.config;
    const defaultDisplayMode = {
      mode: STATS_DISPLAY_MODE.criticity,
      parameters: STATS_DISPLAY_MODE_PARAMETERS,
    };

    return {
      form: cloneDeep(displayMode || defaultDisplayMode),
    };
  },
  computed: {
    displayModes() {
      return Object.values(STATS_DISPLAY_MODE);
    },
  },
  methods: {
    openColorPickerModal(level) {
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.form.parameters.colors[level],
          type: 'hex',
          action: color => this.$set(this.form.parameters.colors, level, color),
        },
      });
    },

    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);
      }

      this.hideModal();
    },
  },
};
</script>
