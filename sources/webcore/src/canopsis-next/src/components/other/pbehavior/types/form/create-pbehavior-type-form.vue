<template lang="pug">
  div
    v-layout(row)
      v-text-field(
        v-field="form.name",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorType.fields.name')",
        :error-messages="errors.collect('name')",
        name="name"
      )
    v-layout(row)
      v-text-field(
        v-field="form.description",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorType.fields.description')",
        :error-messages="errors.collect('description')",
        name="description"
      )
    v-layout(row, justify-space-between)
      v-flex.mr-2(xs6)
        v-select(
          v-field="form.type",
          :label="$t('modals.createPbehaviorType.fields.type')",
          :items="types"
        )
      v-flex.ml-2(xs6)
        v-text-field(
          v-field.number="form.priority",
          v-validate="'required|numeric|min_value:0'",
          :label="$t('modals.createPbehaviorType.fields.priority')",
          :error-messages="errors.collect('priority')",
          type="number",
          name="priority"
        )
    v-layout(row)
      icon-field(
        v-field="form.iconName",
        v-validate="'required'",
        :label="$t('modals.createPbehaviorType.fields.iconName')",
        :hint="$t('modals.createPbehaviorType.iconNameHint')"
      )
        template(slot="no-data")
          v-list-tile
            v-list-tile-content
              v-list-tile-title(v-html="$t('modals.createPbehaviorType.errors.iconName')")
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import formMixin from '@/mixins/form';

import IconField from '@/components/forms/fields/icon-field.vue';

export default {
  components: { IconField },
  inject: ['$validator'],
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
  computed: {
    types() {
      return Object.values(PBEHAVIOR_TYPE_TYPES);
    },
  },
};
</script>
