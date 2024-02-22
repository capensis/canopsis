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
      :error-messages="errors.collect('file')"
      class="mt-2"
      with-files-list
      @change="changeFiles"
    >
      <template #activator="{ on, disabled }">
        <v-tooltip top>
          <template #activator="{ on: tooltipOn }">
            <v-btn
              :color="errors.has('file') ? 'error' : 'primary'"
              :disabled="disabled"
              small
              outlined
              v-on="{ ...on, ...tooltipOn }"
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
  inject: ['$validator'],
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
  created() {
    this.$validator.attach({
      name: 'file',
      rules: 'required:true',
      getter: () => this.form.file,
      vm: this,
    });
  },
  mounted() {
    this.fetchFieldPbehaviorTypesList();
  },
  beforeDestroy() {
    this.$validator.detach('file');
  },
  methods: {
    changeFiles(files = []) {
      this.updateField('file', files[0]);
      this.$nextTick(() => this.$validator.validate('file'));
    },
  },
};
</script>
