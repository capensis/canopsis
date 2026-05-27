<template>
  <v-layout class="gap-3" column>
    <c-name-field
      v-field="form.name"
      disabled
    />
    <c-form-block>
      <c-form-block-row :label="$t('common.description')">
        <c-description-field v-field="form.description" autofocus />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.enabled')">
        <c-enabled-field v-field="form.enabled" />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.availabilityState')">
        <c-alarm-state-field
          v-field="form.sli_avail_state"
          :label="$t('entity.availabilityState')"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.impact')">
        <c-impact-level-field
          v-field="form.impact_level"
          required
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.type')">
        <c-entity-type-field
          v-field="form.type"
          required
          disabled
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.coordinates')">
        <c-coordinates-field
          v-field="form.coordinates"
          row
        />
      </c-form-block-row>

      <c-form-block-row
        v-if="hasStateSetting"
        :label="$t('stateSetting.title')"
      >
        <entity-state-setting
          :form="form"
          :preparer="prepareStateSettingForm"
        />
      </c-form-block-row>
    </c-form-block>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { isEntityComponentType } from '@/helpers/entities/entity/form';

import EntityStateSetting from '@/components/other/state-setting/entity-state-setting.vue';

export default {
  inject: ['$validator'],
  components: {
    EntityStateSetting,
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
    prepareStateSettingForm: {
      type: Function,
      default: data => data,
    },
  },
  setup(props) {
    const hasStateSetting = computed(() => isEntityComponentType(props.form.type));

    return {
      hasStateSetting,
    };
  },
};
</script>
