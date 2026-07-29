<template>
  <div>
    <slot
      :on="{ click: selectFile }"
      :attrs="{ disabled }"
      :disabled="disabled"
      :files="selectedFiles"
      :clear="clear"
      name="activator"
    >
      <v-btn :disabled="disabled" @click="selectFile">
        {{ label || t('common.import') }}
      </v-btn>
    </slot>

    <input
      ref="inputElement"
      :accept="accept"
      :multiple="multiple"
      :disabled="disabled"
      class="d-none"
      type="file"
      @change="change"
    />

    <v-messages v-if="errorMessages.length" :messages="errorMessages" color="error" />
  </div>
</template>

<script>
import { computed, ref, watch } from 'vue';
import { useField } from 'vee-validate';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: [File, Array],
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    accept: {
      type: String,
      default: '',
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    /**
     * File size in kilobytes (KB)
     */
    maxFileSize: {
      type: [Number, String],
      required: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'file',
    },
  },

  emits: ['update:modelValue', 'change'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const inputElement = ref(null);
    const selectedFiles = ref([]);
    const fieldName = computed(() => props.name || 'file');
    const fieldRules = computed(() => ({
      required: props.required,
      size: props.maxFileSize ? Number(props.maxFileSize) : false,
    }));
    const { errors, setValue } = useField(fieldName, fieldRules, {
      initialValue: props.modelValue,
      validateOnMount: false,
      validateOnValueUpdate: false,
    });

    const selectFile = () => {
      inputElement.value?.click();
    };

    const updateFiles = (files) => {
      const value = props.multiple ? files : files[0];

      selectedFiles.value = files;
      setValue(value, false);
      emit('update:modelValue', value);
      emit('change', value);
    };

    const change = ({ target }) => {
      const files = Array.from(target.files ?? []);

      updateFiles(files);
    };

    const dropFiles = (event) => {
      updateFiles(Array.from(event.dataTransfer?.files ?? []));
    };

    const clear = () => {
      if (inputElement.value) {
        inputElement.value.value = '';
      }

      updateFiles([]);
    };

    watch(
      () => props.modelValue,
      (nextValue) => {
        selectedFiles.value = Array.isArray(nextValue) ? nextValue : [nextValue].filter(Boolean);
        setValue(nextValue, false);
      },
    );

    return {
      t,
      inputElement,
      selectedFiles,
      errorMessages: errors,
      selectFile,
      change,
      dropFiles,
      clear,
    };
  },
};
</script>
