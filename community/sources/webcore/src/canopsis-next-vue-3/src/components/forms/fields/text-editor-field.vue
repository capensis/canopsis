<template>
  <text-editor
    v-model="fieldValue"
    :label="label"
    :public="isPublic"
    :buttons="buttons"
    :config="config"
    :extra-buttons="extraButtons"
    :error-messages="fieldErrorMessages"
    :max-file-size="preparedMaxFileSize"
    :variables="variables"
    :sanitize-options="sanitizeOptions"
    :autofocus="autofocus"
    :required="required"
    @blur="handleBlur"
  />
</template>

<script>
import { computed, watch } from 'vue';
import { useField } from 'vee-validate';

import { useInfo } from '@/hooks/store/modules/info';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

export default {
  components: {
    TextEditor,
  },

  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    value: {
      type: String,
      default: undefined,
    },
    name: {
      type: String,
      default: 'text',
    },
    label: {
      type: String,
      default: '',
    },
    public: {
      type: Boolean,
      default: false,
    },
    buttons: {
      type: Array,
      default: undefined,
    },
    extraButtons: {
      type: Array,
      default: undefined,
    },
    config: {
      type: Object,
      default: () => ({}),
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
    maxFileSize: {
      type: Number,
      default: 0,
    },
    variables: {
      type: Array,
      required: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    sanitizeOptions: {
      type: Object,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { fileUploadMaxSize } = useInfo();

    const validateValue = (value) => {
      if (props.required && !value) {
        return `${props.label || props.name} is required`;
      }

      return true;
    };

    const { value, errorMessage, handleBlur, setValue, validate } = useField(
      computed(() => props.name),
      validateValue,
      {
        initialValue: props.modelValue ?? props.value ?? '',
        validateOnValueUpdate: false,
      },
    );

    const fieldValue = computed({
      get: () => props.modelValue ?? props.value ?? value.value,
      set: (nextValue) => {
        setValue(nextValue, false);
        emit('update:modelValue', nextValue);
        validate();
      },
    });

    const isPublic = computed(() => props.public);
    const preparedMaxFileSize = computed(() => props.maxFileSize || fileUploadMaxSize.value);
    const fieldErrorMessages = computed(() => [
      ...props.errorMessages,
      ...(errorMessage.value ? [errorMessage.value] : []),
    ]);

    watch(
      () => props.modelValue,
      (nextValue) => {
        if (nextValue !== undefined && nextValue !== value.value) {
          setValue(nextValue, false);
          validate();
        }
      },
    );

    watch(
      () => props.value,
      (nextValue) => {
        if (nextValue !== undefined && nextValue !== value.value) {
          setValue(nextValue, false);
          validate();
        }
      },
    );

    return {
      fieldValue,
      fieldErrorMessages,
      isPublic,
      preparedMaxFileSize,
      handleBlur,
    };
  },
};
</script>
