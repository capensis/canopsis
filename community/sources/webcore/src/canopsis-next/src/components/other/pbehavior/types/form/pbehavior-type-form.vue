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
        pbehavior-type-default-type-field.mr-2(
          v-field="form.type",
          :disabled="onlyColor",
          @update:color="updateColor"
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
    c-color-picker-field.mt-2(
      v-field="form.color",
      required,
      @input="setColorWasChanged"
    )
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { formMixin } from '@/mixins/form';

import PbehaviorTypeDefaultTypeField from './pbehavior-type-default-type-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorTypeDefaultTypeField },
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
  data() {
    return {
      colorWasChanged: false,
    };
  },
  computed: {
    types() {
      return Object.values(PBEHAVIOR_TYPE_TYPES);
    },
  },
  methods: {
    setColorWasChanged() {
      this.colorWasChanged = true;
    },

    updateColor(color) {
      if (this.colorWasChanged && this.form.color) {
        return;
      }

      this.updateField('color', color);
    },
  },
};
</script>
