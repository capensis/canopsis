<template>
  <v-card>
    <v-card-text>
      <c-name-field
        :model-value="form.name"
        :name="nameFieldName"
        required
        @update:model-value="updateField('name', $event)"
      >
        <template #append-outer>
          <c-action-btn
            type="delete"
            @click="$emit('remove', form)"
          />
        </template>
      </c-name-field>
      <c-text-field
        v-if="!form.copy_from_children"
        :model-value="form.description"
        :label="t('common.description')"
        @update:model-value="updateField('description', $event)"
      />
      <v-row class="ga-4 flex-nowrap align-start ma-0" no-gutters>
        <v-col cols="auto">
          <c-enabled-field
            v-model="copyFromChildren"
            :label="t('metaAlarmRule.copyFromLastChild')"
            hide-details
          >
            <template #append>
              <c-help-icon
                :text="t('metaAlarmRule.copyFromLastChildTooltip')"
                icon="help"
                top
              />
            </template>
          </c-enabled-field>
        </v-col>
        <v-col v-if="!form.copy_from_children">
          <c-mixed-field
            :value="form.value"
            :label="t('common.value')"
            @update:value="updateField('value', $event)"
          />
        </v-col>
      </v-row>
    </v-card-text>
  </v-card>
</template>

<script>
import { computed } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: '',
    },
  },

  emits: ['update:form', 'remove'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const nameFieldName = computed(() => `${props.name}.name`);
    const copyFromChildren = computed({
      get: () => props.form.copy_from_children,
      set: value => updateField('copy_from_children', value),
    });

    return {
      t,
      nameFieldName,
      copyFromChildren,
      updateField,
    };
  },
};
</script>
