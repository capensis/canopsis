<template>
  <v-layout class="gap-3" justify-center column>
    <v-layout v-if="splitted" align-center justify-space-between>
      <v-btn
        key="splitted"
        :disabled="disabled"
        class="mr-3"
        small
        @click="showColorPickerModal"
      >
        {{ displayLabel }}
      </v-btn>
      <div
        :style="style"
        class="pa-1 text-center"
      >
        {{ color }}
      </div>
    </v-layout>
    <v-layout align-center justify-end>
      <v-btn
        :style="style"
        :disabled="disabled"
        @click="showColorPickerModal"
      >
        {{ displayLabel }}
      </v-btn>
    </v-layout>

    <v-messages
      v-if="errors.has(name)"
      :value="errors.collect(name)"
      color="error"
    />
  </v-layout>
</template>

<script>
import { computed, watch, nextTick, onBeforeUnmount } from 'vue';
import { Validator } from 'vee-validate';

import { MODALS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useValidationAttachRequired } from '@/hooks/validator/validation-attach-required';
import { useValidator } from '@/hooks/validator/validator';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  model: {
    prop: 'color',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: false,
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
  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateModel } = useModelField(props, emit);
    const modals = useModals();
    const validator = useValidator();
    const { errors } = validator;
    const { attachRequiredRule, detachRequiredRule } = useValidationAttachRequired(props.name);

    const displayLabel = computed(() => props.label ?? t('common.selectColor'));

    const style = computed(() => ({
      backgroundColor: props.color,
      color: getMostReadableTextColor(props.color, { level: 'AA', size: 'large' }),
    }));

    /**
     * Opens the color picker modal for the current value. On confirm, emits the new color and
     * re-validates the field when it is required.
     */
    const showColorPickerModal = () => {
      modals.show({
        name: MODALS.colorPicker,
        config: {
          color: props.color,
          type: props.type,
          action: (newColor) => {
            updateModel(newColor);

            if (props.required) {
              nextTick(() => validator.validate(props.name));
            }
          },
        },
      });
    };

    watch(() => props.required, (required) => {
      if (required && !props.disabled) {
        attachRequiredRule(() => props.color);

        return;
      }

      detachRequiredRule();
    }, { immediate: true });

    onBeforeUnmount(detachRequiredRule);

    return {
      displayLabel,
      style,
      errors,
      showColorPickerModal,
    };
  },
};
</script>
