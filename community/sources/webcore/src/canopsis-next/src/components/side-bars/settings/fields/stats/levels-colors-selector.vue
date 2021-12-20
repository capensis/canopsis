<template lang="pug">
  v-list-group(data-test="levelsColorsSelector")
    v-list-tile(slot="activator") {{$t('settings.colorsSelector.title')}}
    v-container
      v-layout(wrap)
        v-flex(
          v-for="level in $constants.STATS_CRITICITY",
          xs12,
          :data-test="`levelsColor-${level}`",
          :key="level"
        )
          c-color-picker-field(
            v-field="levelsColors[level]",
            :label="getButtonText(level)",
            :type="colorType",
            splitted
          )
</template>

<script>
import { ALARM_STATS_CALENDAR_COLORS, STATS_CRITICITY } from '@/constants';

export default {
  model: {
    prop: 'levelsColors',
    event: 'input',
  },
  props: {
    levelsColors: {
      type: Object,
      default: () => ({ ...ALARM_STATS_CALENDAR_COLORS }),
    },
    hideSuffix: {
      type: Boolean,
      default: false,
    },
    colorType: {
      type: String,
      default: 'rgba',
    },
  },
  methods: {
    getButtonText(key) {
      let suffix = '';

      if (!this.hideSuffix) {
        if (key === STATS_CRITICITY.ok) {
          suffix = ` / ${this.$t('common.yes')}`;
        } else if (key === STATS_CRITICITY.critical) {
          suffix = ` / ${this.$t('common.no')}`;
        }
      }

      return this.$t(`settings.colorsSelector.statsCriticity.${key}`) + suffix;
    },
  },
};
</script>
