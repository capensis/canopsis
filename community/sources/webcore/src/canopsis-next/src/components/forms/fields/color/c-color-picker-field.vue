<template lang="pug">
  v-layout(:align-center="splitted")
    template(v-if="splitted")
      v-btn.ml-0(
        :disabled="disabled",
        key="splitted",
        small,
        @click="showColorPickerModal"
      ) {{ label }}
      div.pa-1.text-xs-center(:style="style") {{ color }}
    v-btn.ml-0(
      v-else,
      :style="style",
      :disabled="disabled",
      key="not-splitted",
      @click="showColorPickerModal"
    ) {{ label }}
</template>

<script>
import { MODALS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'color',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      default() {
        return this.$t('common.selectColor');
      },
    },
    color: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: 'hex',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    splitted: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    style() {
      return {
        backgroundColor: this.color,
        color: getMostReadableTextColor(this.color, { level: 'AA', size: 'large' }),
      };
    },
  },
  methods: {
    showColorPickerModal() {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          color: this.color,
          type: this.type,
          action: color => this.updateModel(color),
        },
      });
    },
  },
};
</script>
