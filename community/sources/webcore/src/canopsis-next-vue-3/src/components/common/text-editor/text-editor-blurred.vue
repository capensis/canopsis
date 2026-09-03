<template>
  <div class="text-editor-blurred" @click="$emit('click', $event)">
    <v-input :error-messages="errorMessages" :hide-details="hideDetails">
      <v-field
        :label="label"
        :dirty="Boolean(value)"
        :disabled="disabled"
        :error="errorMessages.length > 0"
        variant="filled"
        active
      >
        <template #default="{ props: fieldProps }">
          <div v-bind="fieldProps" @click="wrapperClickHandler">
            <c-compiled-template
              :template="value"
              :class="{ 'text-editor-blurred__disabled': disabled }"
            />
          </div>
        </template>
      </v-field>
    </v-input>
  </div>
</template>

<script>
import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    errorMessages: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['click'],

  setup() {
    const modals = useModals();

    const wrapperClickHandler = (event) => {
      if (event.target.tagName !== 'IMG') {
        return;
      }

      event.preventDefault();
      event.stopPropagation();

      modals.show({
        name: MODALS.imageViewer,
        config: {
          src: event.target.getAttribute('src'),
        },
      });
    };

    return {
      wrapperClickHandler,
    };
  },
};
</script>

<style lang="scss" scoped>
.text-editor-blurred {
  :deep(.v-field__input) {
    min-height: 150px;
    align-items: flex-start;
    overflow: auto;
  }

  :deep(img) {
    cursor: pointer !important;
  }
}

.text-editor-blurred__disabled {
  color: rgba(var(--v-theme-on-surface), 0.38);

  :deep(img) {
    pointer-events: all !important;
  }
}
</style>
