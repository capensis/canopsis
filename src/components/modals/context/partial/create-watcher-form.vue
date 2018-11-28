<template lang="pug">
v-form
  v-layout(wrap, justify-center)
    v-flex(xs11)
      v-text-field(
        :label="$t('modals.createWatcher.displayName')",
        v-model="form.name",
        data-vv-name="name",
        v-validate="'required'",
        :error-messages="errors.collect('name')",
      )
  v-layout(wrap, justify-center)
    v-flex(xs11)
      h3.text-xs-center {{ $t('filterEditor.title') }}
      v-divider
      filter-editor(v-model="form.mfilter")
</template>

<script>
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  components: {
    FilterEditor,
  },
  mixins: [formMixin],
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

