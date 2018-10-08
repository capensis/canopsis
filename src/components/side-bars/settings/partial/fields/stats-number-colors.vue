<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.colorsSelector')}}
    v-container
      v-layout(wrap)
        v-flex(xs12)
          v-layout(align-center)
            v-btn(@click="showColorPickerModal('ok')", small) Ok
            div.pa-1.text-xs-center(:style="{ backgroundColor: value['ok'] }") {{ value['ok'] }}
        v-flex(xs12)
          v-layout(align-center)
            v-btn(@click="showColorPickerModal('minor')", small) Minor
            div.pa-1.text-xs-center(:style="{ backgroundColor: value['minor'] }") {{ value['minor'] }}
        v-flex(xs12)
          v-layout(align-center)
            v-btn(@click="showColorPickerModal('major')", small) Major
            div.pa-1.text-xs-center(:style="{ backgroundColor: value['major'] }") {{ value['major'] }}
        v-flex(xs12)
          v-layout(align-center)
            v-btn(@click="showColorPickerModal('critical')", small) Critical
            div.pa-1.text-xs-center(:style="{ backgroundColor: value['critical'] }") {{ value['critical'] }}
</template>

<script>
import { MODALS } from '@/constants';
import modalMixin from '@/mixins/modal/modal';

export default {
  mixins: [modalMixin],
  props: {
    value: {
      type: Object,
      default() {
        return {};
      },
    },
  },
  methods: {
    showColorPickerModal(key) {
      const newValue = { ...this.value };
      this.showModal({
        name: MODALS.colorPicker,
        config: {
          title: 'modals.colorPicker.title',
          action: (color) => {
            newValue[key] = color;
            this.$emit('input', newValue);
          },
        },
      });
    },
  },
};
</script>
