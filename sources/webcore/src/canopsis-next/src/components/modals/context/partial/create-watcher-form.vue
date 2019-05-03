<template lang="pug">
  v-form
    v-layout(wrap, justify-center)
      v-flex(xs11)
        v-text-field(
        :label="$t('modals.createWatcher.displayName')",
        :value="form.name",
        :error-messages="errors.collect('name')",
        data-vv-name="name",
        v-validate="'required'",
        @input="updateField('name', $event)"
        )
    v-layout(wrap, justify-center)
      v-flex(xs11)
        template(v-if="stack === 'go'")
          v-textarea(
          label="Output template",
          :value="form.output_template",
          :error-messages="errors.collect('output_template')",
          data-vv-name="output_template",
          v-validate="'required'",
          @input="updateField('output_template', $event)"
          )
          h3.text-xs-center Pattern
          v-divider
          patterns-list(
          :patterns="form.entities",
          @input="updateField('entities', $event)"
          )
        template(v-else)
          h3.text-xs-center {{ $t('filterEditor.title') }}
          v-divider
          filter-editor(
          :value="form.mfilter",
          required,
          @input="updateField('mfilter', $event)"
          )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesInfoMixin from '@/mixins/entities/info';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: {
    FilterEditor,
    PatternsList,
  },
  mixins: [formMixin, entitiesInfoMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
};
</script>

