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
        c-priority-field(v-field="form.priority", required)
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
    v-layout.mt-2(wrap)
      v-flex(xs12)
        v-switch(
          v-field="form.isSpecialColor",
          :label="$t('modals.createPbehaviorType.fields.isSpecialColor')",
          color="primary"
        )
      v-flex(xs12)
        c-color-picker-field(
          v-field="form.color",
          :disabled="!form.isSpecialColor"
        )
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import IconField from '@/components/forms/fields/icon-field.vue';

export default {
  inject: ['$validator'],
  components: { IconField },
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
