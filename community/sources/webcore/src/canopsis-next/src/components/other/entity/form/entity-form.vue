<template lang="pug">
  v-tabs(slider-color="primary", centered)
    v-tab {{ $t('entity.form') }}
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
          c-description-field(v-field="form.description")
        v-layout(row, justify-space-between)
          v-flex(xs3)
            c-enabled-field(v-field="form.enabled")
          v-flex(xs9)
            v-layout(row)
              v-flex.pr-3(xs3)
                c-impact-level-field(v-field="form.impact_level", required)
              v-flex.pr-3(xs3)
                c-entity-state-field(
                  v-field="form.sli_avail_state",
                  :label="$t('entity.availabilityState')",
                  required
                )
              v-flex(xs6)
                c-entity-type-field(v-field="form.type", required, disabled)
        v-layout(wrap)
          v-flex(xs12)
            entities-select(
              v-field="form.impact",
              :label="$t('entity.impact')",
              :disabled-entities-ids="form.disabled_impact",
              name="impact"
            )
          v-flex(xs12)
            entities-select(
              v-field="form.depends",
              :label="$t('entity.depends')",
              :disabled-entities-ids="form.disabled_depends",
              :existing-entities-ids="form.impact",
              name="depends"
            )
    v-tab {{ $t('entity.manageInfos') }}
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
