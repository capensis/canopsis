<template>
  <v-layout column>
    <span class="text-subtitle-1 font-weight-bold mb-2">{{ $t('common.infos') }}</span>
    <meta-alarm-rule-infos-item-form
      v-for="(item, index) in infos"
      v-field="infos[index]"
      :key="item.key"
      :name="item.key"
      @remove="removeItemFromArray(index)"
    />
    <v-layout class="my-4">
      <v-btn color="primary" @click="add">
        {{ $t('common.add') }}
      </v-btn>
    </v-layout>
  </v-layout>
</template>
<script>
import { metaAlarmRuleInfosItemToForm } from '@/helpers/entities/meta-alarm/rule/form';

import { useArrayModelField } from '@/hooks/form/array-model-field';

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
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);

    const add = () => addItemIntoArray(metaAlarmRuleInfosItemToForm());

    return {
      add,
      removeItemFromArray,
    };
  },
};
</script>
