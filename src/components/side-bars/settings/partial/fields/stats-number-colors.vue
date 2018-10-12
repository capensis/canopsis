<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.colorsSelector.title')}}
    v-container
      v-layout(wrap)
        v-flex(xs12, v-for="(value, key) in value", :key="key")
          v-layout(align-center)
            v-btn(@click="showColorPickerModal(key)", small)
              template {{ $t('settings.colorsSelector.statsCriticity.' + key) }}
              template(v-if="key === statsCriticity.ok") / {{ $t('common.yes') }}
              template(v-if="key === statsCriticity.critical") / {{ $t('common.no') }}
            div.pa-1.text-xs-center(:style="{ backgroundColor: value }") {{ value }}
</template>

<script>
import { MODALS, STATS_CRITICITY } from '@/constants';
import modalMixin from '@/mixins/modal/modal';
import formMixin from '@/mixins/form';

export default {
  mixins: [modalMixin, formMixin],
  props: {
    value: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  computed: {
    statsCriticity() {
      return { ...STATS_CRITICITY };
    },
  },
  methods: {
    showColorPickerModal(key) {
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: 'modals.colorPicker.title',
          action: (color) => {
            this.updateField(key, color);
          },
        },
      });
    },
  },
};
</script>
