<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.colorsSelector.title')}}
    v-container
      v-layout(wrap)
        v-flex(xs12, v-for="level in statsCriticity", :key="level")
          v-layout(align-center)
            v-btn(@click="showColorPickerModal(level)", small) {{ getButtonText(level) }}
            div.pa-1.text-xs-center(:style="{ backgroundColor: levelsColors[level] }") {{ levelsColors[level] }}
</template>

<script>
import { MODALS, STATS_CRITICITY, STATS_CALENDAR_COLORS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import formMixin from '@/mixins/form';

export default {
  mixins: [modalMixin, formMixin],
  model: {
    prop: 'levelsColors',
    event: 'input',
  },
  props: {
    levelsColors: {
      type: Object,
      default: () => ({ ...STATS_CALENDAR_COLORS.alarm }),
    },
    hideSuffix: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    statsCriticity() {
      return { ...STATS_CRITICITY };
    },
    getButtonText() {
      return (key) => {
        let suffix = '';

        if (!this.hideSuffix) {
          if (key === STATS_CRITICITY.ok) {
            suffix = ` / ${this.$t('common.yes')}`;
          } else if (key === STATS_CRITICITY.critical) {
            suffix = ` / ${this.$t('common.no')}`;
          }
        }

        return this.$t(`settings.colorsSelector.statsCriticity.${key}`) + suffix;
      };
    },
  },
  methods: {
    showColorPickerModal(level) {
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: this.$t('modals.colorPicker.title'),
          color: this.levelsColors[level],
          action: (color) => {
            this.updateField(level, color);
          },
        },
      });
    },
  },
};
</script>
