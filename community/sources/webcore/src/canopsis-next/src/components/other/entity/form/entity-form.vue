<template lang="pug">
  v-tabs(slider-color="primary", centered)
    v-tab {{ $t('entity.fields.form') }}
    v-tab-item
      v-layout.mt-3(column)
        v-layout(row)
          v-text-field(
            v-field="form.name",
            v-validate="'required'",
            :label="$t('common.name')",
            :error-messages="errors.collect('name')",
            name="name",
            disabled
          )
        v-layout(row)
          v-textarea(
            v-field="form.description",
            :label="$t('common.description')",
            :error-messages="errors.collect('description')",
            name="description"
          )
        v-layout(row, justify-space-between)
          v-flex(xs4)
            c-enabled-field(v-field="form.enabled")
          v-flex(xs8)
            v-layout(row)
              v-flex(xs4)
                c-impact-level-field.mr-3(v-field="form.impact_level", required)
              v-flex(xs8)
                c-entity-type-field(v-field="form.type", required, disabled)
        v-layout(wrap)
          v-flex(xs12)
            entities-select(
              v-field="form.impact",
              :label="$t('entity.fields.impact')",
              :disabled-entities-ids="form.disabled_impact",
              name="impact"
            )
          v-flex(xs12)
            entities-select(
              v-field="form.depends",
              :label="$t('entity.fields.depends')",
              :disabled-entities-ids="form.disabled_depends",
              :existing-entities-ids="form.impact",
              name="depends"
            )
    v-tab {{ $t('entity.fields.manageInfos') }}
    v-tab-item
      manage-infos(v-field="form.infos")
</template>

<script>
import ManageInfos from '@/components/widgets/context/manage-infos.vue';
import EntitiesSelect from '@/components/widgets/context/entities-select.vue';

export default {
  inject: ['$validator'],
  components: {
    EntitiesSelect,
    ManageInfos,
  },
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
