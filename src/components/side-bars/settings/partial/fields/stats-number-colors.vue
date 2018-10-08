<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.colorsSelector')}}
    v-container
      v-layout(wrap)
        v-flex(xs6, v-for="(criticity, key) in value", :key="key")
          v-layout(align-center)
            v-btn(@click="showColorPickerModal(key)", small) {{ key }}
            div.pa-1.text-xs-center(:style="{ backgroundColor: value[key] }") {{ value[key] }}
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
