<template lang="pug">
  v-layout
    v-flex.pr-1(xs5)
      date-time-splited-picker-field(
        v-validate="startRules",
        :value="value.tstart",
        :fullDay="fullDay",
        :label="startLabel",
        name="tstart",
        @input="updateField('tstart', $event)"
      )
    template(v-if="!noEnding")
      v-flex.pr-1(xs2)
        div.time-dash –
      v-flex(xs5)
        date-time-splited-picker-field(
          v-validate="endRules",
          :value="value.tstop",
          :fullDay="fullDay",
          :label="endLabel",
          name="tstop",
          reverse,
          @input="updateField('tstop', $event)"
        )
</template>

<script>
import formMixin from '@/mixins/form';

import DateTimeSplitedPickerField from '@/components/forms/fields/date-time-picker/date-time-splited-picker-field.vue';

export default {
  components: {
    DateTimeSplitedPickerField,
  },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      required: true,
    },
    endRules: {
      type: Object,
      required: true,
    },
    startRules: {
      type: Object,
      required: true,
    },
    startLabel: {
      type: String,
      required: true,
    },
    endLabel: {
      type: String,
      required: true,
    },
    noEnding: {
      type: Boolean,
      default: false,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.time-dash {
  line-height: 68px;
  padding: 0 8px;
  text-align: center;
}
</style>
