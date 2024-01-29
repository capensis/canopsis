<template>
  <file-selector
    ref="fileSelector"
    v-bind="$attrs"
    :name="name"
    @change="$emit('change', $event)"
  >
    <template #activator="{ clear, on, files }">
      <slot
        v-if="$scopedSlots.selected && files.length"
        :clear="clear"
        name="selected"
      />
      <v-layout
        v-else
        :class="{ 'drag-zone--dropping': dropping }"
        class="pa-3 drag-zone"
        align-center
        @drop.prevent="drop"
        @dragover.prevent="dragover"
        @dragleave.prevent="dragleave"
      >
        <v-icon
          class="mr-3"
          large
        >
          upload_file
        </v-icon>
        <span class="text-subtitle-2">
          {{ $t('common.fileSelector.dragAndDrop.label') }}
          <a v-on="on">
            {{ $t('common.fileSelector.dragAndDrop.labelAction') }}
          </a>
          {{ fileTypeLabel ? ` ${fileTypeLabel}` : '' }}
        </span>
      </v-layout>
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
      this.$nextTick(() => this.dropping = false);
    },
  },
};
</script>

<style lang="scss" scoped>
.drag-zone {
  position: relative;
  border-radius: 10px;
  border: 3px dashed;

  &--dropping:after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1;
    opacity: .4;

    .theme--light & {
      background: var(--v-application-background-darken2);
    }

    .theme--dark & {
      background: var(--v-application-background-lighten3);
    }
  }

  .theme--light & {
    border-color: var(--v-application-background-darken2);

    &, & ::v-deep .v-icon {
      color: var(--v-application-background-darken2);
    }
  }

  .theme--dark & {
    border-color: var(--v-application-background-lighten3);

    &, & ::v-deep .v-icon {
      color: var(--v-application-background-lighten4);
    }
  }
}
</style>
