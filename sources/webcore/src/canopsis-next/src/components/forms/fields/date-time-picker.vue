<template lang="pug">
  v-menu(
  ref="menu",
  v-model="opened",
  content-class="date-time-picker",
  transition="slide-y-transition",
  max-width="290px",
  :close-on-content-click="false",
  right,
  lazy
  )
    div(slot="activator")
      v-text-field(
      readonly,
      :label="label",
      :error-messages="name ? errors.collect(name) : []",
      :value="value | date('DD/MM/YYYY HH:mm', true)",
      v-validate="rules",
      :data-vv-name="name",
      data-vv-validate-on="none",
      :append-icon="clearable ? 'close' : ''",
      @click:append="clear"
      )
    date-time-picker-field(:value="value", :rules="rules", :roundHours="roundHours", @submit="save")
</template>

<script>
import DateTimePickerField from '@/components/forms/fields/date-picker/date-time-picker-field.vue';

/**
 * Date time picker component
 *
 * @prop {Boolean} [clearable] - if it is true then input field will be have cross button with clear event on click
 * @prop {Date} [value] - v-model
 * @prop {string} [label] - label of the input field
 * @prop {string} [name] - name property in the validation object
 * @prop {string} [rules] - validation rules in vee-validate format
 * @prop {string} [format='DD/MM/YYYY HH:mm'] - date format for display
 *
 * @event value#input
 * @type Date - new date value
 */
export default {
  components: {
    DateTimePickerField,
  },
  props: {
    clearable: Boolean,
    value: Date,
    label: String,
    name: String,
    rules: [String, Object],
    roundHours: Boolean,
  },
  data() {
    return {
      opened: false,
    };
  },
  methods: {
    clear() {
      this.dateTimeObject = null;
      this.dateString = '';
      this.timeString = '';

      this.$emit('input', this.dateTimeObject);
    },

    save(value) {
      this.$emit('input', value);
      this.$refs.menu.save();
    },
  },
};
</script>
