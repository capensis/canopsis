<template>
  <v-layout
    class="flowchart-color-field"
    justify-space-between
    align-center
  >
    <v-checkbox
      v-if="!hideCheckbox"
      :input-value="isFilled"
      :label="label || $t('flowchart.color')"
      class="mt-0"
      color="primary"
      hide-details
      @change="updateIsFilled"
    />
    <span
      v-else
      class="v-label"
    >
      {{ label }}
    </span>
    <v-flex
      v-show="isFilled || hideCheckbox"
      xs3
    >
      <v-btn
        :style="style"
        class="flowchart-color-field__button ma-0 pa-0"
        small
        block
        @click="showColorPickerModal"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { COLORS } from '@/config';
import { MODALS } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: 'transparent',
    },
    palette: {
      type: Array,
      default: () => COLORS.flowchart.shapes,
    },
    label: {
      type: String,
      required: false,
    },
    hideCheckbox: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isFilled() {
      return this.value !== 'transparent';
    },

    style() {
      return {
        backgroundColor: this.value,
      };
    },
  },
  methods: {
    updateIsFilled(checked) {
      this.updateModel(checked ? this.palette[0] : 'transparent');
    },

    showColorPickerModal() {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          color: this.value,
          palette: this.palette,
          action: color => this.updateModel(color),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.flowchart-color-field {
  &__button {
    min-width: unset;
  }
}
</style>
