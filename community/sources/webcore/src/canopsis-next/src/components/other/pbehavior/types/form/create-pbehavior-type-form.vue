<template lang="pug">
  v-layout(column)
    c-name-field(v-field="form.name", :disabled="onlyColor")
    v-text-field(
      v-field="form.description",
      v-validate="'required'",
      :label="$t('modals.createPbehaviorType.fields.description')",
      :error-messages="errors.collect('description')",
      :disabled="onlyColor",
      name="description"
    )
    v-layout(row, justify-space-between)
      v-flex(xs6)
        v-select.mr-2(
          v-field="form.type",
          :label="$t('modals.createPbehaviorType.fields.type')",
          :items="types",
          :disabled="onlyColor"
        )
      v-flex.ml-2(xs6)
        c-priority-field(
          v-field="form.priority",
          :disabled="onlyColor",
          required
        )
    c-icon-field(
      v-field="form.icon_name",
      :label="$t('modals.createPbehaviorType.fields.iconName')",
      :hint="$t('modals.createPbehaviorType.iconNameHint')",
      :disabled="onlyColor",
      required
    )
      template(#no-data="")
        v-list-tile
          v-list-tile-content
            v-list-tile-title(v-html="$t('modals.createPbehaviorType.errors.iconName')")
    v-flex.mt-2(xs12)
      v-alert(:value="onlyColor", color="info") {{ $t('pbehaviorTypes.defaultType') }}
    c-color-picker-field.mt-2(v-field="form.color", required)
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
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
    onlyColor: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    types() {
      return Object.values(PBEHAVIOR_TYPE_TYPES);
    },
  },
};
</script>
