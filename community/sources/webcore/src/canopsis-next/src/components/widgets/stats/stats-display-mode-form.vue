<template lang="pug">
  div
    v-container(data-test="statsDisplayModeParameters")
      v-select(v-field="form.mode", :items="displayModes")
      v-card(dark, color="secondary")
        v-card-title {{ $tc('common.parameter', 2) }}
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
                c-color-picker-field.ml-2(
                  v-field="form.parameters.colors[criticity]",
                  :label="$t('settings.statsNumbers.selectAColor')"
                )
</template>

<script>
import { STATS_DISPLAY_MODE } from '@/constants';

export default {
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
};
</script>
