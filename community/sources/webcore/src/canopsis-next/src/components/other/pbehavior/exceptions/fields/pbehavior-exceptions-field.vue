<template>
  <v-layout column>
    <v-layout class="mt-3">
      <v-flex xs12>
        <slot
          v-if="!exdates.length"
          name="no-data"
        />
        <pbehavior-exception-field
          v-field="exdates[index]"
          v-for="(exdate, index) in exdates"
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
import { uid } from '@/helpers/uid';
import { convertDateToStartOfDayDateObject } from '@/helpers/date/date';

import { formArrayMixin } from '@/mixins/form';
import { entitiesFieldPbehaviorFieldTypeMixin } from '@/mixins/entities/pbehavior/types-field';

import PbehaviorExceptionField from '@/components/other/pbehavior/exceptions/fields/pbehavior-exception-field.vue';

export default {
  inject: ['$validator'],
  components: { PbehaviorExceptionField },
  mixins: [
    formArrayMixin,
    entitiesFieldPbehaviorFieldTypeMixin,
  ],
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
  mounted() {
    this.fetchFieldPbehaviorTypesList();
  },
  methods: {
    addExceptionDate() {
      this.addItemIntoArray({
        key: uid(),
        begin: convertDateToStartOfDayDateObject(),
        end: convertDateToStartOfDayDateObject(),
        type: '',
      });
    },
  },
};
</script>
