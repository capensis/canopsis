<template>
  <v-layout column>
    <c-name-field
      v-field="form.name"
      required
    />
    <c-pbehavior-type-field
      v-field="form.type"
      name="type"
      required
    />
    <file-selector
      class="mt-2"
      name="file"
      required
      with-files-list
      @change="changeFiles"
    >
      <template #activator="{ on, disabled }">
        <v-tooltip top>
          <template #activator="{ on: tooltipOn }">
            <v-btn
              v-on="{ ...on, ...tooltipOn }"
              :color="errors.has('file') ? 'error' : 'primary'"
              :disabled="disabled"
              small
              outlined
            >
              <v-icon>cloud_upload</v-icon>
            </v-btn>
          </template>
          <span>{{ $t('common.chooseFile') }}</span>
        </v-tooltip>
      </template>
    </file-selector>
  </v-layout>
</template>

<script>
import { formMixin } from '@/mixins/form';
import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

import FileSelector from '@/components/forms/fields/file-selector.vue';

export default {
  components: { FileSelector },
  mixins: [
    formMixin,
    entitiesFieldPbehaviorFieldTypeMixin,
  ],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  mounted() {
    this.fetchFieldPbehaviorTypesList();
  },
  methods: {
    changeFiles(files = []) {
      this.updateField('file', files[0]);
    },
  },
};
</script>
