<template>
  <c-layout class="flowchart-color-field" justify-space-between align-center>
    <v-checkbox
      v-if="!hideCheckbox"
      :model-value="isFilled"
      :label="label || $t('common.color')"
      class="mt-0"
      color="primary"
      hide-details
      @update:model-value="updateIsFilled"
    />
    <span v-else class="v-label">
      {{ label }}
    </span>
    <v-col v-show="isFilled || hideCheckbox" cols="3">
      <v-btn
        :style="style"
        class="flowchart-color-field__button ma-0 pa-0"
        size="small"
        block
        @click="showColorPickerModal"
      />
    </v-col>
  </c-layout>
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
  emits: ['input', 'update:value'],
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
