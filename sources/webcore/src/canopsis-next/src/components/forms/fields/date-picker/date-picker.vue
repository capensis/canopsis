<script>
import VDatePicker from 'vuetify/es5/components/VDatePicker';

import DatePickerYears from './date-picker-years.vue';

export default {
  extends: VDatePicker,
  props: {
    ...VDatePicker.props,

    yearFirst: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    const data = VDatePicker.data.call(this);

    return {
      ...data,

      activePicker: this.yearFirst ? 'YEAR' : this.type.toUpperCase(),
    };
  },
  methods: {
    genYears() {
      return this.$createElement(DatePickerYears, {
        props: {
          color: this.color,
          format: this.yearFormat,
          locale: this.locale,
          min: this.minYear,
          max: this.maxYear,
          value: `${this.tableYear}`,
        },
        on: {
          input: this.yearClick,
        },
        ref: 'years',
      });
    },
  },
};
</script>
