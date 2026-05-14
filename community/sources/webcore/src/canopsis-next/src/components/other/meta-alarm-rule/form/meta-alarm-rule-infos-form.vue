<template>
  <v-layout class="gap-3" column>
    <c-label>{{ $t('common.infos') }}</c-label>
    <meta-alarm-rule-infos-item-form
      v-for="(item, index) in infos"
      v-field="infos[index]"
      :key="item.key"
      :name="item.key"
      @remove="removeItemFromArray(index)"
    />
    <div>
      <v-btn color="primary" outlined @click="add">
        {{ $t('common.add') }}
      </v-btn>
    </div>
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
