<template>
  <c-form-block-array-field
    :form="infos"
    :item-to-form="metaAlarmRuleInfosItemToForm"
    :label="t('common.infos')"
    :add-button-label="t('entity.addInfos')"
    @update:form="$emit('update:infos', $event)"
  >
    <template #item="{ index, remove }">
      <meta-alarm-rule-infos-item-form
        :form="infos[index]"
        :name="infos[index].key"
        @update:form="updateItemInArray(index, $event)"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>

<script>
import { metaAlarmRuleInfosItemToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { useI18n } from '@/hooks/i18n';

import MetaAlarmRuleInfosItemForm from './meta-alarm-rule-infos-item-form.vue';

export default {
  components: { MetaAlarmRuleInfosItemForm },

  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:infos'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateItemInArray } = useArrayModelField(props, emit);

    return {
      t,
      metaAlarmRuleInfosItemToForm,
      updateItemInArray,
    };
  },
};
</script>
