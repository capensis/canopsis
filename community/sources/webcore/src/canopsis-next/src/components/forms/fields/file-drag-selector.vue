<template>
  <file-selector
    v-validate="rules"
    v-bind="$attrs"
    :error-messages="errors.collect(name)"
    :name="name"
    accept=".svg"
    @change="$emit('change', $event)"
  >
    <template #activator="{ clear, on, drop, files }">
      <slot
        v-if="$scopedSlots.selected && files.length"
        :clear="clear"
        name="selected"
      />
      <v-layout
        class="pa-3 grey--text drag-zone"
        v-else
        align-center
        @drop.prevent="drop"
        @dragover.prevent=""
      >
        <v-icon
          class="mr-3"
          color="grey"
          large
        >
          upload_file
        </v-icon>
        <span class="text-subtitle-2">
          {{ $t('common.fileSelector.dragAndDrop.label') }}
          <a v-on="on">{{ $t('common.fileSelector.dragAndDrop.labelAction') }}</a>
          {{ fileTypeLabel ? ` ${fileTypeLabel}` : '' }}
        </span>
      </v-layout>
    </template>
  </file-selector>
</template>

<script>
import FileSelector from './file-selector.vue';

export default {
  inject: ['$validator'],
  components: { FileSelector },
  inheritAttrs: false,
  props: {
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'file',
    },
    fileTypeLabel: {
      type: String,
      default: '',
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.drag-zone {
  border-radius: 10px;
  border: 3px dashed #9e9e9e77;
}
</style>
