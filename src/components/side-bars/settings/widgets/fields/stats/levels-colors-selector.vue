<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.colorsSelector.title')}}
    v-container
      v-layout(wrap)
        v-flex(xs12, v-for="level in $constants.STATS_CRITICITY", :key="level")
          v-layout(align-center)
            v-btn(@click="showColorPickerModal(level)", small) {{ getButtonText(level) }}
            div.pa-1.text-xs-center(:style="{ backgroundColor: levelsColors[level] }") {{ levelsColors[level] }}
</template>

<script>
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
      default: () => ({ ...this.$constants.STATS_CALENDAR_COLORS.alarm }),
    },
    hideSuffix: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    getButtonText() {
      return (key) => {
        let suffix = '';

        if (!this.hideSuffix) {
          if (key === this.$constants.STATS_CRITICITY.ok) {
            suffix = ` / ${this.$t('common.yes')}`;
          } else if (key === this.$constants.STATS_CRITICITY.critical) {
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
        name: this.$constants.MODALS.colorPicker,
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
