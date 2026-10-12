<template>
  <c-layout column>
    <c-label :label="t('common.tag', 2)" />
    <c-layout>
      <c-enabled-field
        v-model="copyFromChildren"
        :label="t('metaAlarmRule.copyTagsFromChildren')"
        hide-details
      />
    </c-layout>
    <c-layout v-if="copyFromChildren" gap="3" align="start">
      <c-enabled-field
        v-model="filterByLabelEnabled"
        :label="t('metaAlarmRule.filterByLabelEnabled')"
        hide-details
      >
        <template #append>
          <c-help-icon
            :text="t('metaAlarmRule.filterByLabelEnabledTooltip')"
            icon="help"
            top
          />
        </template>
      </c-enabled-field>
      <v-expand-transition>
        <c-alarm-tag-field
          v-if="filterByLabelEnabled"
          :value="filterByLabel"
          :label="t('common.label')"
          name="filter_by_label"
          multiple
          addable
          combobox
          required
          hide-selected
          only-labels
          @update:value="updateField('filter_by_label', $event)"
        >
          <template #no-data>
            <v-list-item>
              <v-list-item-title v-html="t('common.pressEnterToApply')" />
            </v-list-item>
          </template>
        </c-alarm-tag-field>
      </v-expand-transition>
    </c-layout>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:form'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);
    const filterByLabelEnabled = ref(!!props.form.filter_by_label?.length);

    const copyFromChildren = computed({
      get: () => props.form.copy_from_children,
      set: value => updateField('copy_from_children', value),
    });
    const filterByLabel = computed(() => props.form.filter_by_label);

    return {
      t,
      copyFromChildren,
      filterByLabel,
      filterByLabelEnabled,
      updateField,
    };
  },
};
</script>
