<template lang="pug">
  v-form
    v-layout
      v-text-field(
        v-model="form._id",
        v-validate="'required'",
        :label="$t('common.name')",
        :error-messages="errors.collect('name')",
        name="name",
        data-test="name",
        @input="updateField('_id', $event)"
      )
    v-layout
      v-text-field(
        :value="form.description",
        :label="$t('common.description')",
        data-test="description",
        @input="updateField('description', $event)"
      )
    v-layout
      view-selector(
        :value="form.defaultview",
        @input="updateField('defaultview', $event)"
      )
</template>

<script>
import ViewSelector from '@/components/forms/fields/view-selector.vue';

import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  components: { ViewSelector },
  mixins: [formMixin],
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
};
</script>

