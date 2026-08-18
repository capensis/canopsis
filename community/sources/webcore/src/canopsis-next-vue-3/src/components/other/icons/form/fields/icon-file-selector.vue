<template>
  <file-drag-selector
    :file-type-label="$t('common.fileSelector.fileTypes.svg')"
    :model-value="file"
    :max-file-size="maxFileSize"
    accept=".svg,.svg+xml"
    required
    @change="chooseIcon"
  >
    <template #selected="{ clear, files }">
      <c-layout class="icon-file-selector__selected cursor-default" align="center">
        <div class="icon-file-selector__preview">
          <custom-icon :content="content" />
        </div>

        <div class="icon-file-selector__name">
          {{ files[0]?.name }}
        </div>

        <div class="icon-file-selector__actions">
          <c-action-btn type="delete" @click.prevent="clear" />
        </div>
      </c-layout>
    </template>
  </file-drag-selector>
</template>

<script>
import { onBeforeUnmount, onMounted, ref } from 'vue';

import { registerCustomIcon, unregisterCustomIcon } from '@/plugins/vuetify/icons';

import { getFileTextContent } from '@/helpers/file/file-select';
import { uid } from '@/helpers/uid';

import FileDragSelector from '@/components/forms/fields/file-drag-selector.vue';

import CustomIcon from '../../partials/custom-icon.vue';

export default {
  components: {
    CustomIcon,
    FileDragSelector,
  },

  props: {
    file: {
      type: File,
      default: null,
    },
    maxFileSize: {
      type: [Number, String],
      required: false,
    },
  },

  emits: ['update:file'],

  setup(props, { emit }) {
    const uniqueIconName = ref(uid('icon'));
    const content = ref('');

    const parseIcon = async (file) => {
      if (!file) {
        content.value = '';
        unregisterCustomIcon(uniqueIconName.value);

        return;
      }

      const template = await getFileTextContent(file);

      content.value = template;
      registerCustomIcon(uniqueIconName.value, template);
    };

    const chooseIcon = (file = null) => {
      parseIcon(file);
      emit('update:file', file);
    };

    onMounted(() => parseIcon(props.file));
    onBeforeUnmount(() => unregisterCustomIcon(uniqueIconName.value));

    return {
      content,
      uniqueIconName,

      chooseIcon,
    };
  },
};
</script>

<style lang="scss">
.icon-file-selector {
  &__selected {
    width: 100%;
    min-width: 0;
  }

  &__preview {
    flex: 0 0 48px;
  }

  &__name {
    min-width: 0;
    flex: 1 1 auto;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  &__actions {
    flex: 0 0 40px;
    text-align: right;
  }
}
</style>
