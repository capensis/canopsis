<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", required)
    c-pbehavior-type-field(
      v-field="form.type",
      name="type",
      required
    )
    file-selector.mt-2(
      :error-messages="errors.collect('file')",
      with-files-list,
      @change="changeFiles"
    )
      template(#activator="{ on, disabled }")
        v-tooltip(top)
          template(#activator="{ on: tooltipOn }")
            v-btn.ma-0(
              v-on="{ ...on, ...tooltipOn }",
              :color="errors.has('file') ? 'error' : 'primary'",
              :disabled="disabled",
              small,
              outline
            )
              v-icon cloud_upload
          span {{ $t('common.chooseFile') }}
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
