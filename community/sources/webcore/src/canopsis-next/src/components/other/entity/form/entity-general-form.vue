<template>
  <c-layout gap="3" column>
    <c-layout gap="2">
      <c-name-field
        :model-value="form.name"
        class="flex-grow-1"
        disabled
        @update:model-value="updateField('name', $event)"
      />

      <c-entity-type-field
        :value="form.type"
        class="flex-grow-1"
        disabled
        @update:value="updateField('type', $event)"
      />
    </c-layout>
    <c-form-block>
      <c-form-block-row :label="$t('common.description')" indented>
        <c-description-field
          :model-value="form.description"
          autofocus
          @update:model-value="updateField('description', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.enabled')">
        <c-enabled-field
          :model-value="form.enabled"
          @update:model-value="updateField('enabled', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.availabilityState')">
        <c-alarm-state-field
          :value="form.sli_avail_state"
          :label="$t('entity.availabilityState')"
          required
          @update:value="updateField('sli_avail_state', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('entity.impact')">
        <c-impact-level-field
          :value="form.impact_level"
          required
          @update:value="updateField('impact_level', $event)"
        />
      </c-form-block-row>

      <c-form-block-row :label="$t('common.coordinates')">
        <c-coordinates-field
          :value="form.coordinates"
          row
          @update:value="updateField('coordinates', $event)"
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
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { isEntityComponentType } from '@/helpers/entities/entity/form';

import { useModelField } from '@/hooks/form/model-field';

import EntityStateSetting from '@/components/other/state-setting/entity-state-setting.vue';

export default {
  components: {
    EntityStateSetting,
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

  emits: ['update:form'],

  setup(props, { emit }) {
    const { updateField } = useModelField(props, emit);
    const hasStateSetting = computed(() => isEntityComponentType(props.form.type));

    return {
      hasStateSetting,
      updateField,
    };
  },
};
</script>
