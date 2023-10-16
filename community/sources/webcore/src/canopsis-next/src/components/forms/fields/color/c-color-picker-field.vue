<template>
  <v-layout align-center>
    <template v-if="splitted">
      <v-btn
        class="ml-0"
        :disabled="disabled"
        key="splitted"
        small
        @click="showColorPickerModal"
      >
        {{ label }}
      </v-btn>
      <div
        class="pa-1 text-center"
        :style="style"
      >
        {{ color }}
      </div>
    </template>
    <v-btn
      class="ml-0"
      v-else
      :style="style"
      :disabled="disabled"
      key="not-splitted"
      @click="showColorPickerModal"
    >
      {{ label }}
    </v-btn>
    <v-messages
      v-if="errors.has(name)"
      :value="errors.collect(name)"
      color="error"
    />
  </v-layout>
</template>

<script>
import { Validator } from 'vee-validate';

import { MODALS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';

import { formBaseMixin } from '@/mixins/form';
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  mixins: [
    formBaseMixin,
    validationAttachRequiredMixin,
  ],
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
    name: {
      type: String,
      default: 'color',
    },
    required: {
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
  watch: {
    required: {
      immediate: true,
      handler(required) {
        if (required && !this.disabled) {
          this.attachRequiredRule();

          return;
        }

        this.detachRequiredRule();
      },
    },
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    showColorPickerModal() {
      this.$modals.show({
        name: MODALS.colorPicker,
        config: {
          color: this.color,
          type: this.type,
          action: (color) => {
            this.updateModel(color);

            if (this.required) {
              this.$nextTick(() => this.$validator.validate(this.name));
            }
          },
        },
      });
    },
  },
};
</script>
