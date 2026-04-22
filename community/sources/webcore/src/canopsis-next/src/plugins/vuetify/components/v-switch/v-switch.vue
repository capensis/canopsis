<script>
import VSwitch from 'vuetify/lib/components/VSwitch/VSwitch';

export default {
  extends: VSwitch,
  props: {
    confirmDisable: {
      type: Function,
      default: null,
    },
    confirmEnable: {
      type: Function,
      default: null,
    },
  },
  methods: {
    /**
     * Redefined onChange method to confirm the action if the confirmDisable or confirmEnable function is provided.
     */
    async onChange() {
      if (!this.isInteractive) return;
      const { value } = this;
      let input = this.internalValue;

      if (this.isMultiple) {
        if (!Array.isArray(input)) {
          input = [];
        }

        const { length } = input;
        input = input.filter(item => !this.valueComparator(item, value));

        if (input.length === length) {
          input.push(value);
        }
      } else if (this.trueValue !== undefined && this.falseValue !== undefined) {
        input = this.valueComparator(input, this.trueValue) ? this.falseValue : this.trueValue;
      } else if (value) {
        input = this.valueComparator(input, value) ? null : value;
      } else {
        input = !input;
      }

      /**
       * If the confirmDisable function is provided, confirm the disable action.
       */
      if (this.confirmDisable && !input) {
        const isConfirmed = await this.confirmDisable({
          currentValue: true,
          nextValue: false,
        });

        if (!isConfirmed) {
          return;
        }
      }

      /**
       * If the confirmEnable function is provided, confirm the enable action.
       */
      if (this.confirmEnable && input) {
        const isConfirmed = await this.confirmEnable({
          currentValue: false,
          nextValue: true,
        });

        if (!isConfirmed) {
          return;
        }
      }

      this.validate(true, input);
      this.internalValue = input;
      this.hasColor = input;
    },
  },
};
</script>
