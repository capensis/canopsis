<template>
  <file-selector ref="fileSelector" v-bind="$attrs" :name="name" @change="$emit('change', $event)">
    <template #activator="{ clear, on, files }">
      <slot v-if="$slots.selected && files.length" :clear="clear" :files="files" name="selected" />
      <c-layout
        v-else
        :class="{ 'drag-zone--dropping': dropping }"
        class="pa-3 drag-zone"
        align="center"
        @drop.prevent="drop"
        @dragover.prevent="dragover"
        @dragleave.prevent="dragleave"
      >
        <v-icon class="mr-3" size="large"> upload_file </v-icon>
        <slot :on="on" name="label">
          <span class="text-subtitle-2">
            {{ $t('common.fileSelector.dragAndDrop.label') }}
            <a @click="on.click">
              {{ $t('common.fileSelector.dragAndDrop.labelAction') }}
            </a>
            {{ fileTypeLabel ? ` ${fileTypeLabel}` : '' }}
          </span>
        </slot>
      </c-layout>
    </template>
  </file-selector>
</template>

<script>
import FileSelector from './file-selector.vue';

export default {
  components: { FileSelector },
  inheritAttrs: false,

  props: {
    name: {
      type: String,
      default: 'file',
    },
    fileTypeLabel: {
      type: String,
      default: '',
    },
  },

  emits: ['change'],

  data() {
    return {
      dropping: false,
    };
  },

  methods: {
    drop(event) {
      this.$refs.fileSelector.dropFiles(event);

      this.dragleave();
    },

    dragover() {
      this.dropping = true;
    },

    dragleave() {
      this.$nextTick(() => (this.dropping = false));
    },
  },
};
</script>

<style lang="scss" scoped>
.drag-zone {
  position: relative;
  border-radius: 10px;
  border: 3px dashed rgb(var(--v-theme-secondary));

  &--dropping:after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
    opacity: 0.4;
    background: rgb(var(--v-theme-secondary));
  }
}
</style>
