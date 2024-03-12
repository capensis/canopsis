<template>
  <v-layout wrap>
    <v-flex
      v-show="label"
      xs12
    >
      <v-layout>
        <h4 class="text-subtitle-1 grey--text text--darken-2">
          {{ label }}
        </h4>
        <c-help-icon
          v-if="helpText"
          :text="helpText"
          icon-class="ml-2 storage-help-tooltip"
          right
        />
      </v-layout>
    </v-flex>
    <v-flex xs12>
      <v-layout column>
        <c-storage-field
          v-for="(storage, index) in storages"
          :key="storage.key"
          :value="storage.directory"
          :disabled="disabled"
          @edit="$emit('edit', storage)"
          @remove="removeItemFromArray(index)"
        />
      </v-layout>
    </v-flex>
    <v-flex xs12>
      <v-layout>
        <v-btn
          :disabled="disabled"
          class="ml-0"
          color="primary"
          @click="$emit('add')"
        >
          {{ $t('common.add') }}
        </v-btn>
      </v-layout>
    </v-flex>
  </v-layout>
</template>

<script>
import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'storages',
    event: 'input',
  },
  props: {
    storages: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      required: false,
    },
  },
};
</script>

<style lang="scss">
.storage-help-tooltip {
  pointer-events: auto !important;
}
</style>
