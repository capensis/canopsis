<template>
  <v-layout class="gap-3" column>
    <v-layout v-if="exdates.length">
      <v-flex xs12>
        <pbehavior-exception-field
          v-for="(exdate, index) in exdates"
          v-field="exdates[index]"
          :key="exdate.key"
          :disabled="disabled"
          :with-type="withExdateType"
          class="mb-3"
          @delete="removeItemFromArray(index)"
        />
      </v-flex>
    </v-layout>
    <v-layout v-if="!disabled">
      <slot name="actions">
        <v-flex>
          <v-btn
            class="ml-0"
            color="secondary"
            @click="addExceptionDate"
          >
            {{ $t('modals.createPbehaviorException.addDate') }}
          </v-btn>
        </v-flex>
      </slot>
    </v-layout>
  </v-layout>
</template>

<script>
import { onMounted } from 'vue';

import { uid } from '@/helpers/uid';
import { convertDateToStartOfDayDateObject, convertDateToEndOfDayDateObject } from '@/helpers/date/date';

import { useArrayModelField } from '@/hooks/form/array-model-field';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';

export default {
  components: { PbehaviorExceptionField },
  model: {
    prop: 'exdates',
    event: 'input',
  },
  props: {
    exdates: {
      type: Array,
      default: () => [],
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    withExdateType: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { addItemIntoArray, removeItemFromArray } = useArrayModelField(props, emit);
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    onMounted(fetchPbehaviorTypesFieldList);

    const addExceptionDate = () => {
      addItemIntoArray({
        key: uid(),
        begin: convertDateToStartOfDayDateObject(),
        end: convertDateToEndOfDayDateObject(),
        type: '',
      });
    };

    return {
      addExceptionDate,
      removeItemFromArray,
    };
  },
};
</script>
