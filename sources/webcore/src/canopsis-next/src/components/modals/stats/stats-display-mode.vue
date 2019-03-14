<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Stats - Display mode
    v-card-text
      v-container
        v-select(:items="displayModes", v-model="form.mode")
        v-card(dark, color="secondary")
          v-card-title {{ $t('common.parameters') }}
          v-card-text
            template(v-for="criticity in $constants.STATS_CRITICITY", align-center, justify-space-between)
              v-card.my-1(color="secondary darken-1")
                v-card-title {{ criticity }}
                v-card-text
                  v-layout(align-center)
                    v-text-field(type="number", label="Value", v-model="form.parameters.criticityLevels[criticity]")
                    v-layout
                      v-btn(
                      :style="{ backgroundColor: form.parameters.colors[criticity] }",
                      @click="openColorPickerModal(criticity)"
                      ) Select a color
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { find } from 'lodash';

import { MODALS, STATS_DISPLAY_MODE } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.statsDisplayMode,
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        mode: STATS_DISPLAY_MODE.criticity,
        parameters: {
          criticityLevels: {
            ok: 0,
            minor: 10,
            major: 20,
            critical: 30,
          },
          colors: {
            ok: '#66BB6A',
            minor: '#FFEE58',
            major: '#FFA726',
            critical: '#FF7043',
          },
        },
      },
    };
  },
  computed: {
    displayModes() {
      return Object.values(STATS_DISPLAY_MODE);
    },
  },
  mounted() {
    if (this.config.displayMode) {
      const mode = find(STATS_DISPLAY_MODE, displayMode => displayMode === this.config.displayMode.type);
      this.form = { ...this.config.displayMode, mode };
    }
  },
  methods: {
    openColorPickerModal(level) {
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.form.parameters.colors[level],
          type: 'hex',
          action: (color) => {
            this.$set(this.form.parameters.colors, level, color);
          },
        },
      });
    },

    async submit() {
      if (this.config.action) {
        await this.config.action(this.form);

        this.hideModal();
      }
    },
  },
};
</script>
