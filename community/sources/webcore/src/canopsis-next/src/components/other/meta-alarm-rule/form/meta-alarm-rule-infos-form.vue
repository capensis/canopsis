<template>
  <c-form-block-array-field
    v-field="infos"
    :item-to-form="metaAlarmRuleInfosItemToForm"
    :label="$t('common.infos')"
    :add-button-label="$t('entity.addInfos')"
  >
    <template #item="{ item, index, remove }">
      <meta-alarm-rule-infos-item-form
        v-field="infos[index]"
        :name="item.key"
        @remove="remove"
      />
    </template>
  </c-form-block-array-field>
</template>
<script>
import { metaAlarmRuleInfosItemToForm } from '@/helpers/entities/meta-alarm/rule/form';

import MetaAlarmRuleInfosItemForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-infos-item-form.vue';

export default {
  components: { MetaAlarmRuleInfosItemForm },
  model: {
    prop: 'infos',
    event: 'input',
  },
  props: {
    infos: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    return {
      metaAlarmRuleInfosItemToForm,
    };
  },
};
</script>
